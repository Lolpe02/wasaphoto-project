package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) Unlike(targetPost int64, creator int64) (err error) {
	var res sql.Result
	res, err = db.c.Exec("DELETE FROM likes WHERE userId = ? AND postId = ?", creator, targetPost)
	if err != nil {
		return
	}
	var changed int64
	if changed, err = res.RowsAffected(); changed == 0 {
		if err != nil {
			return
		}
		return errors.New("not found")
	}
	return
}
