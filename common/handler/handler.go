package handler

import "net/http"

// Common function wich will be used to send responses for correct http requests
func SendResponse(w http.ResponseWriter, status int, data []byte, token string) {

	w.Header().Set("Content-Type", "application/json")
	// If there is a token, add it to the header
	if token != "" {
		w.Header().Set("Token", token)
	}
	w.WriteHeader(status)
	w.Write(data)

}

// Common function wich will be used to send errors for incorrect http requests
func SendError(w http.ResponseWriter, status int, data []byte) {

	// data := []byte(`{}`)

	w.Header().Set("Content-Type", "application/")
	w.WriteHeader(status)
	w.Write(data)
}
