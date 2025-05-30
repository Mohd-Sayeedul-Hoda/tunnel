package handler

import (
	"log/slog"
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
)

type envelope map[string]any

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := encoding.EncodeJson(w, r, status, env)
	if err != nil {
		slog.Error("server error", "err", err)
		w.WriteHeader(500)
	}
}

func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	slog.Error("server error",
		slog.String("err", err.Error()),
		slog.String("request", r.Method),
		slog.String("url", r.URL.Path),
	)

	message := "the server encounter a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message)
}
