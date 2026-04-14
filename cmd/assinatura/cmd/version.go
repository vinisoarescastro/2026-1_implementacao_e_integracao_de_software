package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// version é injetada em tempo de build via:
//
//	go build -ldflags "-X github.com/kyriosdata/runner/cmd/assinatura/cmd.version=v1.0.0"
//
// Deve ser variável (não constante) para que o linker consiga sobrescrever.
var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão atual do CLI",
	Long:  `Exibe a versão atual do CLI assinatura junto com o SO e a arquitetura.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("assinatura %s %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
