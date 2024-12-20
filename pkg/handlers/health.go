package handlers

import (
	"net/http"

	"github.com/officiallysidsingh/ecom-server/pkg/utils"
)

func Health(w http.ResponseWriter, r *http.Request) {
	utils.RespondWithJSON(w, http.StatusOK, "Server is Live and Running!!")
}
