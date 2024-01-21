package api

import (
	"encoding/json"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
	// "strconv"
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
	immagine := r.FormValue("photo")
	description := r.FormValue("description")

	// format2 := http.DetectContentType([]byte(immagine))
	if err != nil {
		// could not decode post, bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		err = json.NewEncoder(w).Encode("cant decode image " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	description, err = os.Getwd()
	if err != nil {
		// could not decode post, bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		err = json.NewEncoder(w).Encode("cant mkdir " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	outFile, err := os.Create(description + "tmpimmagini/image.png")
	if err != nil {
		// could not create post, internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("cant create image " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	defer outFile.Close()
	// turn string into image
	var img image.Image
	img, _, err = image.Decode(strings.NewReader(immagine))

	if err != nil {
		// could not decode post, bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		err = json.NewEncoder(w).Encode("cant decode image " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}

	// Encode the image as JPEG and write it to the file
	err = jpeg.Encode(outFile, img, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("cant encode image " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	// create the post

	// _, err = rt.db.CreatePost(img, creator)
	// err = json.NewEncoder(w).Encode("type:  " + format2 + "DESC:  ")
	if err != nil {
		// could not create post, internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("could not create post " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}

	// return the id of the post?? idk
	w.WriteHeader(http.StatusCreated) // 201
	// err = json.NewEncoder(w).Encode(format1 + " " + format2)
}
