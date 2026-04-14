package cmd_test

import (
	"os/exec"
	"strings"
	"testing"
)

// TestVersionCommand verifica que `go run . version` exibe a string "dev"
// (valor padrão da variável version sem injeção de ldflags).
func TestVersionCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "version")
	cmd.Dir = ".." // executa a partir de cmd/assinatura/

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("falha ao executar 'go run . version': %v\nSaída: %s", err, out)
	}

	output := strings.TrimSpace(string(out))
	if !strings.Contains(output, "dev") {
		t.Errorf("esperava 'dev' na saída, obteve: %q", output)
	}
}
