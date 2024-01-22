package database

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) Unpost(creator int64, postId int64) (err error) {
	// get description and timestamp temporarily
	var desc, time string
	err = db.c.QueryRow("SELECT description, time FROM images WHERE postId = ? AND userId = ?;", postId, creator).Scan(&desc, &time)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("not found")
		}
		return err
	}
	var res sql.Result
	res, err = db.c.Exec("DELETE FROM images WHERE userId = ? AND postId = ?;", creator, postId)
	if err != nil {
		return err
	}
	var changed int64
	changed, err = res.RowsAffected()
	if err != nil {
		return err
	}
	if changed == 0 {
		return errors.New("not found")
	}
	// delete likes and comments
	_, err = db.c.Exec("DELETE FROM likes WHERE postId = ?", postId)
	if err != nil {
		return err
	}
	_, err = db.c.Exec("DELETE FROM comments WHERE postId = ?", postId)
	if err != nil {
		return err
	}
	// delete the image
	var path = os.TempDir() + "/" + strconv.FormatInt(postId, 10)
	var names []string
	names, err = filepath.Glob(path + ".*")
	if err != nil || names == nil || len(names) == 0 {
		return err
	}
	err = os.Remove(path + "." + strings.Split(names[0], ".")[1])
	if err != nil {
		// reinsert the image in the database
		_, err = db.c.Exec("INSERT INTO images (userId, postId, description, time) VALUES (?, ?);", creator, postId, desc, time)
		if err != nil {
			return err
		}
	}
	return
}

/*
func (db *appdbimpl) DeletePhoto(username string, photoID string) (errstring string, err error) {

	userID, err := db.GetUserID(username)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error getting user ID: %w", err)
	}

	_, err = db.c.Exec(`DELETE FROM posts WHERE postId = ? AND poster_ID = ?`, photoID, userID)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error deleting photo: %w", err)
	}

	// erase from /tmp/photos

	err = os.Remove("/tmp/photos/" + photoID + ".png")

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error deleting photo: %w", err)
	}

	return "", nil
}
*/
