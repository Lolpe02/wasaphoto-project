package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) SearchById(targetUserId int64) (selUserName string, subscription string, err error) {
	err = db.c.QueryRow("SELECT userName, date FROM users WHERE userId = ?;", targetUserId).Scan(&selUserName, &subscription)

	if err != nil {
		// Handling sql.ErrNoRows
		if errors.Is(err, sql.ErrNoRows) {
			// Return nil values and specify the error
			return "", "", errors.New(NotFound)
		}
		// Return other errors as is
		return "", "", err
	}

	// Return retrieved values and nil error
	return
}
