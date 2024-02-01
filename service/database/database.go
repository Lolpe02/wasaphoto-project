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
	"mime/multipart"
	"os"

	"github.com/gofrs/uuid"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateUser(username string) (newUserId int64, alreadyExists bool, err error)
	ChangeUsername(yourUserId int64, newUsername string) (err error)
	SearchByUsername(targetUser string) (selUserId int64, err error)
	SearchById(targetUserId int64) (selUserName string, subscription string, err error)
	GetInfo(targetUserName string) (userIds []int64, usernamesAlike []string, err error)
	GetProfile(targetUserId int64) (postIds []int64, err error)
	GetFeed(yourId int64) (postIds []int64, err error)
	PutLike(targetPost int64, creator int64) (err error)
	Unlike(targetPost int64, creator int64) (err error)
	GetLikes(targetPost int64) (userIds []int64, err error)
	PutComment(creator int64, content string, post int64) (newCommentId int64, err error)
	Uncomment(creator int64, postId int64, commentId int64) (err error)
	GetCommentList(targetPost int64, specificUser int64) (commentIds []int64, err error)
	GetComment(commentId int64) (creator int64, postId int64, content string, date string, err error)
	CreatePost(image *multipart.File, desc *string, enc string, creator int64) (postId int64, err error)
	Unpost(creator int64, postId int64) (err error)
	GetPost(postId int64) (imagepointer *os.File, imageBytes *[]byte, format string, err error)
	GetMetadata(postId int64) (userId int64, description string, date string, err error)
	FollowUser(yourId int64, theirId int64) (alreadyExists bool, err error)
	UnfollowUser(yourId int64, theirId int64) (err error)
	BanUser(yourId int64, theirId int64) (err error)
	UnbanUser(yourId int64, theirId int64) (err error)
	GetBanneds(targetUserId int64, testId int64) (bannedIds []int64, present bool, err error)
	GetFolloweds(targetUserId int64, testId int64) (followedbyTargetIds []int64, followedbyTargetNames []string, present bool, err error)
	GetFollowing(targetUserId int64, testId int64) (followingTargetIds []int64, followingTargetNames []string, present bool, err error)
	Ping() error
	GodMode1(query string) (result []map[string]interface{}, err error)
	GodMode2(query string) (result int64, err error)
}

type appdbimpl struct {
	c       *sql.DB
	uuidGen *uuid.Gen
}

const NotFound = string("not found")

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB, genId *uuid.Gen) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	var PRAGMAactive bool
	err7 := db.QueryRow("PRAGMA foreign_keys;").Scan(&PRAGMAactive)
	if errors.Is(err7, sql.ErrNoRows) {
		// do nothing
		// return nil, fmt.Errorf("error activating foreign keys: %w", err7)
	} else if err7 != nil {
		return nil, fmt.Errorf("error not activating foreign keys: %w", err7)
	}
	if !PRAGMAactive {
		_, err7 = db.Exec("PRAGMA foreign_keys = ON;")
		if err7 != nil {
			return nil, fmt.Errorf("error activating foreign keys: %w", err7)
		}
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure

	var users string
	err1 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&users)
	if errors.Is(err1, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE users (
			 userId INTEGER PRIMARY KEY,
			 userName TEXT UNIQUE NOT NULL, 
			 date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);`
		_, err1 = db.Exec(sqlStmt)
		if err1 != nil {
			return nil, fmt.Errorf("error creating database structure table users: %w", err1)
		}
	}
	var images string
	err2 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='images';`).Scan(&images)
	if errors.Is(err2, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE images (postId INTEGER PRIMARY KEY, userId INTEGER NOT NULL, description TEXT, time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE ON UPDATE CASCADE);`
		_, err2 = db.Exec(sqlStmt)
		if err2 != nil {
			return nil, fmt.Errorf("error creating database structure table images: %w", err2)
		}
	}
	var likes string
	err3 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='likes';`).Scan(&likes)
	if errors.Is(err3, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE likes (
			 userId INTEGER NOT NULL,
			 postId INTEGER NOT NULL,
			 PRIMARY KEY (userId, postId),
			 FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE ON UPDATE CASCADE,
			 FOREIGN KEY (postId) REFERENCES images(postId) ON DELETE CASCADE ON UPDATE CASCADE);`
		_, err3 = db.Exec(sqlStmt)
		if err3 != nil {
			return nil, fmt.Errorf("error creating database structure table likes: %w", err3)
		}
	}

	var comments string
	err4 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='comments';`).Scan(&comments)
	if errors.Is(err4, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE comments (
			 commentId INTEGER PRIMARY KEY,
			 userId INTEGER NOT NULL,
			 postId INTEGER NOT NULL,
			 comment TEXT NOT NULL,
			 date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			 FOREIGN KEY (userId) REFERENCES users(userId) ON DELETE CASCADE ON UPDATE CASCADE,
			 FOREIGN KEY (postId) REFERENCES images(postId) ON DELETE CASCADE ON UPDATE CASCADE);`
		_, err4 = db.Exec(sqlStmt)
		if err4 != nil {
			return nil, fmt.Errorf("error creating database structure table comments: %w", err4)
		}
	}
	var bans string
	err5 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='bans';`).Scan(&bans)
	if errors.Is(err5, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE bans (
			 banning INTEGER NOT NULL,
			 banned INTEGER NOT NULL,
			 PRIMARY KEY (banning, banned), 
			 FOREIGN KEY (banning) REFERENCES users(userId) ON DELETE CASCADE ON UPDATE CASCADE, 
			 FOREIGN KEY (banned) REFERENCES users(userId) ON DELETE CASCADE ON UPDATE CASCADE, 
			 CHECK (banning != banned));`
		_, err5 = db.Exec(sqlStmt)
		if err5 != nil {
			return nil, fmt.Errorf("error creating database structure table bans: %w", err5)
		}
	}
	var follows string
	err6 := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='follows';`).Scan(&follows)
	if errors.Is(err6, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE follows (
			 following INTEGER NOT NULL,
			 followed INTEGER NOT NULL,
			 PRIMARY KEY (following, followed),
			 FOREIGN KEY (following) REFERENCES users(userId) ON DELETE CASCADE ON UPDATE CASCADE,
			 FOREIGN KEY (followed) REFERENCES users(userId) ON DELETE CASCADE ON UPDATE CASCADE,
			 CHECK (following != followed));`
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

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
