package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

// receive creator from bearer token, photo id from path, write
func (rt *_router) like(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// read security bearer token from header
	likeCreator, err := extractToken(r)
	if err != nil {
		// bad request
		err = json.NewEncoder(w).Encode(http.StatusUnauthorized) // 403
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	// Parse the path parameter for the photo id
	var pid int64
	pid, err = readPath(ps, "postId")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var yourId int64
	yourId, err = readPath(ps, "yourId")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode("Unauthorized" + err.Error()) // 401
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if yourId != likeCreator {
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode("something went wrong with your id and token")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	// check if user follows creator
	var postCreator int64
	postCreator, _, _, err = rt.db.GetMetadata(pid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode("couldnt get post owner")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	var follows bool
	_, _, follows, err = rt.db.GetFollowing(postCreator, likeCreator)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode("Unauthorized" + err.Error()) // 401
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	if !follows || postCreator == likeCreator {
		w.WriteHeader(http.StatusUnauthorized)
		err = json.NewEncoder(w).Encode("you dont follow or you're trying to like your own post, pityful") // 401
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	like := like{likeCreator, pid}
	// Call the database function to create the like
	err = rt.db.PutLike(like.PostId, like.UserId)
	if err != nil {
		if err.Error() == "already liked" {
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode("already liked")
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		log.Println(err)
		err = json.NewEncoder(w).Encode("Failed to create like")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Return a success response
	w.WriteHeader(http.StatusCreated) // 201
	err = json.NewEncoder(w).Encode("like created")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
