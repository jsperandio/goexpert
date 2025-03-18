package client

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"text/template"
)

type Client struct {
	Options *ClientOptions
}

func NewClient(opt *ClientOptions) *Client {
	if opt == nil {
		opt = NewDefaultClientOptions()
	}

	return &Client{
		Options: opt,
	}
}

func (c *Client) GetCotacaoDolar() (*Cotacao, error) {
	tCtx, cancel := context.WithTimeout(context.Background(), c.Options.MaxTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(tCtx, http.MethodGet, c.Options.ExternalRequestUrl, nil)
	if err != nil {
		slog.Error("failed to create request", "error", err)
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("failed to do request", "error", err)

		if errors.Is(err, context.DeadlineExceeded) {
			err = ErrExternalRequestTimeout
		}

		return nil, c.handleError(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("failed to get cotacao", "status_code", resp.StatusCode)

		errE := NewExternalRequestError(resp.StatusCode)
		return nil, c.handleError(errE)
	}

	var csr CotacaoServerResponse
	err = json.NewDecoder(resp.Body).Decode(&csr)
	if err != nil {
		slog.Error("failed to decode response", "error", err)
		return nil, err
	}

	return NewCotacao(csr.UsdBrl.Bid), nil
}

func (c *Client) handleError(err error) error {
	switch {
	case errors.Is(err, ErrExternalRequestTimeout):
		return ErrExternalRequestTimeout
	default:
		return err
	}
}

func (c *Client) SaveToFile(cotacao Cotacao) error {
	file, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		slog.Error("failed to open file", "error", err)
		return err
	}
	defer file.Close()

	tmpl := template.Must(template.New("cotacao").Parse("Dólar: {{.UsdBrl}}\n"))
	err = tmpl.Execute(file, cotacao)
	if err != nil {
		slog.Error("failed to execute template", "error", err)
		return err
	}

	return nil
}

func (c *Client) RunAndSave() {
	cot, err := c.GetCotacaoDolar()
	if err != nil {
		slog.Error("failed to get run client get cotacao", "error", err)
		os.Exit(1)
	}

	if cot == nil {
		slog.Error("failed to get run client get cotacao", "error", "cotacao is nil")
		os.Exit(1)
	}

	if err := c.SaveToFile(*cot); err != nil {
		slog.Error("failed to save to file", "error", err)
		os.Exit(1)
	}

	slog.Debug("success on get cotacao and save to file", "cotacao", cot)
}
