package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// authenticate the user
	_, err := extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	// get the user from the query
	var targetUser string
	targetUser = r.URL.Query().Get("userName")
	if targetUser == "" {
		// no user specified, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	var userIdList []int64
	var userNameList []string
	userIdList, userNameList, err = rt.db.GetInfo(targetUser)
	if err != nil {
		if err.Error() == NotFound {
			// user not found, throw not found
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}
		// could not get names, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("could not get Users alike, " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	if len(userIdList) != len(userNameList) {
		// could not get names, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("Something went wrong")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
		}
		return
	}
	// create map string to id
	var idToName = make(map[string]int64)
	for i := 0; i < len(userIdList); i++ {
		idToName[userNameList[i]] = userIdList[i]
	}
	// return the list of user ids
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(idToName)
	// how do i specify this in the api doc?

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
	}
}
