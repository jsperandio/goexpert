package main

import (
	"log/slog"
	"os"

	"github.com/jsperandio/goexpert/clientserverapi/server"
	"github.com/jsperandio/goexpert/clientserverapi/server/database"
	"github.com/jsperandio/goexpert/clientserverapi/server/handler"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	)))

	database.InitDb()

	c := database.NewCotacao()

	chOpts := handler.NewDefaultCotacaoHandlerOptions()
	ch := handler.NewCotacaoHandler(c, chOpts)

	svrOpts := server.NewDefaultServerOptions()
	svr := server.NewServer(svrOpts)
	svr.RegisterRoute("GET /cotacao", ch)

	if err := svr.Start(); err != nil {
		panic(err)
	}
}
