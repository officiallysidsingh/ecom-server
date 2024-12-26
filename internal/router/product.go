package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/officiallysidsingh/ecom-server/internal/handlers"
	"github.com/officiallysidsingh/ecom-server/internal/services"
	"github.com/officiallysidsingh/ecom-server/internal/store"
)

func productRoutes(db *sqlx.DB) chi.Router {
	// Initialize dependencies
	productStore := store.NewProductStore(db)
	productService := services.NewProductService(productStore)
	productHandler := handlers.NewProductHandler(productService)

	// Set up router
	r := chi.NewRouter()

	// Routes
	r.Get("/", productHandler.GetAllProducts)
	r.Get("/{id}", productHandler.GetProductById)
	r.Post("/", productHandler.AddProduct)
	r.Put("/{id}", productHandler.PutUpdateProduct)
	r.Patch("/{id}", productHandler.PatchUpdateProduct)
	r.Delete("/{id}", productHandler.DeleteProduct)

	return r
}
