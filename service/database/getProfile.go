package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetProfile(targetUserId int64) (idList []int64, err error) {
	var rows *sql.Rows
	rows, err = db.c.Query("SELECT postId FROM images WHERE userId = ? ORDER BY time DESC", targetUserId)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS count
		return
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {

		if err = rows.Err(); err != nil {
			return
		}

		// Scan the Id values from each row into variables
		var postId int64
		if err = rows.Scan(&postId); err != nil {
			return
		}

		// Append the retrieved Id to the list
		idList = append(idList, postId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Print or use the retrieved Id list
	return
}
