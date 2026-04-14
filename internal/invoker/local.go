// Package invoker é responsável por invocar o assinador.jar,
// tanto no modo local (java -jar) quanto no modo servidor (HTTP).
package invoker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Result representa a resposta retornada pelo assinador.jar.
type Result struct {
	Signature string `json:"signature"`
	Valid     bool   `json:"valid"`
	Message   string `json:"message"`
}

// jarPath retorna o caminho para o assinador.jar.
// Procura primeiro em ~/.hubsaude/, depois no diretório corrente.
func jarPath() (string, error) {
	home, err := os.UserHomeDir()
	if err == nil {
		candidate := filepath.Join(home, ".hubsaude", "assinador.jar")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}

	// Fallback: diretório corrente
	local := "assinador.jar"
	if _, err := os.Stat(local); err == nil {
		return local, nil
	}

	return "", fmt.Errorf(
		"assinador.jar não encontrado em ~/.hubsaude/ nem no diretório corrente.\n" +
			"Dica: execute o build do projeto Java com 'mvn package' e copie o jar para ~/.hubsaude/",
	)
}

// Sign invoca o assinador.jar para criar uma assinatura simulada.
// Se local=true usa java -jar diretamente; caso contrário usa HTTP (Sprint 3).
func Sign(javaPath, content, token string, local bool) (*Result, error) {
	if !local {
		// HTTP mode será implementado na Sprint 3
		return nil, fmt.Errorf("modo servidor HTTP ainda não implementado — use --local por enquanto")
	}
	return runJar(javaPath, "sign", content, token, "")
}

// Validate invoca o assinador.jar para validar uma assinatura simulada.
// Se local=true usa java -jar diretamente; caso contrário usa HTTP (Sprint 3).
func Validate(javaPath, content, signature string, local bool) (*Result, error) {
	if !local {
		return nil, fmt.Errorf("modo servidor HTTP ainda não implementado — use --local por enquanto")
	}
	return runJar(javaPath, "validate", content, "", signature)
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

	args := []string{"-jar", jar, operation, content}
	if token != "" {
		args = append(args, token)
	}
	if signature != "" {
		args = append(args, signature)
	}

	cmd := exec.Command(javaPath, args...)
	out, err := cmd.Output()
	if err != nil {
		// Captura stderr para mensagem de erro mais clara
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("assinador.jar retornou erro:\n%s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("falha ao executar assinador.jar: %w", err)
	}

	var result Result
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("resposta inválida do assinador.jar: %w\nSaída bruta: %s", err, out)
	}

	if !result.Valid && result.Message != "" {
		return &result, fmt.Errorf("%s", result.Message)
	}

	return &result, nil
}
