package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) SearchByUsername(targetUser string) (selUserId int64, err error) {
	err = db.c.QueryRow("SELECT userId FROM users WHERE userName = ?;", targetUser).Scan(&selUserId)

	// Handling sql.ErrNoRows
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return nil values and specify the error
			return -1, errors.New(NotFound)
		}
		// Return other errors as is
		return -1, err
	}

	// Return retrieved values and nil error
	return
}
