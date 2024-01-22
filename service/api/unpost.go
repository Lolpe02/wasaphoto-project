package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	// put the comment in the database
	err = rt.db.Unpost(creator, postId)
	if err != nil {
		if err.Error() == "not found" {
			w.WriteHeader(http.StatusNotFound) // 404
			err = json.NewEncoder(w).Encode("post not found")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
			}
			return

		}
		// could not create comment, internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("could not delete post " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	w.WriteHeader(http.StatusOK) // 200
}
