package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var followingNames []string

	// authenticate the user

	_, err := extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	// if the user is not authenticated, throw unauthorized

	// get the userName from the url
	var query string
	query = r.URL.Query().Get("userName")
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
		return
	}

	_, followingNames, _, err = rt.db.GetFolloweds(userId, -1)
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
	err = json.NewEncoder(w).Encode(followingNames)
	if err != nil {
		// could not encode the list, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
}
