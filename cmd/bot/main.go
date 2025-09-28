package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"wa-mcp-bridge/internal/config"
	"wa-mcp-bridge/internal/server"
	"wa-mcp-bridge/internal/store"

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

	dataStore, err := store.New();
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer dataStore.Close()

	routes := server.New()

	port := fmt.Sprintf(":%s", cfg.HTTPPort)

	go http.ListenAndServe(port, routes)

	fmt.Println("Server is running successfully on port", cfg.HTTPPort)
	
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
}