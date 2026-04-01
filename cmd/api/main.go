package main

import (
	"log"

	"github.com/AdityaAWP/IdiomaMate/internal/config"
	"github.com/AdityaAWP/IdiomaMate/internal/server"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	app := server.NewServer(cfg)
	if err := app.Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
