// Package invoker é responsável por invocar o assinador.jar,
// tanto no modo local (java -jar) quanto no modo servidor (HTTP).
package invoker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kyriosdata/runner/internal/jdk"
	"github.com/kyriosdata/runner/internal/release"
)

// Result representa a resposta retornada pelo assinador.jar.
type Result struct {
	Signature string `json:"signature"`
	Valid     bool   `json:"valid"`
	Message   string `json:"message"`
}

// jarPath retorna o caminho para o assinador.jar.
// Tenta provisionar automaticamente via release.json antes de buscar localmente.
func jarPath() (string, error) {
	hubsaudeDir, err := jdk.HubSaudeDir()
	if err == nil {
		// Tenta garantir jar atualizado via download automático
		if path, err := release.EnsureJar(hubsaudeDir); err == nil {
			return path, nil
		}

		// Fallback: jar já presente sem checar versão
		candidate := filepath.Join(hubsaudeDir, "assinador.jar")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	// Fallback: diretório corrente
	if _, err := os.Stat("assinador.jar"); err == nil {
		return "assinador.jar", nil
	}

	return "", fmt.Errorf(
		"assinador.jar não encontrado.\n" +
			"  O Runner tentará baixá-lo automaticamente na próxima execução se houver conexão com a internet.\n" +
			"  Alternativamente, compile com 'mvn -f assinador/pom.xml package' e copie para ~/.hubsaude/",
	)
}

// Sign invoca o assinador.jar para criar uma assinatura simulada.
// Por padrão usa o modo servidor (HTTP); --local força invocação direta via java -jar.
func Sign(javaPath, content, token string, local bool) (*Result, error) {
	if local {
		return runJar(javaPath, "sign", content, token, "")
	}
	return signViaHTTP(content, token)
}

// Validate invoca o assinador.jar para validar uma assinatura simulada.
// Por padrão usa o modo servidor (HTTP); --local força invocação direta via java -jar.
func Validate(javaPath, content, signature string, local bool) (*Result, error) {
	if local {
		return runJar(javaPath, "validate", content, "", signature)
	}
	return validateViaHTTP(content, signature)
}

// runJar executa o assinador.jar com os argumentos fornecidos e interpreta
// a saída JSON como um Result.
//
// Convenção de argumentos esperada pelo Main.java:
//
//	java -jar assinador.jar <operação> <content> [token] [signature]
func runJar(javaPath, operation, content, token, signature string) (*Result, error) {
	jar, err := jarPath()
	if err != nil {
		return nil, err
	}

	args := []string{
		"-Dfile.encoding=UTF-8",
		"-Dstdout.encoding=UTF-8",
		"-jar", jar, operation, content,
	}
	if token != "" {
		args = append(args, token)
	}
	if signature != "" {
		args = append(args, signature)
	}

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(javaPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	runErr := cmd.Run()

	// Tenta interpretar stdout como JSON antes de checar o exit code.
	// O JAR sempre emite JSON válido; exit code != 0 indica resultado de negócio (inválido),
	// não necessariamente falha de sistema.
	var result Result
	if jsonErr := json.Unmarshal(stdout.Bytes(), &result); jsonErr != nil {
		// Sem JSON válido: erro de sistema (JVM não encontrada, JAR corrompido, etc.)
		if runErr != nil && stderr.Len() > 0 {
			return nil, fmt.Errorf("assinador.jar falhou: %s", stderr.String())
		}
		return nil, fmt.Errorf("resposta inválida do assinador.jar: %w\nSaída: %s", jsonErr, stdout.String())
	}

	// JSON obtido com sucesso — retorna o resultado sem erro, mesmo quando valid=false.
	// O chamador decide como apresentar resultado inválido ao usuário.
	_ = runErr
	return &result, nil
}
