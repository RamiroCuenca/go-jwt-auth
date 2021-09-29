package models

import (
	"errors"
	"strings"
	"time"
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
		return errors.New("Email can not be empty")
	}

	// Check that the Password is longer than 6 digits.
	if len(u.Password) < 6 {
		return errors.New("Password must have at least 6 digits")
	}

	return nil
}
