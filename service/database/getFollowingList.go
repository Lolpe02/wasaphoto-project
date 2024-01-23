package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

// GetFollowing retrieves the list of target user Ids that the specified target user is following,
// along with a boolean indicating whether the test Id is present in the list.
// It takes the target user Id and the test Id as parameters.
// It returns the following target user Ids, a boolean indicating presence, and any error encountered.
func (db *appdbimpl) GetFollowing(targetUserId int64, testId int64) (followingTargetIds []int64, present bool, err error) {
	var rows *sql.Rows
	present = false
	rows, err = db.c.Query("SELECT following FROM follows WHERE followed = ?;", targetUserId)

	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {
		var followingId int64

		// Scan the Id values from each row into variables

		if rowerr := rows.Scan(&followingId); rowerr != nil {
			return nil, false, rowerr
		}
		if followingId == testId {
			present = true
		}
		// Append the retrieved Id to the list
		followingTargetIds = append(followingTargetIds, followingId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, false, err
	}

	// Print or use the retrieved Id list
	return
}
