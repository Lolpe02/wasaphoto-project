package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) upload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// take token from the header
	creator, err := extractToken(r)

	if err != nil {
		// could not parse token, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	var newImage []byte
	err = json.NewDecoder(r.Body).Decode(&newImage)
	if err != nil {
		// could not decode post, bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// create the post
	_, err = rt.db.CreatePost(newImage, creator)
	if err != nil {
		// could not create post, internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	// return the id of the post?? idk
	w.WriteHeader(http.StatusCreated) // 201
}
