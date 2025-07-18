package middleware

import (
	"log/slog"
	"net/http"
	"time"
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
		slog.Info("api response info",
			"status_code", lrw.statusCode,
			"status", http.StatusText(lrw.statusCode),
			"method", r.Method,
			"uri", r.RequestURI,
			"remoteAddr", r.RemoteAddr,
			"duration", duration.String(),
		)
	}
}
