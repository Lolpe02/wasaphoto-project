package database

import (
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) BanUser(yourId int64, theirId int64) (err error) {
	_, ban, err := db.GetBanneds(yourId, theirId)
	if err != nil {
		return err
	}
	switch {
	case ban:
		return errors.New("already banned this user")
	default:
		_, err = db.c.Exec("INSERT INTO bans (banning, banned) VALUES (?, ?);", yourId, theirId)
	}
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM follows WHERE followed = ? AND following = ?;", yourId, theirId)

	if err != nil {
		return err
	}
	return
}
