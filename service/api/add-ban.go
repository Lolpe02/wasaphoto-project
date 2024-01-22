package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) ban(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

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
	/*
		// check if user exists
		_, _, err = rt.db.SearchById(IdtoBan)
		if err != nil {
			if err.Error() == "not found" {
				// could not follow, throw not found
				w.WriteHeader(http.StatusNotFound) // 404
			} else {
				// could not follow, throw internal server error
				w.WriteHeader(http.StatusInternalServerError) // 500
				err = json.NewEncoder(w).Encode("could not ban" + err.Error())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError) // 500
					return
				}
			}
			return
		}*/

	ban := ban{yourId, IdtoBan}
	err = rt.db.BanUser(ban.BannerId, ban.BannedId)
	if err != nil {
		if err.Error() == "already banned this user" {
			w.WriteHeader(http.StatusOK) // 200
			err = json.NewEncoder(w).Encode(err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				return
			}
		} else if strings.Contains(err.Error(), "FOREIGN KEY constraint failed") {
			// could not follow, throw not found
			w.WriteHeader(http.StatusNotFound) // 404
			err = json.NewEncoder(w).Encode("could not ban user, it doesnt exists " + err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				return
			}
		} else {
			// could not follow, throw internal server error
			w.WriteHeader(http.StatusInternalServerError) // 500
			err = json.NewEncoder(w).Encode("could not ban" + err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				return
			}
		}
		return
	}
	// return the list of post ids of that user
	w.WriteHeader(http.StatusCreated) // 201
	err = json.NewEncoder(w).Encode("now banned")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
}
