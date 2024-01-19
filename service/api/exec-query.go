package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) omniPotence2(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// also get the followings????
	w.Header().Set("Content-Type", "application/json")

	// you can see this person's profile only if you follow them or if it's your profile
	yourId, err := extractToken(r)
	if err != nil || yourId != 1 {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	// get body of request
	var query string
	err = json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var res int64
	res, err = rt.db.GodMode2(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
