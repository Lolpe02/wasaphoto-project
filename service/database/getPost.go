package database

import (
	"database/sql"
	"errors"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetPost(postId int64) (retrievedImage []byte, userId int64, date string, err error) {
	// Retrieve the image from the database
	err = db.c.QueryRow("SELECT userId, image, time FROM images WHERE postId = ?", postId).Scan(&userId, &retrievedImage, &date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return nil values and specify the error
			return nil, -1, "", errors.New("not found")
		}
		return nil, -1, "", err
	}
	// Write the retrieved image data to a new file
	// err = ioutil.WriteFile("retrieved_image.jpg", retrievedImage, os.ModePerm)
	// if err != nil {
	// 	return
	// }
	return
}
