package database

import "fmt"

func (db *appdbimpl) IsMessageOwner(userId, messageId, conversationId int) (bool, error) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM message_table
												WHERE messageId = ? AND senderId = ? AND conversationId = ?`,
		messageId, userId, conversationId).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error in query execution: %w", err)
	}

	return count > 0, nil
}
