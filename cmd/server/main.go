package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"paybridge-transaction-service/internal/server"
	"syscall"
	"time"
)

func main() {
	e := server.NewRouter()

	go func() {
		log.Println("Server starting on port 8083...")

		if err := e.Start(":8083"); err != nil && err != http.ErrServerClosed {
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
