package database

import (
// "github.com/mattn/go-sqlite3"
// "strings"
)

func (db *appdbimpl) PutComment(creator int64, content string, post int64) (newCommentId int64, err error) {
	err = db.c.QueryRow("INSERT INTO comments (userId, postId, comment) VALUES (?, ?, ?) RETURNING commentId", creator, post, content).Scan(&newCommentId)
	if err != nil {
		return -1, err
	}
	return
}
