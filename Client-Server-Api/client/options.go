package client

import "time"

type ClientOptions struct {
	ExternalRequestUrl string
	MaxTimeout         time.Duration
}

func NewDefaultClientOptions() *ClientOptions {
	return &ClientOptions{
		// O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
		ExternalRequestUrl: "http://localhost:8080/cotacao",
		// o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.
		MaxTimeout: 300 * time.Millisecond,
	}
}
