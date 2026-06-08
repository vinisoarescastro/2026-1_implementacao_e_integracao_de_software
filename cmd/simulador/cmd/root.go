package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "simulador",
	Short: "CLI para gerenciar o Simulador do HubSaúde",
	Long: `Sistema Runner — CLI para iniciar, parar e monitorar o Simulador do HubSaúde.
Gerencia o ciclo de vida do simulador.jar sem necessidade de configurar o ambiente Java manualmente.`,
}

// Execute é o ponto de entrada chamado pelo main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
