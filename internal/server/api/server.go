package api

import (
	"net/http"

	m "github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/middleware"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
)

func NewHTTPServer(cfg *config.Config, userRepo repositories.UserRepo) http.Handler {

	mux := http.NewServeMux()
	AddRoute(mux, cfg, userRepo)

	var handler http.Handler = mux
	handler = m.RecoverPanic(m.NewLoggingMiddleware(handler))

	return handler
}
