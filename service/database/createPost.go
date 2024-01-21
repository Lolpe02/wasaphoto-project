package database

import (
	"image"
	"os"
	"strconv"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) CreatePost(image image.Image, creator int64) (postId int64, err error) {
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
	err = os.WriteFile("images/"+strconv.FormatInt(postId, 10)+".jpg", []byte{}, os.ModePerm)
	if err != nil {
		return
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

	_, err = db.c.Exec(`INSERT OR REPLACE INTO posts (post_ID, poster_ID, description, creation_date) VALUES (?, ?, ?, ?)`, photo_ID, userID, photo.Desc, creation_time)

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
