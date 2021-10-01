package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Receives a password formatted a string and
// hashes it. The it return the generated hash.
//
// It uses the bcrypt library
func PasswordHash(password string) (string, error) {
	// Generates the hash
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to generate hashed password: %v", err)
	}

	return string(hashedPass), nil
}

// Checks if the provided password is correct or not
func PasswordCheck(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
