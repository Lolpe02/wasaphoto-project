package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) unfollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// take id parameters from the path (person to unfollow)
	IdtounFollow, err := readPath(ps, "followId")
	if err != nil {
		// could not parse the id, throw bad request
		w.WriteHeader(http.StatusBadRequest) //400
		return
	}

	var yourId int64
	yourId, err = extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) //401
		return
	}
	err = rt.db.UnfollowUser(yourId, IdtounFollow)
	if err != nil {
		// could not follow, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) //500
		return
	}

	// return the list of post ids of that user
	w.WriteHeader(http.StatusOK) //200
	return
}
