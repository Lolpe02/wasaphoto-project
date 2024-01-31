package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

// GetFollowing retrieves the list of people ids and names who are following target user Id,
// along with a boolean indicating whether the test Id is present in the list.
// It takes the target user Id and the test Id as parameters.
func (db *appdbimpl) GetFollowing(targetUserId int64, testId int64) (followingTargetIds []int64, followingTargetNames []string, present bool, err error) {
	var rows *sql.Rows
	present = false
	rows, err = db.c.Query("SELECT follows.following, users.userName FROM follows JOIN users ON follows.following = users.userId WHERE follows.followed = ?;", targetUserId)

	if err != nil {
		return nil, nil, false, err
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {

		if err = rows.Err(); err != nil {
			return
		}

		// Scan the Id values from each row into variables
		var followingId int64
		var followingName string
		if err = rows.Scan(&followingId, &followingName); err != nil {
			return
		}
		if followingId == testId {
			present = true
		}
		// Append the retrieved Id to the list
		followingTargetIds = append(followingTargetIds, followingId)
		followingTargetNames = append(followingTargetNames, followingName)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, nil, false, err
	}

	// Print or use the retrieved Id list
	return
}
