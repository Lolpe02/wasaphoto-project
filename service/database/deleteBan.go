package database

import (
// "github.com/mattn/go-sqlite3"
// "strings"
)

func (db *appdbimpl) UnbanUser(yourId int64, theirId int64) (err error) {
	_, err = db.c.Exec("DELETE FROM bans WHERE banning = ? AND banned = ?;", yourId, theirId)
	if err != nil {
		return err
	}
	return
}
