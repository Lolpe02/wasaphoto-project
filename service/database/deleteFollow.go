package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) UnfollowUser(yourId int64, theirId int64) (err error) {
	var res sql.Result
	res, err = db.c.Exec("DELETE FROM follows WHERE following = ? AND followed = ?", yourId, theirId)
	if err != nil {
		return err
	}
	var aff int64
	aff, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("not following user")
	}
	return
}
