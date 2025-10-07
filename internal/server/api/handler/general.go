package handler

import (
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
)

func HandleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" || r.Method != http.MethodGet {
			errorResponse(w, r, http.StatusNotFound, "path not found")
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func HealthCheck(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"status":     http.StatusOK,
			"version":    cfg.AppVersion,
			"enviroment": cfg.AppEnv,
			"message":    "service is healthy",
		}
		respondWithJSON(w, r, http.StatusOK, data)
	}
}

func respondWithJSON[T any](w http.ResponseWriter, r *http.Request, status int, data T) {

	err := encoding.EncodeJson(w, r, status, data)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}
