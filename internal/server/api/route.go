package api

import (
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/api/handler"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/repositories"
)

func AddRoute(mux *http.ServeMux, cfg *config.Config, userRepo repositories.UserRepo) {

	// general
	mux.HandleFunc("/", handler.HandleRoot())
	mux.HandleFunc("GET /api/v1/healthcheck", handler.HealthCheck(cfg))

	// users
	mux.HandleFunc("GET /api/v1/users/{id}", handler.GetUsers(userRepo))
	mux.HandleFunc("DELETE /api/v1/users/{id}", handler.DeleteUser(userRepo))

	//login
	mux.HandleFunc("POST /api/v1/users", handler.CreateUser(userRepo))
	mux.HandleFunc("POST /api/v1/login", handler.CreateUser(userRepo))

}
