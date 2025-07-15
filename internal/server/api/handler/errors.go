package handler

import (
	"log/slog"
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/encoding"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/request"
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

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message)
}

func failedValidationResponse(w http.ResponseWriter, r *http.Request, errors *request.Valid) {
	errorResponse(w, r, http.StatusUnprocessableEntity, errors.Errors)
}
