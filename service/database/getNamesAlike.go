package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetInfo(targetUser string) (userIds []int64, usernamesAlike []string, err error) {
	var rows *sql.Rows
	rows, err = db.c.Query(`SELECT userId, userName FROM users WHERE userName LIKE '%'||?||'%'`, targetUser)
	// Handling sql.ErrNoRows
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return nil values and specify the error
			return nil, nil, errors.New(NotFound)
		}
		// Return other errors as is
		return nil, nil, err
	}

	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {

		if err = rows.Err(); err != nil {
			return
		}
		// Scan the Id values from each row into variables
		var userId int64
		var userName string
		if err = rows.Scan(&userId, &userName); err != nil {
			return
		}

		// Append the retrieved Id to the list
		userIds = append(userIds, userId)
		usernamesAlike = append(usernamesAlike, userName)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return
	}
	return
}
