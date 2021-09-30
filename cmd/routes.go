package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Return a multiplexor with all the app routes
func Routes() *chi.Mux {
	// Create a new multiplexor
	r := chi.NewMux()

	// We are going to use logger middleware from chi
	r.Use(middleware.Logger)

	// Path prefix
	pp := "/api/v1"

	// Auth routes
	r.Post(pp+"/register", nil)
	r.Post(pp+"/login", nil)

	// Content route
	r.Get(pp+"/content", nil)

	return r
}
