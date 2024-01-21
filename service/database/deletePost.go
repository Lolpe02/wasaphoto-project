package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) Unpost(creator int64, postId int64) (err error) {
	var res sql.Result
	res, err = db.c.Exec("DELETE FROM images WHERE userId = ? AND postId = ?", creator, postId)
	if err != nil {
		return err
	}
	var changed int64
	if changed, err = res.RowsAffected(); changed == 0 {
		if err != nil {
			return err
		}
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
	return
}

/*
func (db *appdbimpl) DeletePhoto(username string, photoID string) (errstring string, err error) {

	userID, err := db.GetUserID(username)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error getting user ID: %w", err)
	}

	_, err = db.c.Exec(`DELETE FROM posts WHERE post_ID = ? AND poster_ID = ?`, photoID, userID)

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
