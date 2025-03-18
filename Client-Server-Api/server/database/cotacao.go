package database

import (
	"context"
	"log/slog"
)

type Cotacao struct{}

func NewCotacao() *Cotacao {
	return &Cotacao{}
}

func (c *Cotacao) SaveWithContext(ctx context.Context, lastValue string) error {
	rst, err := DB.ExecContext(ctx, "INSERT INTO cotacao (valor, created_at) VALUES (?, datetime('now','localtime'))", lastValue)
	if err != nil {
		slog.Error("failed to save cotacao in database", "error", err)
		return err
	}

	lid, err := rst.LastInsertId()
	if err != nil {
		slog.Error("failed to get last insert id", "error", err)
		return err
	}

	slog.Debug("saved", "id", lid)
	return nil
}
