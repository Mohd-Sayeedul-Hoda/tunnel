package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/handler"
	tools "github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/utils"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/utils"
)

type logRespWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *logRespWriter {
	return &logRespWriter{w, http.StatusOK}
}

func (lrw *logRespWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func NewLoggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lrw := NewLoggingResponseWriter(w)
		startTime := time.Now()

		next.ServeHTTP(lrw, r)

		duration := time.Since(startTime)
		slog.Info("api info",
			"status_code", lrw.statusCode,
			"status", http.StatusText(lrw.statusCode),
			"method", r.Method,
			"uri", r.RequestURI,
			"remoteAddr", r.RemoteAddr,
			"duration", duration.String(),
		)
	}
}

func RecoverPanic(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				handler.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	}
}

func newAuthenticateMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return authenticate(cfg, next)
	}
}

func authenticate(cfg *config.Config, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var accessToken string

		if r.Header.Get("Authorization") != "" {
			slog.Info(r.Header.Get("Authorization"))
			accessToken = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		} else {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				handler.InvalidCredentialsResponse(w, r)
				return
			}
			accessToken = cookie.Value
		}

		tokenDetail, err := utils.ValidateToken(accessToken, cfg.Token.AccessTokenPublicKey)
		if err != nil {
			switch {
			case errors.Is(err, utils.ErrTokenExpired):
				handler.TokenExpireResponse(w, r)
			case errors.Is(err, utils.ErrInvalidClaims):
				handler.InvalidCredentialsResponse(w, r)
			default:
				handler.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = tools.ContextSetToken(r, tokenDetail)
		next.ServeHTTP(w, r)
	})
}

func adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := tools.ContextGetToken(r)
		if !token.IsAdmin {
			handler.NotPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
