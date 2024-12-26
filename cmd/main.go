package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/officiallysidsingh/ecom-server/db"
	"github.com/officiallysidsingh/ecom-server/internal/config"
	"github.com/officiallysidsingh/ecom-server/internal/router"
)

func main() {
	// Load environment variables from .env file (if available)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found or error loading it: %v", err)
	}

	// Init Config
	appConfig := config.LoadConfig()

	// Init DB connection
	dbConn, err := db.InitDB(appConfig.DATABASE_URL)
	if err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}

	// Close DB connection on shutdown
	defer db.CloseDB(dbConn)

	// Setup Router & Middlewares
	r := router.Setup(dbConn)

	// Start server with graceful shutdown
	startServerWithGracefulShutdown(r, appConfig.SERVER_PORT)
}

func startServerWithGracefulShutdown(handler http.Handler, port string) {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Starting server on port %s...\n", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// Wait for termination signals (graceful shutdown)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Println("Shutting down gracefully...")

	// Set a deadline to wait for any ongoing requests to finish before shutting down
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down server: %v", err)
	}

	log.Println("Server stopped gracefully")
}
