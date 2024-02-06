package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	// "strconv"
	// "strings"
)

// receive creator from bearer token, photo id from path, write
func (rt *_router) like(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// read security bearer token from header
	likeCreator, err := extractToken(r)
	if err != nil {
		// bad request
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	// Parse the path parameter for the photo id
	var pid int64
	pid, err = readPath(ps, "postId")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var yourId int64
	yourId, err = readPath(ps, "yourId")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // 400
		err = json.NewEncoder(w).Encode("wrong Id" + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if yourId != likeCreator {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode("something went wrong with your id and token, " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	like := like{likeCreator, pid}
	// Call the database function to create the like
	err = rt.db.PutLike(like.PostId, like.UserId)
	if err != nil {
		if err.Error() == "already liked" {
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode("already liked")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode("Failed to create like")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated) // 201
	err = json.NewEncoder(w).Encode("like created")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
