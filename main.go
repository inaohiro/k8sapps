package main

import (
	"inpodk8soperation/server"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Mount("/servers", server.Routes)

	http.ListenAndServe(":8080", r)
}
