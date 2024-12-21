package handlers

import (
	"net/http"

	"github.com/officiallysidsingh/ecom-server/internal/store"
	"github.com/officiallysidsingh/ecom-server/internal/utils"
)

func GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := store.GetAllProducts()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching products")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, products)
}
