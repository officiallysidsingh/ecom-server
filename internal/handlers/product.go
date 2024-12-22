package handlers

import (
	"encoding/json"
	"fmt"
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

	// Decode Product from JSON to Struct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
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

	// Returning successful response
	res := fmt.Sprintf("Product with id: %s added successfully", productID)
	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": res})
}

func PutUpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	// Get Product ID from URL
	productID := chi.URLParam(r, "id")

	// Decode Product from JSON to Struct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		log.Printf("Error decoding product data: %v", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product data")
		return
	}

	// Call the function to add the product to the database
	if err := store.PutUpdateProductInDB(&product, productID); err != nil {
		log.Printf("Error updating product: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	// Returning successful response
	res := fmt.Sprintf("Product with id: %s updated successfully", productID)
	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": res})
}
