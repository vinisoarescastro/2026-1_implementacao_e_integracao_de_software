package cmd

import (
	"fmt"
	"os"

	"github.com/kyriosdata/runner/internal/invoker"
	"github.com/kyriosdata/runner/internal/jdk"
	"github.com/spf13/cobra"
)

var (
	validateContent   string
	validateSignature string
	validateLocal     bool
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Valida uma assinatura digital simulada",
	Long: `Envia uma requisição de validação ao assinador.jar.
Por padrão usa o modo servidor (HTTP). Use --local para invocação direta via java -jar.`,
	Example: `  # Validar via servidor HTTP (padrão — requer assinatura start)
  assinatura validate --content "contrato.pdf" --signature "MOCKED_SIGNATURE_BASE64_=="

  # Validar invocando o JAR diretamente (sem servidor)
  assinatura validate --content "contrato.pdf" --signature "MOCKED_SIGNATURE_BASE64_==" --local`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if validateContent == "" {
			return fmt.Errorf("o parâmetro --content é obrigatório")
		}
		if validateSignature == "" {
			return fmt.Errorf("o parâmetro --signature é obrigatório")
		}

		javaPath, err := jdk.Resolve()
		if err != nil {
			return fmt.Errorf("não foi possível localizar ou provisionar o JDK: %w", err)
		}

		result, err := invoker.Validate(javaPath, validateContent, validateSignature, validateLocal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao validar assinatura: %v\n", err)
			os.Exit(1)
		}

		if result.Valid {
			fmt.Println("✔ Assinatura VÁLIDA")
		} else {
			fmt.Println("✘ Assinatura INVÁLIDA")
		}
		fmt.Printf("  Mensagem : %s\n", result.Message)
		return nil
	},
}

func init() {
	validateCmd.Flags().StringVar(&validateContent, "content", "", "Conteúdo original que foi assinado (obrigatório)")
	validateCmd.Flags().StringVar(&validateSignature, "signature", "", "Assinatura digital em Base64 (obrigatório)")
	validateCmd.Flags().BoolVar(&validateLocal, "local", false, "Invocar o assinador.jar diretamente (modo local/CLI)")
	_ = validateCmd.MarkFlagRequired("content")
	_ = validateCmd.MarkFlagRequired("signature")
	rootCmd.AddCommand(validateCmd)
}
