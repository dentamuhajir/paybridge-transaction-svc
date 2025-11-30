package main

import (
	"log"
	"paybridge-transaction-service/internal/config"
	"paybridge-transaction-service/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	if err := server.Run(cfg); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
