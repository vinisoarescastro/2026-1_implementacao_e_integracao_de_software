package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "CLI para assinatura digital simulada",
	Long: `Sistema Runner — CLI para criação e validação de assinaturas digitais.
Invoca o assinador.jar sem necessidade de configurar o ambiente Java manualmente.`,
}

// Execute é o ponto de entrada chamado pelo main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Version = version
}
