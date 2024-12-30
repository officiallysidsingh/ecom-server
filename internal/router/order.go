package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/officiallysidsingh/ecom-server/internal/config"
	"github.com/officiallysidsingh/ecom-server/internal/handlers"
	"github.com/officiallysidsingh/ecom-server/internal/middlewares"
	"github.com/officiallysidsingh/ecom-server/internal/services"
	"github.com/officiallysidsingh/ecom-server/internal/store"
)

func orderRoutes(db *sqlx.DB, envConfig *config.EnvConfig) chi.Router {
	// Initialize dependencies
	orderStore := store.NewOrderStore(db)
	orderService := services.NewOrderService(orderStore)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Set up router
	r := chi.NewRouter()

	// JWT Auth Validation Middleware
	r.Use(middlewares.ValidateJWT(db, envConfig))

	// Routes
	r.Get("/", orderHandler.GetAllOrders)
	r.Get("/{id}", orderHandler.GetOrderById)

	return r
}
