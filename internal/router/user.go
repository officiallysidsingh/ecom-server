package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/officiallysidsingh/ecom-server/internal/handlers"
	"github.com/officiallysidsingh/ecom-server/internal/services"
	"github.com/officiallysidsingh/ecom-server/internal/store"
)

func userRoutes(db *sqlx.DB) chi.Router {
	// Initialize dependencies
	userStore := store.NewUserStore(db)
	userService := services.NewUserService(userStore)
	userHandler := handlers.NewUserHandler(userService)

	// Setup a new router
	r := chi.NewRouter()

	// Routes
	r.Post("/login", userHandler.Login)

	return r
}
