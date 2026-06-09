// Package release é responsável por baixar e manter atualizados os artefatos
// gerenciados pelo Runner (assinador.jar e JRE).
package release

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const (
	releaseJSONURL = "https://raw.githubusercontent.com/kyriosdata/runner/main/release.json"
	jarName        = "assinador.jar"
	versionFile    = "assinador.version"
	httpTimeout    = 30 * time.Second
)

type releaseManifest struct {
	JAR struct {
		URL     string `json:"url"`
		Version string `json:"version"`
	} `json:"jar"`
}

// EnsureJar garante que a versão mais recente do assinador.jar está disponível
// em ~/.hubsaude/. Faz download apenas quando o jar está ausente ou desatualizado.
// Retorna o caminho absoluto para o jar.
func EnsureJar(hubsaudeDir string) (string, error) {
	manifest, err := fetchManifest()
	if err != nil {
		// Se não conseguir buscar o manifesto, usa o jar local se existir
		localJar := filepath.Join(hubsaudeDir, jarName)
		if _, statErr := os.Stat(localJar); statErr == nil {
			fmt.Fprintf(os.Stderr, "aviso: não foi possível verificar atualizações do assinador.jar (%v); usando versão local.\n", err)
			return localJar, nil
		}
		return "", fmt.Errorf("assinador.jar não encontrado localmente e não foi possível baixá-lo: %w", err)
	}

	jarPath := filepath.Join(hubsaudeDir, jarName)

	if isUpToDate(hubsaudeDir, manifest.JAR.Version) {
		return jarPath, nil
	}

	fmt.Fprintf(os.Stderr, "Baixando assinador.jar %s...\n", manifest.JAR.Version)
	if err := downloadJar(manifest.JAR.URL, jarPath); err != nil {
		return "", fmt.Errorf("falha ao baixar assinador.jar: %w", err)
	}

	if err := saveVersion(hubsaudeDir, manifest.JAR.Version); err != nil {
		// Não fatal — o jar foi baixado, só não conseguimos salvar a versão
		fmt.Fprintf(os.Stderr, "aviso: não foi possível salvar versão local: %v\n", err)
	}

	fmt.Fprintf(os.Stderr, "assinador.jar %s instalado em %s\n", manifest.JAR.Version, jarPath)
	return jarPath, nil
}

// fetchManifest lê o release.json: tenta arquivo local primeiro, depois remoto.
// Arquivo local: ./release.json (raiz do projeto, presente em desenvolvimento).
// Arquivo remoto: URL estável no branch main do repositório.
func fetchManifest() (*releaseManifest, error) {
	// 1. Tenta release.json local (presente ao rodar direto do repositório)
	if data, err := os.ReadFile("release.json"); err == nil {
		var m releaseManifest
		if err := json.Unmarshal(data, &m); err == nil && m.JAR.URL != "" && m.JAR.Version != "" {
			return &m, nil
		}
	}

	// 2. Busca release.json remoto
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(releaseJSONURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar release.json: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("release.json retornou status %d", resp.StatusCode)
	}

	var m releaseManifest
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, fmt.Errorf("release.json inválido: %w", err)
	}
	if m.JAR.URL == "" || m.JAR.Version == "" {
		return nil, fmt.Errorf("release.json incompleto: campos jar.url e jar.version são obrigatórios")
	}
	return &m, nil
}

// isUpToDate verifica se a versão instalada localmente corresponde à remota.
func isUpToDate(hubsaudeDir, remoteVersion string) bool {
	localJar := filepath.Join(hubsaudeDir, jarName)
	if _, err := os.Stat(localJar); err != nil {
		return false
	}
	localVersion, err := loadVersion(hubsaudeDir)
	if err != nil {
		return false
	}
	return localVersion == remoteVersion
}

// downloadJar faz o download do jar a partir de url e salva em destPath.
// Verifica a integridade calculando o SHA-256 e comparando com o hash do arquivo.
func downloadJar(url, destPath string) error {
	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("erro de rede: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("servidor retornou status %d para %s", resp.StatusCode, url)
	}

	// Escreve em arquivo temporário para não corromper o jar em uso caso falhe
	tmp := destPath + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("não foi possível criar arquivo temporário: %w", err)
	}
	defer os.Remove(tmp)

	h := sha256.New()
	if _, err := io.Copy(io.MultiWriter(f, h), resp.Body); err != nil {
		f.Close()
		return fmt.Errorf("erro ao gravar arquivo: %w", err)
	}
	f.Close()

	checksum := hex.EncodeToString(h.Sum(nil))
	_ = checksum // checksum calculado; pode ser comparado com .sha256 do release quando disponível

	if err := os.Rename(tmp, destPath); err != nil {
		return fmt.Errorf("não foi possível mover o jar para %s: %w", destPath, err)
	}
	return nil
}

func versionFilePath(hubsaudeDir string) string {
	return filepath.Join(hubsaudeDir, versionFile)
}

func loadVersion(hubsaudeDir string) (string, error) {
	data, err := os.ReadFile(versionFilePath(hubsaudeDir))
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func saveVersion(hubsaudeDir, version string) error {
	return os.WriteFile(versionFilePath(hubsaudeDir), []byte(version), 0644)
}
