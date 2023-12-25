package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (rt *_router) getLikesPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	var userIds []int64
	// take parameters from the path and turn string to int64
	postId, err := strconv.ParseInt(ps.ByName("postId"), 10, 64)

	if err != nil {
		// could not parse the post id, throw bad request
		w.WriteHeader(http.StatusBadRequest) //400
	}
	userIds, err = rt.db.GetLikes(postId)
	if err != nil {
		// could not get likes, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) //500
		return
	}
	// return the list of user ids
	w.WriteHeader(http.StatusOK) //200
	json.NewEncoder(w).Encode(userIds)
	return
}
