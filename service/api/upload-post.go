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

	// method 1
	newImage := make([]byte, r.ContentLength) // create a byte array of the size of the image
	_, err = r.Body.Read(newImage)
	if err != nil {
		// could not read file, bad request
		w.WriteHeader(http.StatusBadRequest) //400
		err = json.NewEncoder(w).Encode("cant read file " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
		}
		return
	}
	format1 := http.DetectContentType(newImage)
	// read binary string from body
	// method 2
	file, _, err1 := r.FormFile("image")
	if err1 != nil {
		// could not read file, bad request
		w.WriteHeader(http.StatusBadRequest) //400
		err = json.NewEncoder(w).Encode("cant read file " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) //500
		}
		return
	}
	// turn file into byte array
	_, err = file.Read(newImage)
	format2 := http.DetectContentType(newImage)

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
	err = json.NewEncoder(w).Encode(format1 + " " + format2)
}
