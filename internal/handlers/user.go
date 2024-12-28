package handlers

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

	// Go to Login service
	user, token, err := h.service.Login(r.Context(), &loginReq)
	if err != nil {
		log.Printf("Error logging in user: %v", err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Set JWT in HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  time.Now().Add(15 * time.Minute),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	// Returning successful response
	utils.RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var signupReq models.SignupRequest

	// Decode User from Request Body to Struct
	errMessage, statusCode, err := utils.ParseBodyToJSON(w, r, &signupReq)
	if err != nil {
		log.Printf("Error decoding user data: %v", err)
		utils.RespondWithError(w, statusCode, errMessage)
		return
	}

	// Go to Signup service
	userID, token, err := h.service.Signup(r.Context(), &signupReq)
	if err != nil {
		log.Printf("Error signing up in user: %v", err.Error())
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Set JWT in HTTP-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		Expires:  time.Now().Add(15 * time.Minute),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	})

	// Returning successful response
	res := fmt.Sprintf("User with id: %s signed up successfully", userID)
	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": res})
}
