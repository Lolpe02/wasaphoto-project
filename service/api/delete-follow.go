package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")

	// take id parameters from the path (person to unfollow)
	IdtounFollow, err := readPath(ps, "followId")
	if err != nil {
		// could not parse the id, throw bad request
		w.WriteHeader(http.StatusBadRequest) // 400
		_, err = w.Write([]byte(err.Error() + ", you fool"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		return
	}

	var yourId int64
	if err != nil {
		// not authenticated, throw unauthorized
		w.WriteHeader(http.StatusUnauthorized) // 401
		return
	}
	err = rt.db.UnfollowUser(yourId, IdtounFollow)
	if err != nil {
		if err.Error() == "not following user" {
			w.WriteHeader(http.StatusNotFound) // 404
			/*_, err1 := w.Write([]byte(err.Error() + ", you fool"))
			if err1 != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				return
			}*/
			return
		} else {
			// could not unfollow, throw internal server error
			w.WriteHeader(http.StatusInternalServerError) // 500
			_, err2 := w.Write([]byte(err.Error() + ", sopmething went wrong"))
			if err2 != nil {
				w.WriteHeader(http.StatusInternalServerError) // 500
				return
			}
			return
		}
	}
	_, err = w.Write([]byte("\nUnfollowed(:\n"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
}
