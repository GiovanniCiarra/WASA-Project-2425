package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetMyConversations(userId int) ([]Conversation, error) {

	rows, err := db.c.Query(`SELECT ct.conversationId, ct.type, ct.name, ct.photo
													 FROM conversation_members_table cm
													 JOIN conversation_table ct ON cm.conversationId = ct.conversationId 
													 WHERE cm.userId = ?`, userId)
	if err != nil {
		return nil, fmt.Errorf("error in query execution: %w", err)
	}

	// rows verrà chiuso automaticamente alla fine della funzione
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			if err == nil {
				err = errClose
			}
		}
	}()

	// Creo la Lista delle chat da restituire
	var conversationsList []Conversation
	for rows.Next() {
		var conversation Conversation
		err = rows.Scan(&conversation.ConversationId, &conversation.Type, &conversation.Name, &conversation.Photo)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		conversationsList = append(conversationsList, conversation)
	}
	// Controllo se ci sono stati errori durante l'iterazione
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Eseguo una query per ricevere le informazioni sull'ultimo messaggio di ogni conversazione
	for i := range conversationsList {
		conv := &conversationsList[i]

		var lastMessage Message
		err = db.c.QueryRow(`SELECT mt.messageId, mt.conversationId, mt.senderId, u.name, mt.content, mt.timestamp, mt.type, mt.status, mt.isForwarded 
														FROM message_table mt
														JOIN user_table u ON u.userId = mt.senderId
														WHERE mt.conversationId = ? 
														ORDER BY mt.timestamp DESC
														LIMIT 1`, conv.ConversationId).
			Scan(&lastMessage.MessageId, &lastMessage.ConversationId, &lastMessage.SenderId,
				&lastMessage.SenderUsername, &lastMessage.Content, &lastMessage.Timestamp,
				&lastMessage.Type, &lastMessage.Status, &lastMessage.IsForwarded)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				lastMessage = Message{}
			} else {
				return nil, fmt.Errorf("error retreaving last message: %w", err)
			}
		}
		lastMessage.Comments = nil // Non ho bisono dei commenti
		conv.LastMessage = lastMessage

		// Se la conversazione è una chat singola, allora vado a cambiare il nome e la foto con quelle dell'altro utente
		if conv.Type == "chat" {
			conv.Name, conv.Photo, err = db.GetOtherUserInformation(userId, conv.ConversationId)
			if err != nil {
				return nil, fmt.Errorf("failed to get other user information: %w", err)
			}
		}

	}
	return conversationsList, nil
}
