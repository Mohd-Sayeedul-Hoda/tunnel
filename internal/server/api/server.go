package api

import (
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/cache"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/shared/config"
)

func NewHTTPServer(cfg *config.Config, cacheRepo cache.CacheRepo, userRepo repositories.UserRepo, apiKeyRepo repositories.APIRepo, emailOtpRepo repositories.EmailOtpRepo) http.Handler {

	mux := http.NewServeMux()
	AddRoute(mux, cfg, cacheRepo, userRepo, apiKeyRepo, emailOtpRepo)

	var handler http.Handler = mux
	handler = CORS(RecoverPanic(NewLoggingMiddleware(handler)))

	return handler
}
