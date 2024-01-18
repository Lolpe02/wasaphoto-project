package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// also get the followings????
	w.Header().Set("content-type", "application/json")
	var postIds []int64
	// take username parameters from query
	username := r.URL.Query().Get("userName")
	if username == "" {
		// could not parse the id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	targetId, err := rt.db.SearchByUsername(username)
	if err != nil {
		// could not get id, throw internal server error
		if err.Error() == "user not found" {
			w.WriteHeader(http.StatusNotFound) // 404
			_, err = w.Write([]byte("user not found " + username))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		}
		w.WriteHeader(http.StatusInternalServerError) // 500
		_, err = w.Write([]byte("couldnt search by useranem" + username))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	// you can see this person's profile only if you follow them or if it's your profile
	yourId, err := extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	var followed []int64
	var present bool

	if followed, present, err = rt.db.GetFolloweds(targetId, yourId); present || yourId == targetId {
		// you own this profile or you follow it, you can see it
		if err != nil {
			// could not get follows, throw internal server error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// you are not following this person, throw forbidden
		w.WriteHeader(http.StatusForbidden) // 403
		return

	}
	var selname, sub string
	selname, sub, err = rt.db.SearchById(targetId)
	if err != nil {
		// could not get likes, throw internal server error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	postIds, err = rt.db.GetProfile(targetId)
	if err != nil {
		// could not get likes, throw internal server error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	follows, _, err := rt.db.GetFollowing(targetId, -1)
	if err != nil {
		// could not get follows, throw internal server error
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}
	user := user{targetId, selname, sub, postIds, follows, followed}

	// return the list of post ids of that user
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
