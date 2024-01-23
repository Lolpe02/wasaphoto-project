package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Set the content type
	w.Header().Set("Content-Type", "application/json")
	// get token
	yourId, err := extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		err = json.NewEncoder(w).Encode("not authenticated")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	// Parse the request body
	var newUsername string
	err = json.NewDecoder(r.Body).Decode(&newUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode("cant decode username, " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	if !isValid(newUsername) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode("invalid username")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}

	// Call the changeUsername database function with the new username
	err = rt.db.ChangeUsername(yourId, newUsername)
	if err != nil {
		if err.Error() == NotFound {
			w.WriteHeader(http.StatusNotFound)
		} else if strings.Contains(err.Error(), UQviolation) {
			w.WriteHeader(http.StatusBadRequest) // 400
			err = json.NewEncoder(w).Encode("username already taken")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
			}
		} else {
			// could not change username, internal server error
			w.WriteHeader(http.StatusInternalServerError) // 500
			err = json.NewEncoder(w).Encode("could not change username, " + err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
			}
		}
		return
	}
}
