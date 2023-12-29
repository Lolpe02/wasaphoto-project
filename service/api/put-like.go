package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// receive creator from bearer token, photo id from path, write
func (rt *_router) like(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// read security bearer token from header
	creator, err := extractToken(r)
	if err != nil {
		// bad request
		json.NewEncoder(w).Encode(http.StatusForbidden) // 403
	}
	// Parse the path parameter for the photo id
	var pid int64
	pid, err = readPath(ps, "postId")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	like := like{pid, creator}
	// Call the database function to create the like
	err = rt.db.PutLike(like.PostID, like.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to create like", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
}
