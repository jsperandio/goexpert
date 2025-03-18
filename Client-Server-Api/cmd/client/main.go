package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jsperandio/goexpert/clientserverapi/client"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	)))

	cOpts := client.NewDefaultClientOptions()
	c := client.NewClient(cOpts)
	c.RunAndSave()

	fmt.Println("Client finished successfully")
}
