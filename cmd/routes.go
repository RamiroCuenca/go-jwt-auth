package main

import (
	"net/http"

	"github.com/RamiroCuenca/go-jwt-auth/common/handler"
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

	// Content route
	r.Get(pp+"/test", AuthenticationMiddleware(testAuthMiddleware))

	return r
}

func testAuthMiddleware(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello world! AuthenticationMiddleware working properly! :)")
	data := `{
		"message": "Hello world! AuthenticationMiddleware working properly! :)"
	}`
	handler.SendResponse(w, http.StatusOK, []byte(data), "")
}
