package api

import (
	"net/http"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/config"
)

func NewHTTPServer(cfg *config.Config) http.Handler {

	mux := http.NewServeMux()
	AddRoute(mux, cfg)

	return mux
}
