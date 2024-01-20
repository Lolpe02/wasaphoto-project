package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GodMode2(query string) (result int64, err error) {
	var res sql.Result
	res, err = db.c.Exec(query)
	if err != nil {
		// Handle error
		return -1, err
	}
	result, err = res.RowsAffected()
	if err != nil {
		// Handle error
		return -1, err
	}
	return
}
