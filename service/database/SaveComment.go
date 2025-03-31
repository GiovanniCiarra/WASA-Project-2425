package database

func (db *appdbimpl) SaveComment(messageId, senderId int, content string) (int, error) {
	var commentId int

	err := db.c.QueryRow(`INSERT INTO comment_table (messageId, senderId, content, timestamp)
												VALUES (?, ?, ?, CURRENT_TIMESTAMP) RETURNING commentId`, messageId, senderId, content).Scan(&commentId)

	return commentId, err
}
