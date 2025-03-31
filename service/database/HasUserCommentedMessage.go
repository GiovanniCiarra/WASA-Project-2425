package database

func (db *appdbimpl) HasUserCommentedMessage(senderId, messageId int) (bool, error) {
	var count int

	err := db.c.QueryRow(`SELECT COUNT(*)
												FROM comment_table
												WHERE messageId = ? AND senderId = ?`,
		messageId, senderId).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
