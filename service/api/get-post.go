package api

import (
	"encoding/json"
	"net/http"

	// "os"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "image/*")

	// take id parameters from the path
	postId, err := readPath(ps, "postId")
	if err != nil {
		// could not parse the id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// you can see this person's photo only if you follow them
	var yourId int64
	yourId, err = extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	// get post creator
	var creator int64
	creator, _, _, err = rt.db.GetMetadata(postId)
	if err != nil {
		if err.Error() == NotFound {
			// post not found, throw not found
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}
		// could not get creator, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	// check if you are following this person
	var present bool
	_, present, err = rt.db.GetFolloweds(creator, yourId)
	if err != nil {
		// could not get follows, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	if !present && creator != yourId {
		// you are not following this person, throw forbidden
		w.WriteHeader(http.StatusForbidden) // 403
		err = json.NewEncoder(w).Encode("you are not following this person or its not your post")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	// put retrieved post in post object
	// var pointer *os.File
	var imageBytes *[]byte
	_, imageBytes, err = rt.db.GetPost(postId)
	if err != nil {
		if err.Error() == NotFound {
			// post not found, throw not found
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}
		// could not get creator, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	if imageBytes == nil {
		// image not found, throw not found
		w.WriteHeader(http.StatusNotFound) // 404
		err = json.NewEncoder(w).Encode("nil pointer")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	_, err = w.Write(*imageBytes)
	if err != nil {
		// could not write image, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("couldnt encode image, " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	// return the list of post ids of that user
	w.WriteHeader(http.StatusOK) // 200

}
