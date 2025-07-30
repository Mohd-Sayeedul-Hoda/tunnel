package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/handler"
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

		next.ServeHTTP(w, r)

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
