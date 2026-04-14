// Package jdk é responsável por detectar ou provisionar automaticamente
// o JDK 21 necessário para executar o assinador.jar.
package jdk

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

const (
	jdkDir     = ".hubsaude/jdk"
	minVersion = "21"
)

// Resolve retorna o caminho para o executável java disponível.
// Ordem de busca:
//  1. ~/.hubsaude/jdk/bin/java  (provisionado pelo Runner)
//  2. java no PATH do sistema
//
// Se nenhum for encontrado, retorna um erro orientativo.
func Resolve() (string, error) {
	// 1. Verifica JDK provisionado localmente
	managed, err := managedJavaPath()
	if err == nil {
		if isJavaValid(managed) {
			return managed, nil
		}
	}

	// 2. Verifica java no PATH do sistema
	systemJava, err := exec.LookPath("java")
	if err == nil && isJavaValid(systemJava) {
		return systemJava, nil
	}

	return "", fmt.Errorf(
		"JDK %s não encontrado.\n"+
			"  • Instale o JDK manualmente (https://adoptium.net/) e adicione ao PATH, ou\n"+
			"  • Aguarde: o provisionamento automático será implementado na Sprint 2 (US-04.1).",
		minVersion,
	)
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

	path := filepath.Join(home, jdkDir, "bin", bin)
	if _, err := os.Stat(path); err != nil {
		return "", fmt.Errorf("JDK gerenciado não encontrado em %s", path)
	}
	return path, nil
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
