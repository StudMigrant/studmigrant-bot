package handler

import "net/http"

func NewRouter(webhook http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /webhook", webhook)
	mux.Handle("GET /healthz", healthz())

	return mux
}
