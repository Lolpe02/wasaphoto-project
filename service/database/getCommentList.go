package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetCommentList(targetPost int64, specificUser int64) (commentIds []int64, err error) {
	var rows *sql.Rows
	switch {
	case specificUser != -1:
		rows, err = db.c.Query("SELECT commentId FROM comments WHERE postId = ? AND userId = ? ORDER BY date DESC;", targetPost, specificUser)
	default:
		rows, err = db.c.Query("SELECT commentId FROM comments WHERE postId = ? ORDER BY date DESC;", targetPost)
	}
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS count
		return
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {
		var commentId int64

		// Scan the Id values from each row into variables

		if err = rows.Scan(&commentId); err != nil {
			return
		}

		// Append the retrieved Id to the list
		commentIds = append(commentIds, commentId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return
	}
	return
}
