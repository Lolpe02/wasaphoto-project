package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// receive creator from bearer token, photo id from path, write
func likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// read security bearer token from header
	creator, err := ExtractToken(r)
	if err != nil {
		// bad request
		json.NewEncoder(w).Encode(http.StatusBadRequest) //
	}
	// Parse the request body
	var pid int64
	err := json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the database function to create the like
	err = CreateLike(like)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create like", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}

func CreateLike(like Like) error {
	// Database logic to create the like goes here
	// ...

	return nil
}
