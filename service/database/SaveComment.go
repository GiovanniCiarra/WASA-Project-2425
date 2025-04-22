package database

import "fmt"

func (db *appdbimpl) SaveComment(messageId, senderId int, content string) (int, error) {
	var commentId int

	err := db.c.QueryRow(`INSERT INTO comment_table (messageId, senderId, content, timestamp)
												VALUES (?, ?, ?, CURRENT_TIMESTAMP) RETURNING commentId`, messageId, senderId, content).Scan(&commentId)

	if err != nil {
		return -1, fmt.Errorf("error in wuery execution: %w", err)
	}

	return commentId, nil
}
