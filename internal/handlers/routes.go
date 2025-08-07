package handlers

import (
    "net/http"
)

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
    mux.HandleFunc("/api/shorten", handler.Shorten)
    mux.HandleFunc("/api/shorten/", handler.Shorten)
    mux.HandleFunc("/", handler.Redirect)
    mux.HandleFunc("/ping", handler.Ping)
}