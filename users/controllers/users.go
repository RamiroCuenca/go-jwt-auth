package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/RamiroCuenca/go-jwt-auth/auth"
	"github.com/RamiroCuenca/go-jwt-auth/common/handler"
	"github.com/RamiroCuenca/go-jwt-auth/common/logger"
	"github.com/RamiroCuenca/go-jwt-auth/database/connection"
	"github.com/RamiroCuenca/go-jwt-auth/users/models"
	"github.com/RamiroCuenca/go-jwt-auth/utils"
)

// Registers a new user account
//
// It must recieve username, email and password as parameters (All strings...).
func SignUp(w http.ResponseWriter, r *http.Request) {
	// 1° Decode the json received on an User object
	u := models.User{}

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		sendError(w, http.StatusBadRequest, err, "Could not decode request body")
		return
	}

	// Check if user fields are valid
	err = models.Check(u)
	if err != nil {
		sendError(w, http.StatusBadRequest, err, err.Error())
		return
	}

	// Hash the password and replace it on the User field
	u.Password, err = utils.PasswordHash(u.Password)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not hash the password")
		return
	}

	// 2° Create the sql statement and prepare null fields
	q := `
	INSERT INTO users (username, email, hashed_password, created_at)
	VALUES ($1, $2, $3, now()) 
	RETURNING id, created_at
	`

	// 3° Initialize database connection
	db := connection.NewPostgresClient()
	// defer db.Close()

	// 4° Start a transaction
	tx, err := db.Begin()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not start the transaction")
		return
	}

	// 5° Prepare the transaction
	stmt, err := tx.Prepare(q)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not start the transaction")
		tx.Rollback()
		return
	}
	defer stmt.Close()

	// 6° Execute the query and assign the returned id to the Note object
	// We will use QueryRow because the exec method returns two methods that are
	// not compatible with psql!
	err = stmt.QueryRow(
		u.Username,
		u.Email,
		u.Password,
	).Scan(
		&u.Id,
		&u.CreatedAt,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not create the user")
		tx.Rollback()
		return
	}

	// 7° As there are no errors, commit the transaction
	tx.Commit()
	logger.Log().Infof("User created successfully! :)")

	// 8° As the user is valid, generate a JWT
	token, err := auth.GenerateToken(u)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "User created successfully but error generating JWT, try loging in...")
		tx.Rollback()
		return
	}
	logger.Log().Info("JWT generated successfully! :)")

	// 9° If the token was generated successfully, create a Json to send a response
	// Encode the User into a JSON object
	userJson, _ := json.Marshal(u)

	responseJson := fmt.Sprintf(`{
		"Message": "User created successfully",
		"User": %s,
		"JWT": "%s"
	}`, userJson, token)

	// 9° Send the response
	handler.SendResponse(w, http.StatusCreated, []byte(responseJson))
}

// Log in an existing user
//
// It must recieve email and password as parameters (Both strings...).
func SignIn(w http.ResponseWriter, r *http.Request) {
	// 1° Decode the json received on an User object
	type loginUserCMD struct {
		Username       string `json:"username"`
		Email          string `json:"email"`
		Password       string `json:"password"`
		HashedPassword string `json:"hashed_password"`
	}

	u := loginUserCMD{}

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		sendError(w, http.StatusBadRequest, err, "Could not decode request body")
		return
	}

	// Check if user fields are valid
	if u.Email == "" || u.Password == "" {
		sendError(w, http.StatusBadRequest, errors.New("Email and Password are required"), "Email and Password are required")
		return
	}

	// 2° Create the sql statement and prepare null fields
	q := `
	SELECT username, hashed_password FROM users
	WHERE email = $1
	`

	// 3° Initialize database connection
	db := connection.NewPostgresClient()
	// defer db.Close()

	// 4° Start a transaction
	tx, err := db.Begin()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not start the transaction")
		return
	}

	// 5° Prepare the transaction
	stmt, err := tx.Prepare(q)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not start the transaction")
		tx.Rollback()
		return
	}
	defer stmt.Close()

	// 6° Execute the query and assign the returned id to the Note object
	// We will use QueryRow because the exec method returns two methods that are
	// not compatible with psql!
	err = stmt.QueryRow(
		u.Email,
	).Scan(
		&u.Username,
		&u.HashedPassword,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "User not found")
		tx.Rollback()
		return
	}

	// 7° As there are no errors, commit the transaction
	tx.Commit()

	// 8° Compare password received and hashed password from the server
	err = utils.PasswordCheck(u.Password, u.HashedPassword)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Incorrect password")
		tx.Rollback()
		return
	}

	logger.Log().Infof("User logged successfully! :)")

	// 8° Encode the Note into a JSON object
	// json, _ := json.Marshal(u)
	data := fmt.Sprintf(`{
		"message": "User logged in successfully! :)",
		"username": "%s"
	}`, u.Username)

	// 9° Send the response
	handler.SendResponse(w, http.StatusCreated, []byte(data))
}

func sendError(w http.ResponseWriter, status int, err error, message string) {
	// Log the error
	logger.Log().Infof(message, ": ", err)

	// Set a json with the error message
	data := fmt.Sprintf(`{
	"message": "%s",
	"error": "%s"
}`, message, err)

	json := []byte(data)
	handler.SendError(w, status, json)
}
