package cmd

import (
	"fmt"
	"os"

	"github.com/kyriosdata/runner/internal/invoker"
	"github.com/kyriosdata/runner/internal/jdk"
	"github.com/spf13/cobra"
)

var (
	signContent string
	signToken   string
	signLocal   bool
)

var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Cria uma assinatura digital simulada",
	Long: `Envia uma requisição de assinatura digital ao assinador.jar.
Por padrão usa o modo servidor (HTTP). Use --local para invocação direta via java -jar.`,
	Example: `  # Assinar via servidor HTTP (padrão — requer assinatura start)
  assinatura sign --content "contrato.pdf"

  # Assinar invocando o JAR diretamente (sem servidor)
  assinatura sign --content "contrato.pdf" --local

  # Assinar com token de dispositivo criptográfico
  assinatura sign --content "contrato.pdf" --token "meu-pin" --local`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if signContent == "" {
			return fmt.Errorf("o parâmetro --content é obrigatório")
		}

		// Garantir que o JDK está disponível antes de invocar o jar
		javaPath, err := jdk.Resolve()
		if err != nil {
			return fmt.Errorf("não foi possível localizar ou provisionar o JDK: %w", err)
		}

		result, err := invoker.Sign(javaPath, signContent, signToken, signLocal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao criar assinatura: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✔ Assinatura criada com sucesso")
		fmt.Printf("  Assinatura : %s\n", result.Signature)
		fmt.Printf("  Mensagem   : %s\n", result.Message)
		return nil
	},
}

func init() {
	signCmd.Flags().StringVar(&signContent, "content", "", "Conteúdo a ser assinado (obrigatório)")
	signCmd.Flags().StringVar(&signToken, "token", "", "Token de autenticação / PIN (opcional)")
	signCmd.Flags().BoolVar(&signLocal, "local", false, "Invocar o assinador.jar diretamente (modo local/CLI)")
	_ = signCmd.MarkFlagRequired("content")
	rootCmd.AddCommand(signCmd)
}
