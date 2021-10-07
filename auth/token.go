package auth

import (
	"errors"
	"time"

	"github.com/RamiroCuenca/go-jwt-auth/users/models"
	"github.com/golang-jwt/jwt/v4"
)

// Generates a JWT. It receives the data from the user who has logged in!
// The JWT is a string
func GenerateToken(user models.User) (string, error) {
	claim := models.Claim{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // 2hs ot expire
			Issuer:    "Ramiro Cuenca Salinas",
		},
	}

	// Prepared the token to be signed with RS256 method and including the claim
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)

	// Sign the token with our private key
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	// Return the generated token
	return signedToken, nil
}

// Validate the JWT.
// Returns the claim so that we can access the information of the user.
//
// Before, we need to create a function wich return as the information of
// our public key (Already parsed, it is declared as "verifyKey")
func ValidateToken(t string) (models.Claim, error) {

	token, err := jwt.ParseWithClaims(t, &models.Claim{}, verifyFunction)
	if err != nil {
		return models.Claim{}, err
	}

	// Check if the token is valid
	if !token.Valid {
		return models.Claim{}, errors.New("Invalid token")
	}

	// Obtain the claims from the token
	claim, ok := token.Claims.(*models.Claim)
	if !ok {
		return models.Claim{}, errors.New("Couldn't fetch the claims")
	}

	return *claim, nil
}

// Returns the verifyKey wich is out public key already parsed
func verifyFunction(t *jwt.Token) (interface{}, error) {
	return verifyKey, nil
}
