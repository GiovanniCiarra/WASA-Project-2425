package database

import "fmt"

func (db *appdbimpl) SaveMessage(conversationId, senderId, isForwarded int, content, msgType string) (int, error) {
	var messageId int

	// Eseguo la query per inserire il messaggio nel database
	err := db.c.QueryRow(`INSERT INTO message_table 
												(conversationId, senderId, content, type, timestamp, status, isForwarded)
												VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP, 'sent', ?) RETURNING messageId`,
		conversationId, senderId, content, msgType, isForwarded).Scan(&messageId)

	if err != nil {
		return -1, fmt.Errorf("error in query execution: %w", err)
	}

	// Recupero tutti i membri della conversazione
	members, err := db.GetConversationMembers(conversationId)

	if err != nil {
		return -1, fmt.Errorf("error retreaving conversation members: %w", err)
	}

	// Inserisco in message_status_table lo stato 'sent' per ogni destinatartio del messaggio
	for _, member := range members {
		if member.UserId == senderId {
			continue // Non devo inserire il sender del messaggio
		}

		_, err := db.c.Exec(`INSERT INTO message_status_table (messageId, userId, status)
													VALUES (?, ?, 'sent')`, messageId, member.UserId)

		if err != nil {
			return -1, fmt.Errorf("error in query execution: %w", err)
		}
	}

	return messageId, nil
}
