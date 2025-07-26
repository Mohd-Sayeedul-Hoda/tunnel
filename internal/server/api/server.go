package api

import (
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/cache"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
)

func NewHTTPServer(cfg *config.Config, cacheRepo cache.CacheRepo, userRepo repositories.UserRepo) http.Handler {

	mux := http.NewServeMux()
	AddRoute(mux, cfg, cacheRepo, userRepo)

	var handler http.Handler = mux
	handler = RecoverPanic(NewLoggingMiddleware(handler))

	return handler
}
