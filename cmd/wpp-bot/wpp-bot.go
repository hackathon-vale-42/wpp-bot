package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/hackathon-vale-42/wpp-bot/api"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	host, found := os.LookupEnv("HOST")
	if !found {
		host = "0.0.0.0"
	}

	port, found := os.LookupEnv("PORT")
	if !found {
		port = "8000"
	}

	server := api.NewServer()
	if server == nil {
		slog.Error("Couldn't create server")
		return
	}

	if err := server.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		slog.Error("Couldn't run server", "error", err)
	}
}
