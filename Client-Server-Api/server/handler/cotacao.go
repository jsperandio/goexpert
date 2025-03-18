package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type CotacaoDB interface {
	SaveWithContext(ctx context.Context, lastValue string) error
}

type CotacaoHandler struct {
	DB      CotacaoDB
	Options *CotacaoHandlerOptions
}

func NewCotacaoHandler(cdb CotacaoDB, opt *CotacaoHandlerOptions) *CotacaoHandler {
	if opt == nil {
		opt = NewDefaultCotacaoHandlerOptions()
	}

	return &CotacaoHandler{
		DB:      cdb,
		Options: opt,
	}
}

func (ch *CotacaoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(ctx, ch.Options.MaxTimeoutExternalRequest)
	defer cancel()
	c, err := ch.getCotacao(ctx)
	if err != nil {
		slog.Error("failed to get cotacao", "error", err)
		ch.handleErrorResponse(w, err)
		return
	}

	if c == nil {
		slog.Error("failed to get cotacao", "error", "cotacao is nil")
		ch.handleErrorResponse(w, errors.New("cotacao is nil"))
		return
	}

	dbCtx, cancel := context.WithTimeout(ctx, ch.Options.MaxTimeoutDBRequest)
	defer cancel()
	if err := ch.DB.SaveWithContext(dbCtx, c.UsdBrl.Bid); err != nil {
		slog.Error("failed to save last cotacao", "error", err)
		if errors.Is(err, context.DeadlineExceeded) {
			err = ErrDatabaseTimeout
		}
		ch.handleErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(c)
	if err != nil {
		slog.Error("failed to encode response on getCotacao", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	slog.Debug("Handler cotacao success response", "cotacao", c)
}

func (c *CotacaoHandler) getCotacao(ctx context.Context) (*CotacaoResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.Options.ExternalRequestUrl, nil)
	if err != nil {
		slog.Error("failed to create request on getCotacao", "error", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("failed to do request on getCotacao", "url", c.Options.ExternalRequestUrl, "error", err)
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, ErrExternalRequestTimeout
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("failed to get cotacao", "status_code", resp.StatusCode)
		return nil, ErrExternalRequestReturnedError
	}

	var cr CotacaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cr); err != nil {
		slog.Error("failed to decode response on getCotacao", "error", err)
		return nil, err
	}

	return &cr, nil
}

func (c *CotacaoHandler) handleErrorResponse(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrExternalRequestTimeout):
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusGatewayTimeout,
		})
		return

	case errors.Is(err, ErrDatabaseTimeout):
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusGatewayTimeout,
		})
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Error: err.Error(),
			Code:  http.StatusInternalServerError,
		})
		return
	}
}
