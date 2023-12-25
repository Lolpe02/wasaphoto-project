package api

import (
	"encoding/json"
	"net/http"
)

func GetNameHandler(w http.ResponseWriter, r *http.Request) {

	// Query the database for a user
	var user User

	// Convert the user to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the JSON response
	w.Write(userJSON)
}
