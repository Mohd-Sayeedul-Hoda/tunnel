package main

import (
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
)

func healthCheck(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"status":     http.StatusOK,
			"version":    cfg.AppVersion,
			"enviroment": cfg.AppEnviroment,
			"message":    "service is healthy",
		}
		encoding.Encode(w, r.Header, http.StatusOK, data)
	}
}
