package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"paybridge-transaction-service/internal/config"
	"paybridge-transaction-service/internal/db"
	"strconv"
	"syscall"
	"time"
)

func Run(cfg *config.Config) error {
	// Setup DB
	pg, err := db.NewPostgres(cfg.Database.DSN)
	if err != nil {
		return err
	}

	deps := &Dependencies{
		DB: pg,
	}

	// Setup router
	e := NewRouter(deps)

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %d...", cfg.Server.Port)

		if err := e.Start(":" + strconv.Itoa(cfg.Server.Port)); err != nil &&
			err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutdown initiated...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Println("Forced shutdown:", err)
	} else {
		log.Println("Server shut down gracefully")
	}

	return nil
}
