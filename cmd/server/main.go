package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"paybridge-transaction-service/internal/config"
	"paybridge-transaction-service/internal/db"
	"paybridge-transaction-service/internal/server"
	"strconv"
	"syscall"
	"time"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	pg, err := db.NewPostgres(cfg.Database.DSN)
	if err != nil {
		log.Fatalf("Postgres error: %v", err)
	}

	deps := &server.Dependencies{
		DB: pg,
	}

	e := server.NewRouter(deps)

	go func() {
		log.Printf("Server starting on port %d...", cfg.Server.Port)

		if err := e.Start(":" + strconv.Itoa(cfg.Server.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatalf(" Server failed: %v", err)
		}
	}()

	// Block main — wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Println("Shutdown signal received:", sig.String())

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Println("Forced shutdown:", err)
	} else {
		log.Println("Server shut down gracefully")
	}
}
