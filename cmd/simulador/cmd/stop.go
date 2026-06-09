package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Para o Simulador do HubSaúde",
	Long:  `Encerra o processo do simulador.jar em execução.`,
	Example: `  # Parar o Simulador
  simulador stop`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("simulador stop — a ser implementado na Sprint 4 (US-03.2)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
