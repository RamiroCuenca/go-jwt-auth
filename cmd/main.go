package main

import (
	"github.com/RamiroCuenca/go-jwt-auth/common/logger"
	"github.com/RamiroCuenca/go-jwt-auth/database/connection"
)

func main() {
	// Init zap logger
	logger.InitZapLogger()

	// Init postgre database
	connection.NewPostgresClient()

	// Get the router
	mux := Routes()

	// Setup the server
	sv := NewServer(mux)

	// Run the server
	logger.Log().Info("Server running over port :8000 ...\n")
	sv.Run()
}
