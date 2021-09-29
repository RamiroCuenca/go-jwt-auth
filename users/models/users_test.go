package models

import (
	"testing"
	"time"

	"github.com/RamiroCuenca/go-jwt-auth/utils"
)

func TestUserWithCorrectParams(t *testing.T) {
	u := User{
		Id:        12,
		Username:  utils.GenerateRandomString(15),
		Email:     utils.GenerateRandomString(10) + "@example.com",
		Password:  utils.GenerateRandomString(15),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := validate(u); err != nil {
		t.Errorf("❌ Could not create user with correct params: %v.", err)
	} else {
		t.Log("✅ User created with correct params successfully.")
	}
}

func TestUserWithWrongUsername(t *testing.T) {
	u := User{
		Id:        12,
		Username:  "",
		Email:     utils.GenerateRandomString(10) + "@example.com",
		Password:  utils.GenerateRandomString(15),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := validate(u); err == nil {
		t.Errorf("❌ Created user in spite of sending an empty username: %v.", err)
	} else {
		t.Log("✅ Stopped creation of the user with an empty username.")
	}
}

func TestUserWithWrongEmail(t *testing.T) {
	u := User{
		Id:        12,
		Username:  utils.GenerateRandomString(15),
		Email:     utils.GenerateRandomString(10) + "example.com", // It does not contain the @
		Password:  utils.GenerateRandomString(15),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := validate(u); err == nil {
		t.Errorf("❌ Created user in spite of sending an invalid email: %v.", err)
	} else {
		t.Log("✅ Stopped creation of the user with an invalid email.")
	}
}

func TestUserWithWrongPassword(t *testing.T) {
	u := User{
		Id:        12,
		Username:  utils.GenerateRandomString(15),
		Email:     utils.GenerateRandomString(10) + "@example.com",
		Password:  utils.GenerateRandomString(5), // Less than 6 digits
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := validate(u); err == nil {
		t.Errorf("❌ Created user in spite of sending an invalid password: %v.", err)
	} else {
		t.Log("✅ Stopped creation of the user with an invalid password.")
	}
}
