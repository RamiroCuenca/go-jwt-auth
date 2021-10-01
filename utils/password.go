package utils

import (
	"errors"

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
		return "", errors.New(err.Error())
	}

	return string(hashedPass), nil
}

// Checks if the provided password is correct or not
func PasswordCheck(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
