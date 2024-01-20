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
		return
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {
		var userId int64

		// Scan the Id values from each row into variables

		if rowerr := rows.Scan(&userId); rowerr != nil {
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
