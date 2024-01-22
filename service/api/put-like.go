package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// receive creator from bearer token, photo id from path, write
func (rt *_router) like(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// read security bearer token from header
	creator, err := extractToken(r)
	if err != nil {
		// bad request
		json.NewEncoder(w).Encode(http.StatusForbidden) // 403
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
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode("Unauthorized")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if yourId != creator {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode("something went wrong with your id and token")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	like := like{pid, creator}
	// Call the database function to create the like
	err = rt.db.PutLike(like.PostId, like.UserId)
	if err != nil {
		if err.Error() == "already liked" {
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode("already liked")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		log.Println(err)
		err = json.NewEncoder(w).Encode("Failed to create like")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
