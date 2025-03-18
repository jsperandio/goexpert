package handler

import "time"

type CotacaoHandlerOptions struct {
	ExternalRequestUrl        string
	MaxTimeoutExternalRequest time.Duration
	MaxTimeoutDBRequest       time.Duration
}

func NewDefaultCotacaoHandlerOptions() *CotacaoHandlerOptions {
	return &CotacaoHandlerOptions{
		// https://economia.awesomeapi.com.br/json/last/USD-BRL
		ExternalRequestUrl: EconomiaAwesomeApiDolarURL,
		// o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms
		MaxTimeoutExternalRequest: 200 * time.Millisecond,
		// o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms.
		MaxTimeoutDBRequest: 10 * time.Millisecond,
	}
}
