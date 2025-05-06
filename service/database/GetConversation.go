package database

import "fmt"

func (db *appdbimpl) GetConversation(userId int, conversationId int) ([]Message, error) {

	// Recupero i messaggi della conversazione
	rows, err := db.c.Query(`SELECT mt.messageId, mt.conversationId, mt.content, mt.senderId, mt.timestamp, mt.type, mt.status, mt.isForwarded, u.name, mt.replyTo
													 FROM message_table mt 
													 JOIN user_table u ON u.userId = mt.senderId 
													 WHERE mt.conversationId = ?
													 ORDER BY mt.timestamp ASC`, conversationId)

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

	// Creo la Lista dei messaggi da restituire
	var messageList []Message

	for rows.Next() {
		var message Message
		err = rows.Scan(&message.MessageId, &message.ConversationId, &message.Content, &message.SenderId,
			&message.Timestamp, &message.Type, &message.Status, &message.IsForwarded, &message.SenderUsername, &message.IsReply)

		if err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		messageList = append(messageList, message)

		// Controllo se ci sono stati errori durante l'iterazione
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("error retrieving messages: %w", err)
		}
	}

	// Creo per ogni messaggio il suo array di commenti
	for i := range messageList {
		message := &messageList[i]

		comments, err2 := db.c.Query(`SELECT ct.commentId, ct.messageId,ct.senderId, ct.content, ct.timestamp, ut.name
																	FROM comment_table ct 
																	JOIN user_table ut ON ut.userId = ct.senderId
																	WHERE ct.messageId = ?`, message.MessageId)
		if err2 != nil {
			return nil, fmt.Errorf("error in query execution: %w", err2)
		}

		defer comments.Close()

		var commentsList []Comment
		for comments.Next() {
			var comment Comment
			err2 = comments.Scan(&comment.CommentId, &comment.MessageId, &comment.SenderId, &comment.Content, &comment.Timestamp, &comment.Sender)
			if err2 != nil {
				return nil, fmt.Errorf("error scanning comment rows: %w", err2)
			}
			commentsList = append(commentsList, comment)
		}

		if err2 := comments.Close(); err2 != nil {
			return nil, fmt.Errorf("error closing comments rows: %w", err2)
		}

		// Controllo se ci sono stati errori durante l'iterazione
		if err := comments.Err(); err != nil {
			return nil, err
		}

		message.Comments = commentsList

		// Se è una risposta, prendo i dati del messaggio originale
		if message.IsReply != -1 {
			var info ReplyInfo

			err = db.c.QueryRow(`SELECT mt.content, mt.type, u.name
			                  FROM message_table mt
												JOIN user_table u ON mt.senderId = u.userId
												WHERE mt.messageId = ?`, message.IsReply).Scan(&info.Content, &info.Type, &info.SenderUsername)

			if err != nil {
				return nil, fmt.Errorf("error in query execution: %w", err2)
			}
			message.ReplyTo = info
		}

	}

	return messageList, nil

}
