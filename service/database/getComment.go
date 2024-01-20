package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetComment(commentId int64) (creator int64, content string, date string, err error) {
	err = db.c.QueryRow("SELECT userId, comment, date FROM comments WHERE commentId = ?;", commentId).Scan(&creator, &content, &date)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS countres
		if errors.Is(err, sql.ErrNoRows) {
			// Return nil values and specify the error
			return -1, "", "", errors.New("not found")
		}
		return -1, "", "", err
	}
	return
}
