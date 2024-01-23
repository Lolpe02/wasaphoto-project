package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) ChangeUsername(yourUserId int64, newUsername string) (err error) {
	var res sql.Result
	res, err = db.c.Exec("UPDATE users SET userName = ? WHERE userName != ? AND userId = ?;", newUsername, yourUserId, yourUserId)
	if err != nil {
		// could not update the user, throw internal server error
		return err
	}
	// check if the user was updated
	var rows int64
	rows, err = res.RowsAffected()
	if err != nil {
		// could not update the user, throw internal server error
		return err
	}
	if rows == 0 {
		// the user was not updated, throw bad request
		return errors.New(NotFound)
	}
	return
}
