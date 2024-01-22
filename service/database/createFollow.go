package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) FollowUser(yourId int64, theirId int64) (alreadyExists bool, err error) {
	// check if user exists
	_, _, err = db.SearchById(theirId)
	if err != nil {
		return false, err
	}
	var banned bool
	_, banned, err = db.GetBanneds(theirId, yourId)
	if err != nil {
		return false, err
	}
	var res sql.Result
	switch {
	case banned:
		return true, errors.New("banned by this user")
	default:
		res, err = db.c.Exec("INSERT OR IGNORE INTO follows (following, followed) VALUES (?, ?);", yourId, theirId)
	}
	if err != nil {
		/*if strings.Contains(err.Error(), "FOREIGN KEY") {
			return true, nil
		}*/
		return true, err
	}
	var changed int64
	if changed, err = res.RowsAffected(); changed == 0 {
		if err != nil {
			return true, err
		}
		alreadyExists = true
		return
	}
	alreadyExists = false
	return
}
