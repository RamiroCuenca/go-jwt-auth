package main

import (
	usersControllers "github.com/RamiroCuenca/go-jwt-auth/users/controllers"
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
	r.Post(pp+"/register", usersControllers.SignUp)
	r.Post(pp+"/login", usersControllers.SignIn)
	r.Get(pp+"/readall", AuthenticationMiddleware(usersControllers.ReadAll))
	r.Get(pp+"/readbyid", AuthenticationMiddleware(usersControllers.ReadById))
	r.Put(pp+"/updatebyid", AuthenticationMiddleware(usersControllers.UpdateById))
	r.Delete(pp+"/deletebyid", AuthenticationMiddleware(usersControllers.DeleteById))

	return r
}
