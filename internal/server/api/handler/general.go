package handler

import (
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
)

func HandleRoot() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"message": "welcome to the tunnel nat traversal"
		}
		respondWithJSON(w, r, http.StatusOK, data)
	}
}

func HealthCheck(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"status":     http.StatusOK,
			"version":    cfg.AppVersion,
			"enviroment": cfg.AppEnviroment,
			"message":    "service is healthy",
		}
		respondWithJSON(w, r, http.StatusOK, data)
	}
}

func respondWithJSON[T any](w http.ResponseWriter, r *http.Request, status int, data T) {

	err := encoding.EncodeJson(w, r, status, data)
	if err != nil {
		serverErrorResponse(w, r, err)
	}
}
