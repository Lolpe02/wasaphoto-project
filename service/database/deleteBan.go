package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) UnbanUser(yourId int64, theirId int64) (err error) {
	var res sql.Result
	res, err = db.c.Exec("DELETE FROM bans WHERE banning = ? AND banned = ?;", yourId, theirId)
	if err != nil {
		return err
	}
	var rowsAffected int64
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New(NotFound)
	}

	return
}
