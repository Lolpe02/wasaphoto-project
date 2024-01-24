package database

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) CreatePost(image *multipart.File, desc *string, enc string, creator int64) (postId int64, err error) {
	// Insert the image into the database
	result, err := db.c.Exec("INSERT INTO images (userId, description) VALUES (?, ?);", creator, *desc)
	if err != nil {
		return -1, err
	}

	// Get the Id of the inserted image
	postId, err = result.LastInsertId()
	if err != nil {
		return -1, err
	}
	// create the folder if it doesn't exist
	_, err = os.Stat(os.TempDir() + "/wasaPhotos")
	if os.IsNotExist(err) {
		err = os.Mkdir(os.TempDir()+"/wasaPhotos", 0755)
		if err != nil {
			return -1, err
		}
	}
	var dest *os.File
	dest, err = os.Create(os.TempDir() + "/wasaPhotos/" + strconv.FormatInt(postId, 10) + "." + enc)
	if err != nil {
		err = db.Unpost(creator, postId)
		if err != nil {
			return -1, err
		}
		return -1, err
	}
	// insert the image in images folder
	var bytesWritten int64
	bytesWritten, err = io.Copy(dest, *image)
	if err != nil || bytesWritten == 0 {
		err = db.Unpost(creator, postId)
		if err != nil {
			return -1, err
		}
		return -1, errors.New("nothing written")
	}

	return
}
