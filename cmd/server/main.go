package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/officiallysidsingh/ecom-server/pkg/handlers"
)

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)

	// Health Routes
	r.Get("/", handlers.Health)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}
