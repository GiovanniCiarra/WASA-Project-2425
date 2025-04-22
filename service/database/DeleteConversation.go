package database

import "fmt"

func (db *appdbimpl) DeleteConversation(conversationId int) error {
	// Cancello i membri della conversazione
	_, err := db.c.Exec(`DELETE FROM conversation_members_table
												WHERE conversationId = ?`, conversationId)

	if err != nil {
		return fmt.Errorf("error deleting users from conversation: %w", err)
	}

	// Cancello la conversazione
	_, err = db.c.Exec(`DELETE FROM conversation_table WHERE conversationId = ?`, conversationId)
	if err != nil {
		return fmt.Errorf("error in query execution: %w", err)
	}

	return nil
}
