package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getFollowersOf(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	
	var followerNames []string

	// authenticate the user
	
	authUserId, err := extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	// get the userName from the url
	var query string
	query = r.URL.Query().Get("userName")
	if query == "" {
		// if the query is empty,  
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	// get the userId from the userName
	var userId int64
	userId, _, err = rt.db.SearchByUsername(query)
	if err != nil {
		// could not get the userId, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	var list []string
	_, list, _, err = rt.db.GetFollowing(userId, -1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 401
		err = json.NewEncoder(w).Encode("Unauthorized" + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	
	// get possible query parameters named user
	query := r.URL.Query().Get("userName")
	if query != "" {
		// if the query is empty,  
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	} else {
		// get user names
		followerNames, err = rt.db.GetFollowing(query, -1)
		if err != nil {
			// could not get user names, throw internal server error
			w.WriteHeader(http.StatusInternalServerError) // 500
			err = json.NewEncoder(w).Encode("Internal Server Error" + err.Error())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

	// iterate over the list of comment ids and create comment objects list
	var follower []string
	for _, commentId := range commentIds {
		// get the comment object
		var creator int64
		var content string
		var date string
		creator, _, content, date, err = rt.db.GetComment(commentId)
		if err != nil {
			// could not get comment, throw internal server error
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		// get userName
		var creatorName string
		creatorName, _, err = rt.db.SearchById(creator)
		if err != nil {
			// could not get userName, throw internal server error
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		// create the comment object
		comment := comment{commentId, creatorName, content, date}
		// append the comment object to the list
		comments = append(comments, comment)
	}

	// return the list of user ids
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		// could not encode the list, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
}
