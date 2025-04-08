package database

import "fmt"

func (db *appdbimpl) IsUserInConversation(userId int, conversationId int) (bool, error) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM conversation_members_table
												WHERE conversationId = ? AND userId = ?`,
		conversationId, userId).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error checking conversation membership: %w", err)
	}

	return count > 0, nil
}
