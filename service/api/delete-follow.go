package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) unfollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// take id parameters from the path (person to unfollow)
	IdtounFollow, err := readPath(ps, "followId")
	if err != nil {
		// could not parse the id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		err = json.NewEncoder(w).Encode(err.Error() + ", you fool")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		return
	}

	var yourId int64
	yourId, err = extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	err = rt.db.UnfollowUser(yourId, IdtounFollow)
	if err != nil {
		if err.Error() == "not following user" {
			w.WriteHeader(http.StatusNotFound) // 404
			err = json.NewEncoder(w).Encode(err.Error() + ", you fool")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				return
			}
			return
		} else {
			// could not unfollow, throw internal server error
			w.WriteHeader(http.StatusInternalServerError) // 500
			err = json.NewEncoder(w).Encode(err.Error() + ", sopmething went wrong")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				return
			}
			return
		}
	}
	err = json.NewEncoder(w).Encode("Unfollowed")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
}
