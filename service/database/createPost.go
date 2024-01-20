package database

import (
	"fmt"
	"os"
	"strconv"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) CreatePost(image []byte, creator int64) (postId int64, err error) {
	// Insert the image into the database
	result, err := db.c.Exec("INSERT INTO images (userId) VALUES (?)", creator)
	if err != nil {
		return -1, err
	}

	// Get the Id of the inserted image
	postId, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}
	// get format of image

	// insert the image in images folder
	err = os.WriteFile("images/"+strconv.FormatInt(postId, 10)+".jpg", image, os.ModePerm)
	if err != nil {
		fmt.Println("Error writing image file:", err)
		return
	}

	return
}
