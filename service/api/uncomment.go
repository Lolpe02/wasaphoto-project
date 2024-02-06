package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncomment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	// take token from the header
	creator, err := extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}

	// take parameters from the path and turn string to int64
	var commentId int64
	commentId, err = readPath(ps, "commentId")
	if err != nil {
		// could not parse the post id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}

	var postId int64
	postId, err = readPath(ps, "postId")
	if err != nil {
		// could not parse the post id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		return
	}
	// check if comment is owned by the creator
	var creatorId int64
	var postOfComment int64
	creatorId, postOfComment, _, _, err = rt.db.GetComment(commentId)
	if err != nil {
		if err.Error() == NotFound {
			// could not find the post, throw not found
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}
		// throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if creatorId != creator {
		// not authorized, throw forbidden
		w.WriteHeader(http.StatusForbidden) // 403
		err = json.NewEncoder(w).Encode("not authorized to delete this comment, its not your's")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if postOfComment != postId {
		// impossible something went wrong
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("WHAT HAPPENED? Comment does not belong to the post")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// put the comment in the database
	err = rt.db.Uncomment(creator, postId, commentId)
	if err != nil {
		if err.Error() == NotFound {
			// could not find the post, throw not found
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}
		// could not create comment, internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	// return the id of the comment?? idk
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode("comment deleted")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
