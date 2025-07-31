package api

import (
	"context"
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/utils"
)

type contextKey string

const tokenContextKey = contextKey("tokenDetails")

func ContextSetToken(r *http.Request, token *utils.TokenDetails) *http.Request {
	ctx := context.WithValue(r.Context(), tokenContextKey, token)
	return r.WithContext(ctx)
}

func ContextGetToken(r *http.Request) *utils.TokenDetails {
	token, ok := r.Context().Value(tokenContextKey).(*utils.TokenDetails)
	if !ok {
		panic("missing user value in request context")
	}

	return token
}
