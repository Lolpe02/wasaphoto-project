package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	var commentIds []int64
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
	if !follows && authUserId != postCreator {
		w.WriteHeader(http.StatusUnauthorized) // 401
		err = json.NewEncoder(w).Encode("you dont follow the owner")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	} */
	// get possible query parameters named user
	query := r.URL.Query().Get("commenter")
	if query != "" {
		// if the query is not empty, parse it to int64
		commenter, err := strconv.ParseInt(query, 10, 64)
		if err != nil {
			// could not parse the user id, throw bad request
			w.WriteHeader(http.StatusBadRequest) // 400
			return
		}
		// get the list of comment ids
		commentIds, err = rt.db.GetCommentList(postId, commenter)
		if err != nil {
			// could not get likes, throw internal server error
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
	} else {
		// get the list of comment ids
		commentIds, err = rt.db.GetCommentList(postId, -1)
		if err != nil {
			// could not get likes, throw internal server error
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
	}

	// iterate over the list of comment ids and create comment objects list
	var comments []comment
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
