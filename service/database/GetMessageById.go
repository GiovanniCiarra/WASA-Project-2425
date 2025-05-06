package database

import "fmt"

func (db *appdbimpl) GetMessageById(messageId int) (Message, error) {
	var message Message
	message.MessageId = messageId

	err := db.c.QueryRow(`SELECT conversationId, senderId, content, timestamp, type FROM message_table
												WHERE messageId = ?`, messageId).Scan(
		&message.ConversationId,
		&message.SenderId,
		&message.Content,
		&message.Timestamp,
		&message.Type,
	)
	if err != nil {
		return Message{}, fmt.Errorf("error in query execution: %w", err)
	}

	return message, nil
}
