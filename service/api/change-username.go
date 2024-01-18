package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Set the content type
	w.Header().Set("content-type", "application/json")
	// get token
	yourId, err := extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		_, err = w.Write([]byte("not authenticated\n"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusUnauthorized) // 01
		return
	}
	// Parse the request body
	var newUsername string
	err = json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !isValid(newUsername) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the changeUsername database function with the new username
	err = rt.db.ChangeUsername(yourId, newUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}
