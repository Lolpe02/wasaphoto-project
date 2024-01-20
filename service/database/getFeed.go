package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetFeed(yourId int64) (postIds []int64, err error) {
	var rows *sql.Rows
	rows, err = db.c.Query("SELECT postId FROM images JOIN follows ON images(userId) = followed WHERE following = ? ORDER BY time DESC", yourId)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS count
		return
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {
		var postId int64

		// Scan the Id values from each row into variables

		if rowerr := rows.Scan(&postId); rowerr != nil {
			return
		}

		// Append the retrieved Id to the list
		postIds = append(postIds, postId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return
	}

	// Print or use the retrieved Id list
	return
}
