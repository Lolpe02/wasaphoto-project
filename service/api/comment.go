package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) comment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// take token from the header

	creator, err := extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	// take parameters from the path and turn string to int64
	var postId int64
	postId, err = readPath(ps, "postId")
	if err != nil {
		// could not parse the post id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	// read the body of the request
	var content string
	err = json.NewDecoder(r.Body).Decode(&content)
	if err != nil {
		// could not get likes, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	// put the comment in the database
	_, err = rt.db.PutComment(creator, content, postId)
	if err != nil {
		// could not create comment, internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	// return the id of the comment?? idk
	w.WriteHeader(http.StatusCreated) // 201
}
