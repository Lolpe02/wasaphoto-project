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

/*
func (db *appdbimpl) UploadPhoto(username string, photo components.Photo, photo_ID string) (errstring string, err error) {

	var data string = photo.Data

	encoded_data, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error encoding data: %w", err)
	}

	// Get user ID

	userID, err := db.GetUserID(username)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error getting user ID: %w", err)
	}

	// Get current time

	creation_time := time.Now().Format(time.RFC3339)

	// Insert photo

	_, err = db.c.Exec(`INSERT OR REPLACE INTO posts (postId, poster_ID, description, creation_date) VALUES (?, ?, ?, ?)`, photo_ID, userID, photo.Desc, creation_time)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error inserting photo: %w", err)
	}

	PNG_reader := bytes.NewReader(encoded_data)

	img, err := png.Decode(PNG_reader)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error decoding PNG: %w", err)
	}

	_, err = os.Stat("/tmp/photos")

	if os.IsNotExist(err) {
		err = os.Mkdir("/tmp/photos", 0755)
		if err != nil {
			return components.InternalServerError, fmt.Errorf("error creating /tmp/photos: %w", err)
		}
	}

	f, err := os.OpenFile("/tmp/photos/"+photo_ID+".png", os.O_WRONLY|os.O_CREATE, 0777)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error creating file: %w", err)
	}

	defer f.Close()

	err = png.Encode(f, img)

	if err != nil {
		return components.InternalServerError, fmt.Errorf("error encoding PNG: %w", err)
	}

	return "", nil
}
*/
