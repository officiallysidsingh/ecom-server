package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/officiallysidsingh/ecom-server/internal/handlers"
)

func Setup(db *sqlx.DB) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	setupMiddlewares(r)

	// Routes
	setupRoutes(db, r)

	return r
}

func setupMiddlewares(r *chi.Mux) {
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(15 * time.Second))
}

func setupRoutes(db *sqlx.DB, r *chi.Mux) {
	// Health Check
	r.Get("/", handlers.Health)

	// Sub-Routers
	r.Mount("/products", productRoutes(db))
	r.Mount("/user", userRoutes(db))
}
