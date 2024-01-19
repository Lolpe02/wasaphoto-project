package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// take id parameters from the path
	postId, err := readPath(ps, "postId")
	if err != nil {
		// could not parse the id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	// you can see this person's photo only if you follow them
	var yourId int64
	yourId, err = extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	// put retrieved post in post object

	targetPost, creator, date, err := rt.db.GetPost(postId)
	if err != nil {
		// could not get creator, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	retrieved := post{targetPost, postId, creator, date}
	_, present, err1 := rt.db.GetFolloweds(retrieved.Creator, yourId)
	if err1 != nil {
		// could not get follows, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	if !present {
		// you are not following this person, throw forbidden
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	// return the list of post ids of that user
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(retrieved.Image)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

}
