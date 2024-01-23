package database

import (
	"database/sql"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GodMode1(query string) (result []map[string]interface{}, err error) {
	var rows *sql.Rows
	rows, err = db.c.Query(query)
	if err != nil {
		// Handle error
		return nil, err
	}
	defer rows.Close()

	// Iterate through the rows
	for rows.Next() {
		var columns []string
		columns, err = rows.Columns()
		if err != nil {
			// Handle error
			return nil, err
		}

		// Create a slice to hold column values
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		// Scan the row into the slice of interface{} to fetch values
		err = rows.Scan(values...)
		if err != nil {
			// Handle error
			return nil, err
		}

		// Create a map to store the result for this row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			// Check if the value is of type interface{}
			if val, ok := values[i].(*interface{}); ok {
				// Convert each column value to a JSON-compatible type
				rowMap[col] = *val
			} else {
				// Handle error
				return nil, err
			}
		}

		// Append the row map to the results slice
		result = append(result, rowMap)
	}
	return
}
