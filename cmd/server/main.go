package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/officiallysidsingh/ecom-server/database"
	"github.com/officiallysidsingh/ecom-server/pkg/config"
	"github.com/officiallysidsingh/ecom-server/pkg/handlers"
)

func main() {
	// Load environment variables from .env file (if available)
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found or error loading it: %v", err)
	}

	// Initialize the database connection with proper connection string from env
	connStr := config.GetEnv("DATABASE_URL", "postgres://postgres:sidsingh@localhost/dbname?sslmode=disable")
	if err := database.InitDB(connStr); err != nil {
		log.Fatalf("Error initializing the database: %v", err)
	}
	defer database.CloseDB() // Ensure the database connection is closed on shutdown

	// Set up the router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)

	// Health Routes
	r.Get("/", handlers.Health)

	// Product Routes
	r.Get("/products", handlers.GetProducts)

	// Handle graceful shutdown signals
	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Starting the server in a goroutine
	go func() {
		log.Println("Starting server on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting the server: %v", err)
		}
	}()

	// Wait for termination signals (graceful shutdown)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received
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
