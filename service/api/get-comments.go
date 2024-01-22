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
	if err != nil {
		// could not get likes, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	// iterate over the list of comment ids and create comment objects list
	var comments []comment
	for _, commentId := range commentIds {
		// get the comment object
		creator, _, content, date, errcom := rt.db.GetComment(commentId)
		if errcom != nil {
			// could not get comment, throw internal server error
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		// create the comment object
		comment := comment{postId, creator, content, date}
		// append the comment object to the list
		comments = append(comments, comment)
	}

	// return the list of user ids
	w.WriteHeader(http.StatusOK) // 200
	json.NewEncoder(w).Encode(comments)
}
