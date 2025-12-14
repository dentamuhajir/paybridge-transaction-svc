package main

import (
	"log"
	"paybridge-transaction-service/internal/app"
	"paybridge-transaction-service/internal/config"
)

// @title Paybridge Transaction Service API
// @version 1.0
// @description API documentation for Transaction Services
// @BasePath /api/v1
// @securityDefinitions.apikey InternalTokenAuth
// @in header
// @name Authorization
// @description Provide token as: Bearer <token>
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	bootstrap := app.NewBootstrap(cfg)

	if err := bootstrap.Start(); err != nil {
		log.Fatalf("Startup error: %v", err)
	}
}
