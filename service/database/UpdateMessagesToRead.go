package database

import "fmt"

func (db *appdbimpl) UpdateMessagesToRead(userId, conversationId int) error {

	//Seleziono i messaggi che devono essere aggiornati
	rows, err := db.c.Query(`SELECT messageId 
														FROM message_status_table
														WHERE userId = ? AND status = 'received'
														AND messageId IN (
															SELECT messageId FROM message_table WHERE conversationId = ?)`, userId, conversationId)

	if err != nil {
		return err
	}

	//rows verrà chiuso automaticamente alla fine della funzione
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			if err == nil {
				err = errClose
			}
		}
	}()

	//Creo un array con tutti i messageId che sto aggiornando,
	//per poi controllare uno a uno se hanno tutti lo stesso stato per poter aggiornare lo stato generale
	var messageIds []int
	for rows.Next() {
		var messageId int
		if err := rows.Scan(&messageId); err != nil {
			return fmt.Errorf("error scanning rows: %w", err)
		}
		messageIds = append(messageIds, messageId)
	}

	//Aggiorno i messaggi a 'read'nella message_status_table
	_, err = db.c.Exec(`UPDATE message_status_table
											SET status = 'read'
											WHERE userId = ? AND status = 'received'
											AND messageId IN (
												SELECT messageId FROM message_table WHERE conversationId = ?)`, userId, conversationId)

	if err != nil {
		return fmt.Errorf("error updating message status to 'read': %w", err)
	}

	//Controllo se devo aggiornare lo stato generale
	for _, messageId := range messageIds {
		err := db.UpdateMessageStatus(messageId)
		if err != nil {
			return fmt.Errorf("error updating general message status: %w", err)
		}
	}

	return nil
}
