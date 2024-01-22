package database

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	// "github.com/mattn/go-sqlite3"
	// "strings"
)

func (db *appdbimpl) GetPost(postId int64) (imagepointer *os.File, imageBytes *[]byte, err error) {
	// Retrieve the imageId from the database
	_, _, _, err = db.GetMetadata(postId)
	if err != nil {
		return nil, nil, err
	}
	var path = os.TempDir() + "/" + strconv.FormatInt(postId, 10)
	var names []string
	names, err = filepath.Glob(path + ".*")
	if err != nil || names == nil || len(names) == 0 {
		return nil, nil, err
	}
	// read file and send it to the client
	imagepointer, err = os.Open(path + "." + strings.Split(names[0], ".")[1])
	if err != nil {
		return nil, nil, err
	}
	var imageBytesOb []byte
	imageBytesOb, err = os.ReadFile(path + "." + strings.Split(names[0], ".")[1])
	if err != nil {
		return nil, nil, err
	}
	imageBytes = &imageBytesOb
	return
}
