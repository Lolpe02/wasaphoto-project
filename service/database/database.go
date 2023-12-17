/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"github.com/gofrs/uuid"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	//GetName(username string) (string, error)
	CreateUser(name string) (int64, error)
	//CreateComment(comment string, uid string) error
	//CreateLike(uid string) error
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
	uuidGen *uuid.Gen
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB, genId *uuid.Gen) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure

	var users string
	err1 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&users)
	if errors.Is(err1, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE users (userId INTEGER PRIMARY KEY, userName STRING NOT NULL, date TIMESTAMP NOT NULL);`
		_, err1 = db.Exec(sqlStmt)
		if err1 != nil {
			return nil, fmt.Errorf("error creating database structure table users: %w", err1)
		}
	}
	var images string
	err2 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='images';`).Scan(&images)
	if errors.Is(err2, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE images (postId INTEGER PRIMARY KEY, userId INTEGER, image BLOB NOT NULL, time TIMESTAMP NOT NULL, FOREIGN KEY (userId) REFERENCES users(userId));`
		_, err2 = db.Exec(sqlStmt)
		if err2 != nil {
			return nil, fmt.Errorf("error creating database structure table images: %w", err2)
		}
	}
	var likes string
	err3 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='likes';`).Scan(&likes)
	if errors.Is(err3, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE likes (userId INTEGER, postId INTEGER, PRIMARY KEY (userId, postId), FOREIGN KEY (userId) REFERENCES users(userId), FOREIGN KEY (postId) REFERENCES images(postId));`
		_, err3 = db.Exec(sqlStmt)
		if err3 != nil {
			return nil, fmt.Errorf("error creating database structure table likes: %w", err3)
		}
	}

	var comments string
	err4 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments';`).Scan(&comments)
	if errors.Is(err4, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE comments (commentId INTEGER PRIMARY KEY, userId INTEGER, postId INTEGER, comment TEXT NOT NULL, date TIMESTAMP NOT NULL, FOREIGN KEY (userId) REFERENCES users(userId), FOREIGN KEY (postID) REFERENCES images(postId));`
		_, err4 = db.Exec(sqlStmt)
		if err4 != nil {
			return nil, fmt.Errorf("error creating database structure table comments: %w", err4)
		}
	}
	var bans string
	err5 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='bans';`).Scan(&bans)
	if errors.Is(err5, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE bans (banning STRING, banned STRING, PRIMARY KEY (banning, banned), FOREIGN KEY (banning) REFERENCES users(userId), FOREIGN KEY (banned) REFERENCES users(userId));`
		_, err5 = db.Exec(sqlStmt)
		if err5 != nil {
			return nil, fmt.Errorf("error creating database structure table bans: %w", err5)
		}
	}
	var follows string
	err6 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='follows';`).Scan(&follows)
	if errors.Is(err6, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE follows (following STRING, followed STRING, PRIMARY KEY (following, followed), FOREIGN KEY (following) REFERENCES users(userId), FOREIGN KEY (followed) REFERENCES users(userId));`
		_, err6 = db.Exec(sqlStmt)
		if err6 != nil {
			return nil, fmt.Errorf("error creating database structure table follows: %w", err6)
		}
	}
	return &appdbimpl{
		c: db,
		uuidGen: genId,
	}, nil
}
// Create a new user and return their Id
func (db *appdbimpl) CreateUser(username string) (int64, error) {
	newUId, error := db.uuidGen.NewV4()
	if error != nil {
		fmt.Println(newUId.String())
	}
	var newUserID int64
	err := db.c.QueryRow("INSERT INTO users (userId, userName) VALUES (?, ?) RETURNING id", nil, username).Scan(&newUserID)
	if err != nil {
		fmt.Println("Error inserting into database:", err)
		return -1, err
	}
	return newUserID, nil
}

// This function searches for a specific user in the database given its username.
// It retruns the user if present and ane error.
func (db *appdbimpl) SearchByUsername(targetUser string) (selUserId string, err error) {
	err = db.c.QueryRow("SELECT user_id FROM users WHERE username = ?;", targetUser).Scan(&selUserId)

	// Handling sql.ErrNoRows
	if err != nil {
		if err == sql.ErrNoRows {
			// Return nil values and specify the error
			return "", errors.New("user not found")
		}
		// Return other errors as is
		return "",  err
	}

	// Return retrieved values and nil error
	return selUserId, nil
}
// This function searches for a specific user in the database given its username.
// It retruns the user if present and ane error.
func (db *appdbimpl) SearchById(targetUserId int64) (selUserName string, err error) {
	err = db.c.QueryRow("SELECT user_id, user_date FROM users WHERE username = ?;", targetUserId).Scan(&selUserName)

	// Handling sql.ErrNoRows
	if err != nil {
		if err == sql.ErrNoRows {
			// Return nil values and specify the error
			return "", errors.New("user not found")
		}
		// Return other errors as is
		return "", err
	}

	// Return retrieved values and nil error
	return selUserName, nil
}
func (db *appdbimpl) GetLikes(targetPost int) (userIds []int64, err error) {
	rows, err := db.c.Query("SELECT userId FROM likes WHERE postId = ?;", targetPost)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS countres
		fmt.Println("Error querying database:", err)
		return
	}
	defer rows.Close()
	
	var idList []int64
	fmt.Println("count of likes", rows)
	// Iterate through the rows retrieved
	for rows.Next() {
		var userId int64

		// Scan the Id values from each row into variables
		
		if rowerr := rows.Scan(&userId); rowerr != nil {
			fmt.Println("Error scanning row for likes:", rowerr)
			return  
		}

		// Append the retrieved Id to the list
		idList = append(idList, userId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		fmt.Println("Error during row iteration:", err)
		return
	}

	// Print or use the retrieved Id list
	return idList, err
}

func (db *appdbimpl) GetComments(targetPost int) (commentIds []int64, err error) {
	rows, err := db.c.Query("SELECT commentId FROM comments WHERE postId = ?;", targetPost)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS countres
		fmt.Println("Error querying database:", err)
		return
	}
	defer rows.Close()
	
	var idList []int64
	fmt.Println("count of likes", rows)
	// Iterate through the rows retrieved
	for rows.Next() {
		var commentId int64

		// Scan the Id values from each row into variables
		
		if rowerr := rows.Scan(&commentId); rowerr != nil {
			fmt.Println("Error scanning row for likes:", rowerr)
			return  
		}

		// Append the retrieved Id to the list
		idList = append(idList, commentId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		fmt.Println("Error during row iteration:", err)
		return
	}

	// Print or use the retrieved Id list
	return idList, err
}
func (db *appdbimpl) CreatePost(image []byte, userId int64, date string) (postId int64, err error) {
	// Insert the image into the database
	result, err := db.Exec("INSERT INTO images(data) VALUES(?)", imageData)
	if err != nil {
		fmt.Println("Error inserting image into database:", err)
		return
	}

	// Get the ID of the inserted image
	imageID, _ := result.LastInsertId()
	fmt.Println("Image inserted with ID:", imageID)

	// Retrieve the image from the database
	var retrievedImage []byte
	err = db.QueryRow("SELECT data FROM images WHERE id = ?", imageID).Scan(&retrievedImage)
	if err != nil {
		fmt.Println("Error retrieving image from database:", err)
		return
	}

	// Write the retrieved image data to a new file
	err = ioutil.WriteFile("retrieved_image.jpg", retrievedImage, os.ModePerm)
	if err != nil {
		fmt.Println("Error writing retrieved image file:", err)
		return
	}

	fmt.Println("Image retrieved and saved as retrieved_image.jpg")
}
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
