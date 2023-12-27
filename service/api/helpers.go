package api

// in order not to write hard stuff
import (
	"encoding/json"
	"errors"

	//"fmt"
	"strconv"
	//"github.com/Lolpe02/wasaphoto-project/service/api/reqcontext"
	//"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func bodyReader(r *http.Request, args int) (err error) {
	err = json.NewDecoder(r.Body).Decode(&args)
	return
}
func readPath(ps httprouter.Params, field string) (Id int64, err error) {
	Id, err = strconv.ParseInt(ps.ByName(field), 10, 64)
	if err != nil {
		return -1, errors.New("No field found")
	}
	return Id, nil
}

func extractToken(r *http.Request) (bearer int64, err error) {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		// the request body was not a parseable JSON or is missing, rejecting the request and return error
		return -1, errors.New("No Authorization header")
	}
	// Normally Authorization the_token
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		intId, err1 := strconv.ParseInt(strArr[1], 10, 64)
		if err1 != nil {
			return -1, errors.New("ParseInt error")
		}
		return intId, nil
	}
	return -1, errors.New("Something wrong with Authorization header")
}
