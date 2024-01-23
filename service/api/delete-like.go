package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// receive creator from bearer token, photo id from path, write
func (rt *_router) unlike(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// read security bearer token from header
	creator, err := extractToken(r)
	if err != nil {
		// bad request
		err = json.NewEncoder(w).Encode(http.StatusForbidden) // 403
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	// Call the database function to delete the like
	err = rt.db.Unlike(pid, creator)
	if err != nil {
		if err.Error() == NotFound {
			w.WriteHeader(http.StatusNotFound)
			err = json.NewEncoder(w).Encode("like not found")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		log.Println(err)
		http.Error(w, "Failed to delete like", http.StatusInternalServerError)
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("like deleted")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
