package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RamiroCuenca/go-jwt-auth/auth"
	"github.com/RamiroCuenca/go-jwt-auth/common/handler"
	"github.com/RamiroCuenca/go-jwt-auth/common/logger"
	"github.com/RamiroCuenca/go-jwt-auth/database/connection"
	"github.com/RamiroCuenca/go-jwt-auth/users/models"
	"github.com/RamiroCuenca/go-jwt-auth/utils"
	"github.com/lib/pq"
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
	handler.SendResponse(w, http.StatusCreated, []byte(responseJson), token)
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

	// 8° As the user is valid, generate a JWT
	user := models.User{
		Username: u.Username,
		Email:    u.Email,
		Password: u.HashedPassword,
	}
	token, err := auth.GenerateToken(user)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "User logged successfully but error generating JWT, try loging in again...")
		tx.Rollback()
		return
	}
	logger.Log().Info("JWT generated successfully! :)")

	// 9° If the token was generated successfully, create a Json to send a response
	// Encode the User into a JSON object
	responseJson := fmt.Sprintf(`{
		"Message": "User logged in successfully",
		"Username": %s,
		"JWT": "%s"
	}`, u.Username, token)

	// 9° Send the response
	handler.SendResponse(w, http.StatusCreated, []byte(responseJson), token)
}

// Get all users
//
// The user should be authenticated so it must sent the jwt through the headers
func ReadAll(w http.ResponseWriter, r *http.Request) {
	// 1° Create the sql statement and prepare null fields
	q := `SELECT id, username, email, created_at, updated_at FROM users`

	// 2° Initialize the connection to the database and start a transaction
	db := connection.NewPostgresClient()

	tx, err := db.Begin()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not start the transaction")
		return
	}

	// 3° Prepare the transaction, remember to close it
	stmt, err := tx.Prepare(q)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not prepare the transaction")
		tx.Rollback()
		return
	}
	defer stmt.Close()

	// 4° Execute the query and assign each value to a user object and the append it to an array
	// We will use QueryRow because the exec method returns two methods that are
	// not compatible with psql!
	rows, err := stmt.Query()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not get rows")
		tx.Rollback()
		return
	}
	defer rows.Close()

	// 5° Create a users arr and assign each value
	var usersArr []models.User

	for rows.Next() {
		var u models.User

		// Prepare null values
		var nullUpdated pq.NullTime

		err := rows.Scan(
			&u.Id,
			&u.Username,
			&u.Email,
			&u.CreatedAt,
			&nullUpdated,
		)
		if err != nil {
			sendError(w, http.StatusBadRequest, err, "Could not start the transaction")
			tx.Rollback()
			return
		}

		u.UpdatedAt = nullUpdated.Time

		usersArr = append(usersArr, u)
	}

	// 5° Commit transaction
	tx.Commit()
	logger.Log().Infof("Users fetched successfully! :)")

	// 6° Encode the usersArr in a json
	json, _ := json.Marshal(usersArr)

	// 7° Send response
	handler.SendResponse(w, http.StatusOK, json, "")
}

// Get a specific user by id
//
// The user should be authenticated so it must sent the jwt through the headers
func ReadById(w http.ResponseWriter, r *http.Request) {
	// 1° Get the id from request url
	urlParam := r.URL.Query().Get("id") // Return a string... should convert it to int
	id, err := strconv.Atoi(urlParam)   // Convert it to int
	if err != nil {
		sendError(w, http.StatusBadRequest, err, "Could not fetch the id from url params")
		return
	}

	// 2° Create a used object where the fetched user will be stored
	u := models.User{Id: int64(id)}

	// We are going to create a var where store the updated field in case it's null
	nullUpdateAt := pq.NullTime{}

	// 3° Setup the query to fetch the user
	q := `SELECT id, username, email, created_at, updated_at
			FROM users WHERE id = $1`

	// 4° Init the connection to the database and start a transaction
	db := connection.NewPostgresClient()

	tx, err := db.Begin()
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not start transaction with database")
		return
	}

	// 5° Prepare the transaction for the query (return a statement)
	stmt, err := tx.Prepare(q)
	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not prepare the query to fetch the user")
		tx.Rollback()
		return
	}
	defer stmt.Close()

	// 6° Execute the query and assign the returned values to "u" var

	err = stmt.QueryRow(id).Scan(
		&u.Id,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&nullUpdateAt,
	)

	u.UpdatedAt = nullUpdateAt.Time

	if err != nil {
		sendError(w, http.StatusInternalServerError, err, "Could not fetch a user with sent id")
		tx.Rollback()
		return
	}

	// 7° Commit the transaction
	tx.Commit()

	// 8° Send the response
	json, _ := json.Marshal(u)

	handler.SendResponse(w, http.StatusOK, json, "")

	/*
		// 6° Execute the query and scan the row and assign the values to the note
		err = stmt.QueryRow(n.ID).Scan(
			&n.ID,
			&n.OwnerName,
			&n.Title,
			&nullDetails, // In case it's null
			&n.CreatedAt,
			&nullUpdateAt, // In case it's null
		)
		if err != nil {
			logger.Log().Infof("Error scanning the row: %v", err)
			handler.SendError(w, 500) // Internal Server Error
			return
		}

		n.Details = nullDetails.String
		n.UpdatedAt = nullUpdateAt.Time

		// 8° Encode the Note as Json using Marshal
		json, _ := json.Marshal(n)

		// 7° Commit the transaction
		logger.Log().Info("Record fetched successfully! :)")
		tx.Commit()

		// 9° Send response
		handler.SendResponse(w, http.StatusOK, json)
	*/
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
