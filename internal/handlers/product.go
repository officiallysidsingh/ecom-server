package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/officiallysidsingh/ecom-server/internal/models"
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

func AddProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Printf("Error decoding product data: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product data")
		return
	}

	// Call the function to add the product to the database
	productID, err := store.AddProductToDB(&product)
	if err != nil {
		log.Printf("Error adding product: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add product")
		return
	}

	// Respond with the new product ID
	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"product_id": productID})
}
