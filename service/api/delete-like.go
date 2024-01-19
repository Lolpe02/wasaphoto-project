package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// receive creator from bearer token, photo id from path, write
func (rt *_router) unlike(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// read security bearer token from header
	creator, err := extractToken(r)
	if err != nil {
		// bad request
		err = json.NewEncoder(w).Encode(http.StatusForbidden) // 403
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	// Parse the path parameter for the photo id
	var pid int64
	pid, err = readPath(ps, "postId")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the database function to delete the like
	err = rt.db.Unlike(pid, creator)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to delete like", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}
