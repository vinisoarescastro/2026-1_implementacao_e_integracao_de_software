// Package jdk é responsável por detectar ou provisionar automaticamente
// o JDK 21 necessário para executar o assinador.jar.
package jdk

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	jdkSubDir      = "jdk"
	minVersion     = "21"
	releaseJSONURL = "https://raw.githubusercontent.com/kyriosdata/runner/main/release.json"
	httpTimeout    = 5 * time.Minute
)

type jreURLs struct {
	WindowsX64 string `json:"windows_x64"`
	LinuxX64   string `json:"linux_x64"`
	MacX64     string `json:"mac_x64"`
}

type jreManifest struct {
	JRE jreURLs `json:"jre"`
}

// Resolve retorna o caminho para o executável java disponível.
// Ordem de busca:
//  1. ~/.hubsaude/jdk/bin/java  (provisionado pelo Runner)
//  2. java no PATH do sistema
//  3. Download automático do JRE via Eclipse Temurin (Adoptium)
func Resolve() (string, error) {
	// 1. JDK já provisionado pelo Runner
	managed, err := managedJavaPath()
	if err == nil && isJavaValid(managed) {
		return managed, nil
	}

	// 2. java no PATH do sistema
	if systemJava, err := exec.LookPath("java"); err == nil && isJavaValid(systemJava) {
		return systemJava, nil
	}

	// 3. Download automático
	fmt.Fprintf(os.Stderr, "JDK %s não encontrado. Iniciando download automático do JRE...\n", minVersion)
	javaPath, err := downloadJRE()
	if err != nil {
		return "", fmt.Errorf(
			"JDK %s não encontrado e download automático falhou: %w\n"+
				"  Instale manualmente em https://adoptium.net/ e adicione ao PATH.",
			minVersion, err,
		)
	}
	return javaPath, nil
}

// downloadJRE baixa o JRE 21 do Eclipse Temurin e extrai em ~/.hubsaude/jdk/.
func downloadJRE() (string, error) {
	url, err := jreDownloadURL()
	if err != nil {
		return "", err
	}

	hubsaudeDir, err := HubSaudeDir()
	if err != nil {
		return "", err
	}

	jdkDir := filepath.Join(hubsaudeDir, jdkSubDir)

	fmt.Fprintf(os.Stderr, "Baixando JRE de %s\n", url)
	tmpFile, err := downloadToTemp(url)
	if err != nil {
		return "", fmt.Errorf("falha no download: %w", err)
	}
	defer os.Remove(tmpFile)

	fmt.Fprintf(os.Stderr, "Extraindo JRE em %s\n", jdkDir)
	if err := os.MkdirAll(jdkDir, 0755); err != nil {
		return "", fmt.Errorf("não foi possível criar %s: %w", jdkDir, err)
	}

	if runtime.GOOS == "windows" {
		if err := extractZip(tmpFile, jdkDir); err != nil {
			return "", fmt.Errorf("falha ao extrair zip: %w", err)
		}
	} else {
		if err := extractTarGz(tmpFile, jdkDir); err != nil {
			return "", fmt.Errorf("falha ao extrair tar.gz: %w", err)
		}
	}

	// O Temurin extrai em subdiretório — encontra o bin/java dentro de jdkDir
	javaPath, err := findJavaBinary(jdkDir)
	if err != nil {
		return "", fmt.Errorf("JRE extraído mas java não encontrado: %w", err)
	}

	fmt.Fprintf(os.Stderr, "JRE instalado em %s\n", javaPath)
	return javaPath, nil
}

// jreDownloadURL retorna a URL de download do JRE para a plataforma atual via release.json.
func jreDownloadURL() (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(releaseJSONURL)
	if err != nil {
		return adoptiumFallbackURL(), nil // fallback direto para a API do Adoptium
	}
	defer resp.Body.Close()

	var m jreManifest
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return adoptiumFallbackURL(), nil
	}

	switch runtime.GOOS + "/" + runtime.GOARCH {
	case "windows/amd64":
		return m.JRE.WindowsX64, nil
	case "linux/amd64":
		return m.JRE.LinuxX64, nil
	case "darwin/amd64", "darwin/arm64":
		return m.JRE.MacX64, nil
	default:
		return "", fmt.Errorf("plataforma %s/%s não suportada para download automático do JRE", runtime.GOOS, runtime.GOARCH)
	}
}

// adoptiumFallbackURL retorna a URL direta da API do Adoptium caso o release.json não esteja acessível.
func adoptiumFallbackURL() string {
	switch runtime.GOOS {
	case "windows":
		return "https://api.adoptium.net/v3/binary/latest/21/ga/windows/x64/jre/hotspot/normal/eclipse"
	case "darwin":
		return "https://api.adoptium.net/v3/binary/latest/21/ga/mac/x64/jre/hotspot/normal/eclipse"
	default:
		return "https://api.adoptium.net/v3/binary/latest/21/ga/linux/x64/jre/hotspot/normal/eclipse"
	}
}

// downloadToTemp faz download de url para um arquivo temporário e retorna seu caminho.
func downloadToTemp(url string) (string, error) {
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("servidor retornou status %d", resp.StatusCode)
	}

	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}
	tmp, err := os.CreateTemp("", "jre-*"+ext)
	if err != nil {
		return "", err
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, resp.Body); err != nil {
		os.Remove(tmp.Name())
		return "", err
	}
	return tmp.Name(), nil
}

// extractTarGz extrai um arquivo .tar.gz em destDir.
func extractTarGz(src, destDir string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(destDir, filepath.FromSlash(hdr.Name))
		// Proteção contra path traversal
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(target, 0755)
		case tar.TypeReg:
			os.MkdirAll(filepath.Dir(target), 0755)
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, hdr.FileInfo().Mode())
			if err != nil {
				return err
			}
			io.Copy(out, tr)
			out.Close()
		}
	}
	return nil
}

// extractZip extrai um arquivo .zip em destDir.
func extractZip(src, destDir string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		target := filepath.Join(destDir, filepath.FromSlash(f.Name))
		if !strings.HasPrefix(target, filepath.Clean(destDir)+string(os.PathSeparator)) {
			continue
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(target, 0755)
			continue
		}

		os.MkdirAll(filepath.Dir(target), 0755)
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			out.Close()
			return err
		}
		io.Copy(out, rc)
		rc.Close()
		out.Close()
	}
	return nil
}

// findJavaBinary localiza o executável java dentro de um diretório extraído do Temurin.
// O Temurin extrai em subdiretório como jdk-21.0.x+y-jre/bin/java.
func findJavaBinary(baseDir string) (string, error) {
	bin := "java"
	if runtime.GOOS == "windows" {
		bin = "java.exe"
	}

	var found string
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Base(path) == bin && strings.Contains(path, "bin") {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil || found == "" {
		return "", fmt.Errorf("executável java não encontrado em %s", baseDir)
	}
	return found, nil
}

// managedJavaPath retorna o caminho para o java gerenciado em ~/.hubsaude/jdk/.
func managedJavaPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	bin := "java"
	if runtime.GOOS == "windows" {
		bin = "java.exe"
	}

	// Busca direto em bin/ ou em subdiretório do Temurin
	jdkDir := filepath.Join(home, ".hubsaude", jdkSubDir)
	direct := filepath.Join(jdkDir, "bin", bin)
	if _, err := os.Stat(direct); err == nil {
		return direct, nil
	}

	return findJavaBinary(jdkDir)
}

// isJavaValid executa `java -version` e verifica se retorna sem erro.
func isJavaValid(javaPath string) bool {
	cmd := exec.Command(javaPath, "-version")
	return cmd.Run() == nil
}

// HubSaudeDir retorna o caminho do diretório de trabalho do Runner (~/.hubsaude/).
func HubSaudeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("não foi possível determinar o diretório home: %w", err)
	}
	dir := filepath.Join(home, ".hubsaude")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("não foi possível criar %s: %w", dir, err)
	}
	return dir, nil
}
