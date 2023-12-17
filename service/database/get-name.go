package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetName(uid string) (string, error) {
	var name string
	err := db.c.QueryRow("SELECT name FROM example_table WHERE id=?", uid).Scan(&name)
	return name, err
}
