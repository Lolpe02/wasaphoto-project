package api

// in order not to write hard stuff
import (
	"errors"

	//"fmt"
	"strconv"
	//"github.com/Lolpe02/wasaphoto-project/service/api/reqcontext"
	//"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func readPath(ps httprouter.Params, field string) (Id int64, err error) {
	Id, err = strconv.ParseInt(ps.ByName(field), 10, 64)
	if err != nil {
		return -1, errors.New("no field found")
	}
	return Id, nil
}

func extractToken(r *http.Request) (bearer int64, err error) {
	bearerToken := r.Header.Get("Authorization")
	bearerToken = strings.TrimPrefix(bearerToken, ": ")
	if bearerToken == "" || bearerToken == " " {
		// the request body was not a parseable JSON or is missing, rejecting the request and return error
		return -1, errors.New("no authorization header")
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
	return -1, errors.New("something wrong with authorization header")
}
func isValid(str string) bool {
	return 3 <= len(str) && len(str) <= 20
}
