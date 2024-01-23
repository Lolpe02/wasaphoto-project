package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetLikes(targetPost int64) (userIds []int64, err error) {
	var rows *sql.Rows
	rows, err = db.c.Query("SELECT userId FROM likes WHERE postId = ?;", targetPost)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS countres
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {

		if err = rows.Err(); err != nil {
			return
		}
		// Scan the Id values from each row into variables
		var userId int64
		if err = rows.Scan(&userId); err != nil {
			return
		}

		// Append the retrieved Id to the list
		userIds = append(userIds, userId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return
	}

	// Print or use the retrieved Id list
	return
}
