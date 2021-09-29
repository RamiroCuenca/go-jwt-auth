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
}
