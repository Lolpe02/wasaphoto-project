package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetMetadata(postId int64) (userId int64, description string, date string, err error) {
	// Retrieve the image from the database
	err = db.c.QueryRow("SELECT userId, description, time FROM images WHERE postId = ?;", postId).Scan(&userId, &description, &date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return nil values and specify the error
			return -1, "", "", errors.New(NotFound)
		}
		return -1, "", "", err
	}
	return
}
