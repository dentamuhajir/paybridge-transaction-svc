package main

import (
	"log"
	"paybridge-transaction-service/internal/app"
	"paybridge-transaction-service/internal/config"
)

// @title           Paybridge Transaction Service
// @version         1.0
// @description     Internal transaction & account service
// @BasePath        /internal

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Example "Bearer b32fd54c1593e77fb4a84"

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
