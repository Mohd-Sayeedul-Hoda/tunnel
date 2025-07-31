package handler

import (
	"fmt"
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
)

func HandleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data map[string]any
		if r.URL.Path != "/" || r.Method != "GET" {
			data = map[string]any{
				"message": fmt.Sprintf("%s path not found", r.URL.Path),
			}
			respondWithJSON(w, r, http.StatusNotFound, data)
			return
		}

		data = map[string]any{
			"message": "welcome to the tunnel nat traversal",
		}
		respondWithJSON(w, r, http.StatusOK, data)
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
