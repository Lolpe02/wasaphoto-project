package api

// in order not to write hard stuff
import (
	"errors"
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
func isValid(name string) bool {
	// check substring
	var valid = true
	swears := []string{"dick", "cock", "asshole", "bitch", "pussy"}
	for _, swear := range swears {
		if strings.Contains(name, swear) {
			valid = false
		}
	}
	return valid && 3 <= len(name) && len(name) <= 20
}

// Function to check if the error is a foreign key violation
/*func isForeignKeyViolation(err error) bool {
	// SQLite error codes for foreign key violations
	sqliteForeignKeyErrorCodes := []string{
		"foreign key constraint failed",
		"constraint failed",
		"violated constraint",
		"conflicted with FOREIGN KEY constraint",
		"violates foreign key",
		"duplicate key value violates constraint",
	}

	// Check if the error message contains any of the known SQLite foreign key violation messages
	for _, code := range sqliteForeignKeyErrorCodes {
		if strings.Contains(err.Error(), code) {
			return true
		}
	}
	return false
}*/
