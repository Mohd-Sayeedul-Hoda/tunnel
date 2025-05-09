package main

import (
	"log/slog"
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
)

type envelope map[string]any

func errorResponse(w http.ResponseWriter, status int, message any) {
	header := make(http.Header)
	env := envelope{"error": message}

	err := encoding.Encode(w, header, status, env)
	if err != nil {
		slog.Error("server error", "err", err)
		w.WriteHeader(500)
	}
}

func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	slog.Error("server error", "err", err, "request", r.Method, r.URL)

	message := "the server encounter a problem and could not process your request"
	errorResponse(w, http.StatusInternalServerError, message)
}
