package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getFollowersOf(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var followerNames []string

	// authenticate the user

	token, err := extractToken(r)
	if err != nil || token < 0 {
		// could not extract the token, throw unauthorized error
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	query := r.URL.Query().Get("userName")
	if query == "" {
		// if the query is empty,
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	// get the userId from the userName
	var userId int64
	userId, err = rt.db.SearchByUsername(query)
	if err != nil {
		// could not get the userId, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("Internal Server Error" + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	_, followerNames, _, err = rt.db.GetFollowing(userId, -1)
	if err != nil {
		// could not get user names, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("Internal Server Error" + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// return the list of user ids
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(followerNames)
	if err != nil {
		// could not encode the list, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
}
