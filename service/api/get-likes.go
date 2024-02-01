package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var userIds []int64
	// take parameters from the path and turn string to int64
	postId, err := readPath(ps, "postId")
	if err != nil {
		// could not parse the post id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	// authenticate the user
	_, err = extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	// check if user follows creator
	/* var postCreator int64
	postCreator, _, _, err = rt.db.GetMetadata(postId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode("couldnt get post owner")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	var follows bool
	_, _, follows, err = rt.db.GetFollowing(postCreator, authUserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 401
		err = json.NewEncoder(w).Encode("Unauthorized" + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if !follows {
		w.WriteHeader(http.StatusUnauthorized) // 401
		err = json.NewEncoder(w).Encode("you dont follow the owner")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	} */
	userIds, err = rt.db.GetLikes(postId)
	if err != nil {
		// could not get likes, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	// return the list of userNames
	var userNames []string
	for _, id := range userIds {
		var userName string
		userName, _, err = rt.db.SearchById(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		userNames = append(userNames, userName)
	}
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(userNames)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
}
