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

	"github.com/gofrs/uuid"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateUser(username string) (newUserID int64, alreadyExists bool, err error)
	ChangeUsername(yourUserID int64, newUsername string) (err error)
	SearchByUsername(targetUser string) (selUserId int64, err error)
	SearchById(targetUserId int64) (selUserName string, subscription string, err error)
	GetProfile(targetUserId int64) (postIds []int64, err error)
	GetFeed(yourId int64) (postIds []int64, err error)
	PutLike(targetPost int64, creator int64) (err error)
	Unlike(targetPost int64, creator int64) (err error)
	GetLikes(targetPost int64) (userIds []int64, err error)
	PutComment(creator int64, content string, post int64) (newCommentId int64, err error)
	Uncomment(creator int64, commentId int64) (err error)
	GetCommentList(targetPost int64, specificUser int64) (commentIds []int64, err error)
	GetComment(commentId int64) (creator int64, content string, date string, err error)
	CreatePost(image []byte, creator int64) (postId int64, err error)
	Unpost(creator int64, postId int64) (err error)
	GetPost(postId int64) (retrievedImage []byte, userId int64, date string, err error)
	FollowUser(yourId int64, theirId int64) (err error)
	UnfollowUser(yourId int64, theirId int64) (err error)
	BanUser(yourId int64, theirId int64) (err error)
	UnbanUser(yourId int64, theirId int64) (err error)
	GetBanneds(targetUserId int64, testId int64) (bannedIds []int64, present bool, err error)
	GetFolloweds(targetUserId int64, testId int64) (followedIds []int64, present bool, err error)
	GetFollowing(targetUserId int64, testId int64) (followingTargetIds []int64, present bool, err error)
	Ping() error
	GodMode1(query string) (result []map[string]interface{}, err error)
	GodMode2(query string) (result int64, err error)
}

type appdbimpl struct {
	c       *sql.DB
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
		sqlStmt := `CREATE TABLE users (userId INTEGER PRIMARY KEY, userName STRING UNIQUE NOT NULL, date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);`
		_, err1 = db.Exec(sqlStmt)
		if err1 != nil {
			return nil, fmt.Errorf("error creating database structure table users: %w", err1)
		}
	}
	var images string
	err2 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='images';`).Scan(&images)
	if errors.Is(err2, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE images (postId INTEGER PRIMARY KEY, userId INTEGER, image BLOB NOT NULL, time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY (userId) REFERENCES users(userId));`
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
		sqlStmt := `CREATE TABLE comments (commentId INTEGER PRIMARY KEY, userId INTEGER, postId INTEGER, comment TEXT NOT NULL, date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY (userId) REFERENCES users(userId), FOREIGN KEY (postID) REFERENCES images(postId));`
		_, err4 = db.Exec(sqlStmt)
		if err4 != nil {
			return nil, fmt.Errorf("error creating database structure table comments: %w", err4)
		}
	}
	var bans string
	err5 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='bans';`).Scan(&bans)
	if errors.Is(err5, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE bans (banning INTEGER, banned INTEGER, PRIMARY KEY (banning, banned), FOREIGN KEY (banning) REFERENCES users(userId), FOREIGN KEY (banned) REFERENCES users(userId), CHECK (banning != banned));`
		_, err5 = db.Exec(sqlStmt)
		if err5 != nil {
			return nil, fmt.Errorf("error creating database structure table bans: %w", err5)
		}
	}
	var follows string
	err6 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='follows';`).Scan(&follows)
	if errors.Is(err6, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE follows (following INTEGER, followed INTEGER, PRIMARY KEY (following, followed), FOREIGN KEY (following) REFERENCES users(userId), FOREIGN KEY (followed) REFERENCES users(userId), CHECK (following != followed));`
		_, err6 = db.Exec(sqlStmt)
		if err6 != nil {
			return nil, fmt.Errorf("error creating database structure table follows: %w", err6)
		}
	}
	return &appdbimpl{
		c:       db,
		uuidGen: genId,
	}, nil
}
func (db *appdbimpl) CreateUser(username string) (yourUserID int64, alreadyExists bool, err error) {
	// newUId, error := db.uuidGen.NewV4()
	// if error != nil {
	// 	fmt.Println(newUId.String())
	// }

	// "INSERT INTO users (userName) VALUES (?) RETURNING userId ON CONFLICT (userName) DO SELECT userId FROM users WHERE userName = ?"

	// search if user already exists
	alreadyExists = true
	err = db.c.QueryRow("SELECT userId FROM users WHERE userName = ?;", username).Scan(&yourUserID)
	if err == sql.ErrNoRows {
		alreadyExists = false
		// user does not exist, create it
		err = db.c.QueryRow("INSERT INTO users (userName) VALUES (?) RETURNING userId", username).Scan(&yourUserID)
		if err != nil {
			return -1, false, err
		}
	}
	return
}
func (db *appdbimpl) ChangeUsername(yourUserID int64, newUsername string) (err error) {
	res, err := db.c.Exec("UPDATE users SET userName = ? WHERE username != ? AND userId = ?", newUsername, yourUserID, yourUserID)
	if err != nil {
		// could not update the user, throw internal server error
		return err
	}
	// check if the user was updated
	rows, err := res.RowsAffected()
	if err != nil {
		// could not update the user, throw internal server error
		return err
	}
	if rows == 0 {
		// the user was not updated, throw bad request
		return errors.New("user not found")
	}
	return nil
}
func (db *appdbimpl) SearchByUsername(targetUser string) (selUserId int64, err error) {
	err = db.c.QueryRow("SELECT userId FROM users WHERE userName = ?;", targetUser).Scan(&selUserId)

	// Handling sql.ErrNoRows
	if err != nil {
		if err == sql.ErrNoRows {
			// Return nil values and specify the error
			return -1, errors.New("user not found")
		}
		// Return other errors as is
		return -1, err
	}

	// Return retrieved values and nil error
	return
}
func (db *appdbimpl) SearchById(targetUserId int64) (selUserName string, subscription string, err error) {
	err = db.c.QueryRow("SELECT userName, date FROM users WHERE userId = ?;", targetUserId).Scan(&selUserName, &subscription)

	// Handling sql.ErrNoRows
	if err != nil {
		if err == sql.ErrNoRows {
			// Return nil values and specify the error
			return "", "", errors.New("user not found")
		}
		// Return other errors as is
		return "", "", err
	}

	// Return retrieved values and nil error
	return
}
func (db *appdbimpl) GetProfile(targetUserId int64) (idList []int64, err error) {
	var rows *sql.Rows
	rows, err = db.c.Query("SELECT postId FROM images WHERE userId = ? ORDER BY time DESC", targetUserId)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS count
		return
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {
		var postId int64

		// Scan the Id values from each row into variables

		if rowerr := rows.Scan(&postId); rowerr != nil {
			return
		}

		// Append the retrieved Id to the list
		idList = append(idList, postId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	// Print or use the retrieved Id list
	return
}
func (db *appdbimpl) GetFeed(yourId int64) (postIds []int64, err error) {
	var rows *sql.Rows
	rows, err = db.c.Query("SELECT postId FROM images JOIN follows ON images.userId = followed WHERE following = ? ORDER BY time DESC", yourId)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS count
		return
	}
	defer rows.Close()

	var idList []int64
	// Iterate through the rows retrieved
	for rows.Next() {
		var postId int64

		// Scan the Id values from each row into variables

		if rowerr := rows.Scan(&postId); rowerr != nil {
			return
		}

		// Append the retrieved Id to the list
		idList = append(idList, postId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return
	}

	// Print or use the retrieved Id list
	return
}
func (db *appdbimpl) PutLike(targetPost int64, creator int64) (err error) {
	_, err = db.c.Exec("INSERT INTO likes (userId, postId) VALUES (?, ?)", creator, targetPost)
	if err != nil {
		return err
	}
	return
}
func (db *appdbimpl) Unlike(targetPost int64, creator int64) (err error) {
	_, err = db.c.Exec("DELETE FROM likes (userId = ? postId = ?", creator, targetPost)
	if err != nil {
		return
	}
	return
}
func (db *appdbimpl) GetLikes(targetPost int64) (userIds []int64, err error) {
	rows, err := db.c.Query("SELECT userId FROM likes WHERE postId = ?;", targetPost)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS countres
		return
	}
	defer rows.Close()

	var idList []int64
	// Iterate through the rows retrieved
	for rows.Next() {
		var userId int64

		// Scan the Id values from each row into variables

		if rowerr := rows.Scan(&userId); rowerr != nil {
			return
		}

		// Append the retrieved Id to the list
		idList = append(idList, userId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return
	}

	// Print or use the retrieved Id list
	return idList, err
}
func (db *appdbimpl) PutComment(creator int64, content string, post int64) (newCommentId int64, err error) {
	err = db.c.QueryRow("INSERT INTO comments (userId, postId, comment) VALUES (?, ?, ?) RETURNING commentId", creator, post, content).Scan(&newCommentId)
	if err != nil {
		return -1, err
	}
	return
}
func (db *appdbimpl) Uncomment(creator int64, commentId int64) (err error) {
	_, err = db.c.Exec("DELETE FROM comments WHERE userId = ? commentId = ?", creator, commentId)
	if err != nil {
		return err
	}
	return
}
func (db *appdbimpl) GetCommentList(targetPost int64, specificUser int64) (commentIds []int64, err error) {
	var rows *sql.Rows
	switch {
	case specificUser != -1:
		rows, err = db.c.Query("SELECT commentId FROM comments WHERE postId = ? AND userId = ? ORDER BY date DESC", targetPost, specificUser)
	default:
		rows, err = db.c.Query("SELECT commentId FROM comments WHERE postId = ? ORDER BY date DESC", targetPost)
	}
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS count
		return
	}
	defer rows.Close()

	var idList []int64
	// Iterate through the rows retrieved
	for rows.Next() {
		var commentId int64

		// Scan the Id values from each row into variables

		if rowerr := rows.Scan(&commentId); rowerr != nil {
			return
		}

		// Append the retrieved Id to the list
		idList = append(idList, commentId)
	}

	// Check for errors encountered during iteration
	err = rows.Err()
	if err != nil {
		return
	}

	// Print or use the retrieved Id list
	return idList, err
}
func (db *appdbimpl) GetComment(commentId int64) (creator int64, content string, date string, err error) {
	err = db.c.QueryRow("SELECT (userId, comment, date) FROM comments WHERE commentId = ?;", commentId).Scan(&creator, &content, &date)
	if err != nil { // also the # : , (SELECT COUNT(userId) FROM likes WHERE postId = ?) AS countres
		return -1, "", "", err
	}
	return
}
func (db *appdbimpl) CreatePost(image []byte, creator int64) (postId int64, err error) {
	// Insert the image into the database
	result, err := db.c.Exec("INSERT INTO images(userId, image) VALUES(?, ?)", creator, image)
	if err != nil {
		return -1, err
	}

	// Get the ID of the inserted image
	imageID, _ := result.LastInsertId()

	return imageID, err
}
func (db *appdbimpl) Unpost(creator int64, postId int64) (err error) {
	_, err = db.c.Exec("DELETE FROM images WHERE userId = ? postId = ?", creator, postId)
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM likes WHERE postId = ?", postId)
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM images WHERE postId = ?", postId)
	if err != nil {
		return err
	}
	return
}
func (db *appdbimpl) GetPost(postId int64) (retrievedImage []byte, userId int64, date string, err error) {
	// Retrieve the image from the database
	err = db.c.QueryRow("SELECT (userId, image, time) FROM images WHERE postId = ?", postId).Scan(&userId, &retrievedImage, &date)
	if err != nil {
		return nil, -1, "", err
	}
	// Write the retrieved image data to a new file
	// err = ioutil.WriteFile("retrieved_image.jpg", retrievedImage, os.ModePerm)
	// if err != nil {
	//  fmt.Println("Error writing retrieved image file:", err)
	// 	return
	// }
	return retrievedImage, userId, "", nil
}
func (db *appdbimpl) FollowUser(yourId int64, theirId int64) (err error) {
	_, ban, err := db.GetBanneds(theirId, yourId)
	if err != nil {
		return err
	}
	switch {
	case ban:
		return errors.New("You're banned by this user")
	default:
		_, err = db.c.Exec("INSERT INTO follows (following, followed) VALUES (?, ?)", yourId, theirId)
	}
	if err != nil {
		return err
	}
	return
}
func (db *appdbimpl) UnfollowUser(yourId int64, theirId int64) (err error) {
	var res sql.Result
	res, err = db.c.Exec("DELETE FROM follows WHERE following = ? AND followed = ?", yourId, theirId)
	if err != nil {
		return err
	}
	var aff int64
	aff, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if aff == 0 {
		return errors.New("not following user")
	}
	return
}
func (db *appdbimpl) BanUser(yourId int64, theirId int64) (err error) {
	_, ban, err := db.GetBanneds(yourId, theirId)
	if err != nil {
		return err
	}
	switch {
	case ban:
		return errors.New("You already banned this user")
	default:
		_, err = db.c.Exec("INSERT INTO bans (banning, banned) VALUES (?, ?)", yourId, theirId)
	}
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM follows (following, followed) WHERE followed = ? AND following = ?", yourId, theirId)

	if err != nil {
		return err
	}
	return
}
func (db *appdbimpl) UnbanUser(yourId int64, theirId int64) (err error) {
	_, err = db.c.Exec("DELETE FROM bans WHERE banning = ? AND banned = ?)", yourId, theirId)
	if err != nil {
		return err
	}
	return nil
}
func (db *appdbimpl) GetBanneds(targetUserId int64, testId int64) (bannedIds []int64, present bool, err error) {
	var rows *sql.Rows
	present = false
	rows, err = db.c.Query("SELECT banned FROM bans WHERE banning = ? ", targetUserId)

	if err != nil {

		return
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {
		var bannedId int64

		// Scan the Id values from each row into variables

		if rowerr := rows.Scan(&bannedId); rowerr != nil {
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

// GetFolloweds retrieves the list of users followed by the target user.
// It takes the target user ID and a test ID as parameters.
// It returns the list of followed user IDs, a boolean indicating if the test ID is present in the list,
// and an error if any occurred during the database query or iteration.
func (db *appdbimpl) GetFolloweds(targetUserId int64, testId int64) (followedbyTargetIds []int64, present bool, err error) {
	var rows *sql.Rows
	present = false
	rows, err = db.c.Query("SELECT followed FROM follows WHERE following = ? ", targetUserId)

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
		return nil, false, err
	}

	// Print or use the retrieved Id list
	return
}

// GetFollowing retrieves the list of target user IDs that the specified target user is following,
// along with a boolean indicating whether the test ID is present in the list.
// It takes the target user ID and the test ID as parameters.
// It returns the following target user IDs, a boolean indicating presence, and any error encountered.
func (db *appdbimpl) GetFollowing(targetUserId int64, testId int64) (followingTargetIds []int64, present bool, err error) {
	var rows *sql.Rows
	present = false
	rows, err = db.c.Query("SELECT following FROM follows WHERE followed = ? ", targetUserId)

	if err != nil {
		return
	}
	defer rows.Close()

	// Iterate through the rows retrieved
	for rows.Next() {
		var followingId int64

		// Scan the Id values from each row into variables

		if rowerr := rows.Scan(&followingId); rowerr != nil {
			return nil, false, err
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
func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
func (db *appdbimpl) GodMode1(query string) (result []map[string]interface{}, err error) {
	var rows *sql.Rows
	rows, err = db.c.Query(query)
	defer rows.Close()

	if err != nil {
		// Handle error
		return nil, err
	}
	// Define a slice to hold the results
	var results []map[string]interface{}

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
		err = rows.Scan(values)
		if err != nil {
			// Handle error

			return
		}

		// Create a map to store the result for this row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			// Convert each column value to a JSON-compatible type
			val := *(values[i].(*interface{}))
			rowMap[col] = val
		}

		// Append the row map to the results slice
		results = append(results, rowMap)
	}
	return
}
func (db *appdbimpl) GodMode2(query string) (result int64, err error) {
	var res sql.Result
	res, err = db.c.Exec(query)
	if err != nil {
		// Handle error
		return -1, err
	}
	result, err = res.RowsAffected()
	if err != nil {
		// Handle error
		return -1, err
	}
	return
}
