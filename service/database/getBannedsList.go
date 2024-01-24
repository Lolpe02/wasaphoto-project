package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetBanneds(targetUserId int64, testId int64) (bannedIds []int64, present bool, err error) {
	var rows *sql.Rows
	present = false
	rows, err = db.c.Query("SELECT banned FROM bans WHERE banning = ?;", targetUserId)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {

		if err = rows.Err(); err != nil {
			return nil, false, err
		}
		// Scan the Id values from each row into variables
		var bannedId int64
		if err = rows.Scan(&bannedId); err != nil {
			return nil, false, err
		}
		if bannedId == testId {
			present = true
		}
		// Append the retrieved Id to the list
		bannedIds = append(bannedIds, bannedId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, false, err
	}

	// Print or use the retrieved Id list
	return
}
