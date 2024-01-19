package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) omniPotence1(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		err = json.NewEncoder(w).Encode("error decoding body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var res []map[string]interface{}
	res, err = rt.db.GodMode1(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode("error executing query " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode("error encoding response " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}
}
