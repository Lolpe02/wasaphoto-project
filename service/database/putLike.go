package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) PutLike(targetPost int64, creator int64) (err error) {
	var res sql.Result
	res, err = db.c.Exec("INSERT OR IGNORE INTO likes (userId, postId) VALUES (?, ?);", creator, targetPost)
	if err != nil {
		return err
	}
	var changed int64
	if changed, err = res.RowsAffected(); changed == 0 {
		if err != nil {
			return err
		}
		return errors.New("already liked")
	}
	return
}
