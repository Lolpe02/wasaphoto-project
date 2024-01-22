package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var postIds []int64
	// authenticate the user
	yourId, err := extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	// get the list of post ids
	postIds, err = rt.db.GetFeed(yourId)
	if err != nil {
		// could not get likes, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("could not get feed " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	// return the list of user ids
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(postIds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
	}
}
