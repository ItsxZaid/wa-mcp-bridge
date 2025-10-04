package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wa-mcp-bridge/internal/config"
	"wa-mcp-bridge/internal/server"
	"wa-mcp-bridge/internal/store"
	"wa-mcp-bridge/internal/whatsapp"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("INFO: No .env file found, reading from system environment")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dataStore, err := store.New()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer dataStore.Close()

	waBot, err := whatsapp.New(dataStore)
	if err != nil {
		log.Fatalf("Error initializing whatsapp: %v", err)
	}

	// TODO: before server we need to make whatsapp instance

	srv := server.New(cfg.HTTPPort, dataStore, waBot)

	go func() {
		log.Printf("Server starting on port: %s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server crashed reason: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Shutting down server...")
}
