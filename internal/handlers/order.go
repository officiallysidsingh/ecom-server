package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/officiallysidsingh/ecom-server/internal/services"
	"github.com/officiallysidsingh/ecom-server/internal/utils"
)

type OrderHandler struct {
	service services.OrderService
}

func NewOrderHandler(service services.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.GetAll(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderById(w http.ResponseWriter, r *http.Request) {
	orderID := chi.URLParam(r, "id")

	order, err := h.service.GetByID(r.Context(), orderID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, order)
}
