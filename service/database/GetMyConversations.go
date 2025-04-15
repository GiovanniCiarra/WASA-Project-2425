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

	//rows verrà chiuso automaticamente alla fine della funzione
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			if err == nil {
				err = errClose
			}
		}
	}()

	//Creo la Lista delle chat da restituire
	var ConversationsList []Conversation
	for rows.Next() {
		var conversation Conversation
		err = rows.Scan(&conversation.ConversationId, &conversation.Type, &conversation.Name, &conversation.Photo)
		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		//Eseguo una query per ricevere le informazioni sull'ultimo messaggio della conversazione
		var lastMessage Message
		err = db.c.QueryRow(`SELECT mt.messageId, mt.conversationId, mt.senderId, u.name, mt.content, mt.timestamp, mt.type, mt.status, mt.isForwarded 
													FROM message_table mt
													JOIN user_table u ON u.userId = mt.senderId
													WHERE mt.conversationId = ? 
													ORDER BY mt.timestamp DESC
													LIMIT 1`, conversation.ConversationId).
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
		lastMessage.Comments = nil //Non ho bisono dei commenti
		conversation.LastMessage = lastMessage

		//Se la conversazione è una chat singola, allora vado a cambiare il nome e la foto con quelle dell'altro utente
		if conversation.Type == "chat" {
			conversation.Name, conversation.Photo, err = db.GetOtherUserInformation(userId, conversation.ConversationId)
			if err != nil {
				return nil, fmt.Errorf("failed to get other user information: %w", err)
			}
		}
		ConversationsList = append(ConversationsList, conversation)
	}

	// Controllo se ci sono stati errori durante l'iterazione
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ConversationsList, nil
}
