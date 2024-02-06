package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
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
		if err.Error() == NotFound {
			w.WriteHeader(http.StatusNotFound) // 404
			err = json.NewEncoder(w).Encode("user not found " + username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError) // 500
			err = json.NewEncoder(w).Encode("couldnt search by useranem " + username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
			}
		}
		return
	}
	// you can see this person's profile only if you're authenticated and not banned
	var yourId int64
	yourId, err = extractToken(r)
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	// check if you're banned
	var banned bool
	_, banned, err = rt.db.GetBanneds(targetId, yourId)
	if err != nil {
		// could not check if youre banned, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("couldnt check if youre banned, " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if banned {
		// you're banned, throw forbidden
		w.WriteHeader(http.StatusForbidden) // 403
		err = json.NewEncoder(w).Encode("you're banned by " + username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	var selname, sub string
	selname, sub, err = rt.db.SearchById(targetId)
	if err != nil {
		// could not get likes, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("couldnt search by id, " + err.Error())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	postIds, err = rt.db.GetProfile(targetId)
	if err != nil {
		// could not get likes, throw internal server error
		w.WriteHeader(http.StatusInternalServerError) // 500
		err = json.NewEncoder(w).Encode("couldnt search profile")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	user := user{targetId, selname, sub, postIds}

	// return the list of post ids of that user
	w.WriteHeader(http.StatusOK) // 200
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}

}
