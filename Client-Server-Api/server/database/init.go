package database

import (
	"database/sql"
	"log/slog"
)

var DB *sql.DB

func InitDb() {
	var err error
	DB, err = sql.Open("sqlite3", "cotacao.db")
	if err != nil {
		slog.Error("failed to open database", "error", err)
		panic(err)
	}

	if err := checkDb(); err != nil {
		slog.Error("failed in check database", "error", err)
		panic(err)
	}

	statment := `
		CREATE TABLE IF NOT EXISTS cotacao (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			valor TEXT NOT NULL,
			created_at DATETIME NOT NULL
		);
	`
	_, err = DB.Exec(statment)
	if err != nil {
		slog.Error("failed to create table", "error", err)
		panic(err)
	}
}

func checkDb() error {
	if err := DB.Ping(); err != nil {
		slog.Error("failed to ping database", "error", err)
		return err
	}

	return nil
}
