package handler

import (
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
)

func GetUsers(userRepo repositories.UserRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"working": "fine",
		}
		respondWithJSON(w, r, http.StatusOK, response)
	}
}
