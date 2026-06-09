package invoker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	DefaultServerPort = 8080
	defaultServerURL  = "http://localhost:8080"
	httpTimeout       = 10 * time.Second
)

type signHTTPRequest struct {
	Content string `json:"content"`
	Token   string `json:"token,omitempty"`
}

type validateHTTPRequest struct {
	Content   string `json:"content"`
	Signature string `json:"signature"`
}

// signViaHTTP envia uma requisição POST /sign ao assinador.jar em modo servidor.
func signViaHTTP(content, token string) (*Result, error) {
	body, err := json.Marshal(signHTTPRequest{Content: content, Token: token})
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}
	return postJSON(defaultServerURL+"/sign", body)
}

// validateViaHTTP envia uma requisição POST /validate ao assinador.jar em modo servidor.
func validateViaHTTP(content, signature string) (*Result, error) {
	body, err := json.Marshal(validateHTTPRequest{Content: content, Signature: signature})
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %w", err)
	}
	return postJSON(defaultServerURL+"/validate", body)
}

func postJSON(url string, body []byte) (*Result, error) {
	client := &http.Client{Timeout: httpTimeout}

	resp, err := client.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf(
			"não foi possível conectar ao servidor do assinador (porta %d).\n"+
				"  Inicie o servidor com:  assinatura start\n"+
				"  Ou use o modo local:    assinatura sign --local ...",
			DefaultServerPort,
		)
	}
	defer resp.Body.Close()

	var result Result
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("resposta inválida do servidor: %w", err)
	}

	if !result.Valid && result.Message != "" {
		return &result, fmt.Errorf("%s", result.Message)
	}

	return &result, nil
}
