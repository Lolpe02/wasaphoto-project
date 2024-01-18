package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) ban(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// take id parameters from the body ()
	var IdtoBan int64
	err := json.NewDecoder(r.Body).Decode(&IdtoBan)
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
	ban := ban{IdtoBan, yourId}
	err = rt.db.BanUser(ban.BannerId, ban.BannedId)
	if err != nil {
		// could not follow, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	// return the list of post ids of that user
	w.WriteHeader(http.StatusCreated) // 200

}
