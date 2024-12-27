package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/officiallysidsingh/ecom-server/internal/models"
	"github.com/officiallysidsingh/ecom-server/internal/services"
	"github.com/officiallysidsingh/ecom-server/internal/utils"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest

	// Decode User from Request Body to Struct
	errMessage, statusCode, err := utils.ParseBodyToJSON(w, r, &loginReq)
	if err != nil {
		log.Printf("Error decoding user data: %v", err)
		utils.RespondWithError(w, statusCode, errMessage)
		return
	}

	user, err := h.service.Login(r.Context(), &loginReq)
	if err != nil {
		log.Printf("Error logging in user: %v", err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Returning successful response
	res := fmt.Sprintf("User with Email: %s logged in successfully", user.Email)
	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": res})
}
