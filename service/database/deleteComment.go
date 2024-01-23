package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) Uncomment(creatorId int64, postId int64, commentId int64) (err error) {
	var res sql.Result
	res, err = db.c.Exec("DELETE FROM comments WHERE userId = ? AND commentId = ? AND postId = ?;", creatorId, commentId, postId)
	if err != nil {
		return err
	}
	var changed int64
	if changed, err = res.RowsAffected(); changed == 0 {
		if err != nil {
			return err
		}
		return errors.New(NotFound)
	}
	return
}
