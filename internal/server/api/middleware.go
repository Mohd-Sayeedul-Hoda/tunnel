package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
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

func newAuthenticateAndVerifyMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return authenticate(cfg, requireVerifiedUser(next))
	}
}

func requireVerifiedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := tools.ContextGetToken(r)
		if !token.Verified {
			handler.InactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func authenticate(cfg *config.Config, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var accessToken string

		if r.Header.Get("Authorization") != "" {
			headerParts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				handler.InvalidCredentialsResponse(w, r)
				return
			}
			accessToken = headerParts[1]
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

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:5173",
			"http://localhost:4173",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:5173",
			"http://127.0.0.1:4173",
		}

		allowed := slices.Contains(allowedOrigins, origin) || origin == ""

		if allowed {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			// Add CORS headers for all requests (including OPTIONS)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
		} else if origin != "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
