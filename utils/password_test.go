package utils

import "testing"

// Verify that PasswordHasher is hashing properly.
func TestPasswordHasher(t *testing.T) {
	p := "pass123"

	// Generate the hashed password
	hashedPassword, err := PasswordHash(p)
	if err != nil {
		t.Errorf("There was an error hashing the password: %v", err)
	}
	if hashedPassword == "" {
		t.Errorf("The hashed password can not be empty")
	} else {
		t.Logf("Password hashed succesfully! :) \n HashedPassword: %v", hashedPassword)
	}
}

// Verify that PasswordVerify works properly.
func TestPasswordVerifierWithCorrectPassword(t *testing.T) {
	p := "pass123"

	// Generate the hashed password
	hashedPassword, err := PasswordHash(p)
	if err != nil {
		t.Errorf("There was an error hashing the password: %v", err)
	}
	if hashedPassword == "" {
		t.Errorf("The hashed password can not be empty")
	} else {
		t.Logf("Password hashed succesfully! :) \n HashedPassword: %v", hashedPassword)
	}

	// Check the password
	err = PasswordCheck(p, hashedPassword)
	if err != nil {
		t.Errorf("The password check failed: %v", err)
	} else {
		t.Log("Password checked succesfully! :)")
	}
}

// Verify that PasswordVerify works properly.
func TestPasswordVerifierWithWrongPassword(t *testing.T) {
	p := "pass123"

	// Generate the hashed password
	hashedPassword, err := PasswordHash(p)
	if err != nil {
		t.Errorf("There was an error hashing the password: %v", err)
	}
	if hashedPassword == "" {
		t.Errorf("The hashed password can not be empty")
	} else {
		t.Logf("Password hashed succesfully! :) \n HashedPassword: %v", hashedPassword)
	}

	// Check an incorrect password
	incorrectPassword := "pass12345"
	err = PasswordCheck(incorrectPassword, hashedPassword)
	if err == nil {
		t.Errorf("The password check failed: %v", err)
	} else {
		t.Log("Password checked succesfully! :)")
	}
}

// Verify that every different password generates a different hash
func TestPasswordHasherWithTwoDifferentPasswords(t *testing.T) {
	p := "pass123"

	// Generate the hashed password
	hashedPassword, err := PasswordHash(p)
	if err != nil {
		t.Errorf("There was an error hashing the password: %v", err)
	}
	if hashedPassword == "" {
		t.Errorf("The hashed password can not be empty")
	} else {
		t.Logf("Password hashed succesfully! :) \n HashedPassword: %v", hashedPassword)
	}

	p2 := "pass12345"
	hashedPassword2, err := PasswordHash(p2)
	if err != nil {
		t.Errorf("There was an error hashing the password: %v", err)
	}
	if hashedPassword2 == "" {
		t.Errorf("The hashed password can not be empty")
	} else {
		t.Logf("Password2 hashed succesfully! :) \n HashedPassword2: %v", hashedPassword2)
	}

	if hashedPassword == hashedPassword2 {
		t.Errorf("The hash of diferent passwords can not be equal")
	} else {
		t.Log("PasswordHasher is working properly")
	}
}
