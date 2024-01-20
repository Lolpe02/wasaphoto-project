package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) upload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// take token from the header
	creator, err := extractToken(r)
	if err != nil {
		// could not parse token, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) //401
		return
	}

	var newImage []byte
	// read binary string from body

	err = json.NewDecoder(r.Body).Decode(&newImage)
	if err != nil {
		// could not decode post, bad request
		w.WriteHeader(http.StatusBadRequest) //400
		err = json.NewEncoder(w).Encode("cant decode image " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
		}
		return
	}

	// create the post
	_, err = rt.db.CreatePost(newImage, creator)
	if err != nil {
		// could not create post, internal server error
		w.WriteHeader(http.StatusInternalServerError) //500
		err = json.NewEncoder(w).Encode("could not create post " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
		}
		return
	}

	// return the id of the post?? idk
	w.WriteHeader(http.StatusCreated) //201
}
