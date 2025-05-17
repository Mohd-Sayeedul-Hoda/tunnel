package handler

import (
	"net/http"
)

func GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"working": "fine",
		}
		respondWithJSON(w, r, http.StatusOK, response)
	}
}
