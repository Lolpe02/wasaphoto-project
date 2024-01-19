package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) follow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// take id parameters from the path (person to follow)
	var IdtoFollow int64
	err := json.NewDecoder(r.Body).Decode(&IdtoFollow)
	if err != nil {
		// could not parse the id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	var yourId int64
	yourId, err = extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	follow := follow{yourId, IdtoFollow}
	var exists bool
	exists, err = rt.db.FollowUser(follow.FollowingId, follow.FollowedId)
	if err != nil {
		// could not follow, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("could not follow user " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		return
	}

	// return success
	switch exists {
	case true:
		w.WriteHeader(http.StatusOK) // 200
		err = json.NewEncoder(w).Encode("user already followed")
	case false:
		w.WriteHeader(http.StatusCreated) // 201
		err = json.NewEncoder(w).Encode("now following user")
	}
	if err != nil {
		// could not write response, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 00
		err =  json.NewEncoder(w).Encode("could not write response " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		return
	}
}
