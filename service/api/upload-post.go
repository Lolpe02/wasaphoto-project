package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"mime/multipart"
	"net/http"
	"strings"
)

func (rt *_router) upload(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// take token from the header
	creator, err := extractToken(r)
	if err != nil {
		// could not parse token, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	creator += 0
	// read multipart form
	var file multipart.File
	var handler *multipart.FileHeader
	description := r.FormValue("description")
	file, handler, err = r.FormFile("photo")

	if err != nil {
		// could not read file, bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		err = json.NewEncoder(w).Encode("cant read file " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}

	mime := handler.Header.Get("Content-Type")
	parts := strings.Split(mime, "/")
	if len(parts) != 2 {
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	typeM, enc := parts[0], parts[1]
	if typeM != "image" {
		// could not read file, bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		err = json.NewEncoder(w).Encode("not an image " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		if enc != "png" && enc != "jpeg" && enc != "gif" {
			w.WriteHeader(http.StatusBadRequest) // 400
			err = json.NewEncoder(w).Encode("format not supported, " + err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
			}
		}
		return
	}

	var postId int64
	postId, err = rt.db.CreatePost(&file, &description, enc, creator)
	if err != nil {
		// could not create post, internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("could not create post " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	if err = file.Close(); err != nil {
		// could not read file, bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		err = json.NewEncoder(w).Encode("cant close file " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}

	// return the id of the post?? idk
	w.WriteHeader(http.StatusCreated) // 201
	err = json.NewEncoder(w).Encode(postId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
	}
}
