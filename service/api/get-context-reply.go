package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getContextReply is an example of HTTP endpoint that returns "Hello World!" as a plain text. The signature of this
// handler accepts a reqcontext.RequestContext (see httpRouterHandler).
func (rt *_router) getContextReply(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("Hello World!\n"))
	if err != nil {
		http.Error(w, "getContextReply: error writing response", http.StatusInternalServerError)
	}

}
