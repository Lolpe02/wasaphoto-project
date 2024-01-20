package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) CreateUser(username string) (yourUserId int64, alreadyExists bool, err error) {
	// newUId, error := db.uuidGen.NewV4()
	// if error != nil {
	// 	fmt.Println(newUId.String())
	// }

	// "INSERT INTO users (userName) VALUES (?) RETURNING userId ON CONFLICT (userName) DO SELECT userId FROM users WHERE userName = ?"

	// search if user already exists
	alreadyExists = true
	err = db.c.QueryRow("SELECT userId FROM users WHERE userName = ?;", username).Scan(&yourUserId)
	if errors.Is(err, sql.ErrNoRows) {
		alreadyExists = false
		// user does not exist, create it
		err = db.c.QueryRow("INSERT INTO users (userName) VALUES (?) RETURNING userId", username).Scan(&yourUserId)
		if err != nil {
			return -1, false, err
		}
	}
	return
}
