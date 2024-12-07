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

	host, ok := os.LookupEnv("HOST")
	if ok != true {
		host = "0.0.0.0"
	}

	port, ok := os.LookupEnv("PORT")
	if ok != true {
		port = "8000"
	}

	server := api.NewServer()

	if err := server.Run(fmt.Sprintf("%s:%s", host, port)); err != nil {
		slog.Error("Server run error", "errorKind", err)
		panic(err)
	}
}
