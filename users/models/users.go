package models

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Id        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func validate(u User) error {
	// Check that the Username is not empty.
	if u.Username == "" {
		return errors.New("Username can not be empty")
	}

	// Check that the username has less than 51 digits.
	if len(u.Username) > 50 {
		return errors.New("Username can not be larger than 50 digits")
	}
	// Check that the Email is not empty.
	if u.Email == "" {
		return errors.New("Email can not be empty")
	}

	// Check that the Email contains an @.
	if !strings.Contains(u.Email, "@") {
		return errors.New("Email must be valid (include @)")
	}

	// Check that the Password is longer than 6 digits.
	if len(u.Password) < 6 {
		return errors.New("Password must have at least 6 digits")
	}

	return nil
}

func Check(u User) error {
	return validate(u)
}

// Claim is the information that will be sent through the JWT
// In this case we are going to add only the username. The
// other parameters are going to be automatically added by
// the library that we are going to use.
type Claim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
