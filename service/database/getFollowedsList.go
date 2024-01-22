package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

// GetFolloweds retrieves the list of users followed by the target user.
// It takes the target user Id and a test Id as parameters.
// It returns the list of followed user Ids, a boolean indicating if the test Id is present in the list,
// and an error if any occurred during the database query or iteration.
func (db *appdbimpl) GetFolloweds(targetUserId int64, testId int64) (followedbyTargetIds []int64, present bool, err error) {
	var rows *sql.Rows
	present = false
	rows, err = db.c.Query("SELECT followed FROM follows WHERE following = ?;", targetUserId)

	if err != nil {
		return
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {
		var followedId int64

		// Scan the Id values from each row into variables

		if rowerr := rows.Scan(&followedId); rowerr != nil {
			return nil, false, err
		}
		if followedId == testId {
			present = true
		}
		// Append the retrieved Id to the list
		followedbyTargetIds = append(followedbyTargetIds, followedId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return
	}

	// Print or use the retrieved Id list
	return
}
