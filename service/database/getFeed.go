package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetFeed(yourId int64) (postIds []int64, err error) {
	var rows *sql.Rows
	rows, err = db.c.Query("SELECT images.postId FROM images JOIN follows ON images.userId = follows.followed WHERE follows.following = ? ORDER BY images.time DESC;", yourId)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS count
		return
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {
		
		if err = rows.Err(); err != nil {
			return nil, err
		}

		// Scan the Id values from each row into variables
		var postId int64
		if err = rows.Scan(&postId); err != nil {
			return nil, err
		}

		// Append the retrieved Id to the list
		postIds = append(postIds, postId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return
	}

	return
}
