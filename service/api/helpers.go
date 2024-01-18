package api

// in order not to write hard stuff
import (
	"errors"

	// "fmt"
	"strconv"
	// "github.com/Lolpe02/wasaphoto-project/service/api/reqcontext"
	// "github.com/julienschmidt/httprouter"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func readPath(ps httprouter.Params, field string) (Id int64, err error) {
	Id, err = strconv.ParseInt(ps.ByName(field), 10, 64)
	if err != nil {
		return -1, errors.New("no field found")
	}
	return
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
func isValid(name string) bool {
	// check substring
	var valid = true
	swears := []string{"dick", "cock", "bitch", "ass", "pussy", "fuck", "stupid"}
	for _, swear := range swears {
		if strings.Contains(name, swear) {
			valid = false
			break
		}
	}
	return valid && 3 <= len(name) && len(name) <= 20
}
