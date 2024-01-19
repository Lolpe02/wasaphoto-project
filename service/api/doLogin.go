package api

/*
go run ./cmd/webapi/
curl -v -H 'Content-Type: application/json' -d '{"username": "Lollo"}' localhost:3000/session
*/

/*
Possible outcomes:

1. checking if decoding operation ended successfully
   curl -v -H 'Content-Type: application/json' -d '{"username": "Lollo}' localhost:3000/session
   (the JSON data is missing a closing quote resulting in an invalid JSON structure)

2. checking if the username is valid
   a. curl -v -H 'Content-Type: application/json' -d '{"username": "     "}' localhost:3000/session
      (the client has enterd white spaces only, hence the username is not valid)

   b. (username doesn't match string pattern: '^.*?$': it contains a new line)

   c. curl -v -H 'Content-Type: application/json' -d '{"username": "ab"}' localhost:3000/session
      (username hasn't got required length: is <3 or >16)

3. if the user altready exists, return the ID
   (post an alredy existing username)

	4. if encoding operation is unsuccessful though the user is present
	   (server error)

5. if the user doesn't exist yet, create it and return the ID
   (post a new username)

	6. if user creation or ID retrieval is unsuccessful
	   (server error)

	7. if encoding operation is unsuccessful though the user has been created
	   (server error)
*/

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/Lolpe02/wasaphoto-project/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

/*
Summary: logs in the user

Description:
The login endpoint accepts a username like “Maria” without any password.
If the username does not exist, it will be created, and an identifier is returned.
If the username exists, the user identifier is returned.
*/
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// setting response header

	w.Header().Set("Content-Type", "application/json")

	// extracting username from the request
	var userName string
	err := json.NewDecoder(r.Body).Decode(&userName)

	// 1.
	// checking if decoding operation ended successfully
	if err != nil {
		// the request body was not a parseable JSON or is missing, rejecting the request
		w.WriteHeader(http.StatusBadRequest) //400
		ctx.Logger.WithError(err).Error("doLogin: the request body was not a parseable JSON or is missing")
		fmt.Fprint(w, "\ndoLogin: the request body was not a parseable JSON or is missing\n")
		return
	}
	// validating username (removing white spaces and new lines)
	userName = strings.TrimSpace(strings.Replace(userName, "\n", "", -1))
	// 2.
	// checking if the username is valid
	if !isValid(userName) {
		// the username is not valid, rejecting request
		w.WriteHeader(http.StatusBadRequest) //400
		fmt.Fprint(w, "\ndoLogin: the username is not valid\n\n")
		return
	}
	// 3.
	// return the ID
	var userId int64
	var existed bool
	userId, existed, err = rt.db.CreateUser(userName)

	// 6.
	// if user creation or ID retrieval is unsuccessful
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("doLogin: user creation or ID retrieval is unsuccessful")
		fmt.Fprint(w, "\ndoLogin: user creation or ID retrieval is unsuccessful\n")
		fmt.Fprint(w, "\ndoLogin: user creation or ID retrieval is unsuccessful\n\n")
		return
	}
	// send it back

	// if user didn't exist, return 201, otherwise 200 false:
	if !existed {
		w.WriteHeader(http.StatusCreated) // 201
	}
	//fmt.Fprintln(w)
	err = json.NewEncoder(w).Encode(userId)

	// 4.
	// if encoding operation is unsuccessful though the user is present
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) //500
		ctx.Logger.WithError(err).Error("doLogin: unable to encode JSON response though the user is present")
		_, err = fmt.Fprint(w, "doLogin: unable to encode JSON response though the user is present\n")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError) // 500
			ctx.Logger.WithError(err).Error("doLogin: unable to write response")

			return
		}
		return
	}

	w.WriteHeader(http.StatusOK) //200
	fmt.Fprint(w, "\nUser log-in action successful.\nThe user ID is returned in the content.\n\n")

}
