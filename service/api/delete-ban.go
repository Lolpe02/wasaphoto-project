package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unban(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// take id parameters from the path (person to follow)
	IdtoUnban, err := readPath(ps, "bannedId")
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

	err = rt.db.UnbanUser(yourId, IdtoUnban)
	if err != nil {
		if err.Error() == NotFound {
			// could not unban, throw not found
			w.WriteHeader(http.StatusNotFound) // 404
			err = json.NewEncoder(w).Encode("user not banned " + err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
			}
			return
		}
		// could not follow, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("could not unban, " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		return
	}
	w.WriteHeader(http.StatusOK) // 200
}
