package main

import "net/http"

func newServer() http.Handler {

	handler := http.NewServeMux()

	return handler
}
