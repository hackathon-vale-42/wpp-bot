package main

import (
	"log"

	"github.com/hackathon-vale-42/wpp-bot/api"
)

func main() {
	server := api.NewServer()

	if err := server.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
