package database

import "fmt"

func (db *appdbimpl) GetMessageById(messageId int) (Message, error) {
	var message Message

	err := db.c.QueryRow(`SELECT content, type FROM message_table
												WHERE messageId = ?`, messageId).Scan(
		&message.Content,
		&message.Type,
	)
	if err != nil {
		return Message{}, fmt.Errorf("error in query execution: %w", err)
	}

	return message, nil
}
