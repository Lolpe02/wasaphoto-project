package api

// in order not to write hard stuff
import (
	"encoding/json"
	"errors"
	//"fmt"

	//"github.com/Lolpe02/wasaphoto-project/service/api/reqcontext"
	//"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func bodyReader(r http.Request, args int) (err error) {
	err = json.NewDecoder(r.Body).Decode(&args)
	return
}
func ExtractToken(r *http.Request) (bearer string, err error) {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		// the request body was not a parseable JSON or is missing, rejecting the request and return error
		return "", errors.New("No Authorization header")
	}
	// Normally Authorization the_token
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}
	return "", errors.New("Something wrong with Authorization header")
}
