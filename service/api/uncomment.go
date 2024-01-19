package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncomment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// take token from the header
	var commentId int64
	creator, err := extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	// take parameters from the path and turn string to int64
	commentId, err = readPath(ps, "commentId")
	if err != nil {
		// could not parse the post id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	// put the comment in the database
	err = rt.db.Uncomment(creator, commentId)
	if err != nil {
		// could not create comment, internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	// return the id of the comment?? idk
	w.WriteHeader(http.StatusCreated) // 201
}
