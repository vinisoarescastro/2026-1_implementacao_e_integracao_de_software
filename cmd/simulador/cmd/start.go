package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startSource string

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Inicia o Simulador do HubSaúde",
	Long: `Inicia o simulador.jar. Se o jar não estiver disponível localmente,
faz o download automático da versão mais recente no GitHub Releases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("simulador start — a ser implementado na Sprint 4 (US-03.1)")
		return nil
	},
}

func init() {
	startCmd.Flags().StringVar(&startSource, "source", "", "URL alternativa para download do simulador.jar")
	rootCmd.AddCommand(startCmd)
}
