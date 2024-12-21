package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/officiallysidsingh/ecom-server/internal/store"
	"github.com/officiallysidsingh/ecom-server/internal/utils"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := store.GetAllProductsFromDB()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching products")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}

func GetProductById(w http.ResponseWriter, r *http.Request) {
	// Get Product ID from URL
	productID := chi.URLParam(r, "id")

	product, err := store.GetProductByIdFromDB(productID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product Not Found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, product)
}
