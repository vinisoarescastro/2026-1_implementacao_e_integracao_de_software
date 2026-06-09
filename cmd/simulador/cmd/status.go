package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Exibe o status atual do Simulador",
	Long:  `Exibe se o simulador.jar está em execução e informações do processo (PID, porta).`,
	Example: `  # Ver se o Simulador está rodando
  simulador status`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("simulador status — a ser implementado na Sprint 4 (US-03.2)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
