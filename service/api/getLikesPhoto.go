package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"encoding/json"
)

// getHelloWorld is an example of HTTP endpoint that returns "Hello world!" as a plain text
func (rt *_router) getLikesPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write([]byte("Hello World!"))
}
