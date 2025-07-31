package api

import (
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/handler"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/cache"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
)

func AddRoute(mux *http.ServeMux, cfg *config.Config, cacheRepo cache.CacheRepo, userRepo repositories.UserRepo) {

	// general
	mux.HandleFunc("/", handler.HandleRoot())
	mux.HandleFunc("GET /api/v1/healthcheck", handler.HealthCheck(cfg))

	// users
	mux.HandleFunc("GET /api/v1/users", handler.ListUsers(userRepo))
	mux.HandleFunc("GET /api/v1/users/{id}", handler.GetUsers(userRepo))
	mux.HandleFunc("DELETE /api/v1/users/{id}", handler.DeleteUser(userRepo))
	mux.HandleFunc("POST /api/v1/users", handler.CreateUser(userRepo))

	mux.HandleFunc("POST /api/v1/auth/login", handler.AuthenticateUser(cfg, cacheRepo, userRepo))
	mux.HandleFunc("GET /api/v1/auth/refresh-token", handler.RefreshUserAccessToken(cfg, cacheRepo, userRepo))
	mux.Handle("GET /api/v1/auth/logout", handler.LogoutUser(cfg, cacheRepo, userRepo))

}
