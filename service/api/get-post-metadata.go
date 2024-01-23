package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetPostMetadata(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	var userId int64
	var description string
	var date string
	userId, description, date, err = rt.db.GetMetadata(postId)
	if err != nil {
		if err.Error() == NotFound {
			// post not found, throw not found
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}
		// could not get creator, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	retrieved := post{userId, description, date}
	var present bool
	if _, present, err = rt.db.GetFolloweds(retrieved.Creator, yourId); err != nil {
		// could not get follows, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

	if !present && retrieved.Creator != yourId {
		// you are not following this person, throw forbidden
		w.WriteHeader(http.StatusForbidden) // 403
		return
	}

	// return the list of post ids of that user
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(retrieved)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
	}

}
