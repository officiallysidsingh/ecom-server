package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/services"
	"github.com/officiallysidsingh/ecom-server/internal/utils"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.GetAll(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")

	product, err := h.service.GetByID(r.Context(), productID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, product)
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	// Decode Product from Request Body to Struct
	errMessage, statusCode, err := utils.ParseBodyToJSON(w, r, &product)
	if err != nil {
		log.Printf("Error decoding product data: %v", err)
		utils.RespondWithError(w, statusCode, errMessage)
		return
	}

	productID, err := h.service.Create(r.Context(), &product)
	if err != nil {
		log.Printf("Error adding product: %v", err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Returning successful response
	res := fmt.Sprintf("Product with id: %s added successfully", productID)
	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": res})
}

func (h *ProductHandler) PutUpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	// Get ProductID from URL
	productID := chi.URLParam(r, "id")

	// Decode Product from Request Body to Struct
	errMessage, statusCode, err := utils.ParseBodyToJSON(w, r, &product)
	if err != nil {
		log.Printf("Error decoding product data: %v", err)
		utils.RespondWithError(w, statusCode, errMessage)
		return
	}

	// Call the function to add the product in DB
	if err := h.service.PutUpdate(r.Context(), &product, productID); err != nil {
		log.Printf("Error updating product (ID: %s): %v", productID, err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Returning successful response
	res := fmt.Sprintf("Product with id: %s updated successfully", productID)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": res})
}

func (h *ProductHandler) PatchUpdateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product

	// Get ProductID from URL
	productID := chi.URLParam(r, "id")

	// Decode Product from JSON to Struct
	errMessage, statusCode, err := utils.ParseBodyToJSON(w, r, &product)
	if err != nil {
		log.Printf("Error decoding product data: %v", err)
		utils.RespondWithError(w, statusCode, errMessage)
		return
	}

	// Call the function to update the product in DB
	if err := h.service.PatchUpdate(r.Context(), &product, productID); err != nil {
		log.Printf("Error updating product (ID: %s): %v", productID, err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Returning successful response
	res := fmt.Sprintf("Product with id: %s updated successfully", productID)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": res})
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get ProudctID from URL
	productID := chi.URLParam(r, "id")

	// Call the function to delete the product in DB
	if err := h.service.Delete(r.Context(), productID); err != nil {
		log.Printf("Error updating product (ID: %s): %v", productID, err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Returning successful response
	res := fmt.Sprintf("Product with id: %s deleted successfully", productID)
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": res})
}
