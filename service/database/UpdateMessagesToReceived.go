package database

import "fmt"

func (db *appdbimpl) UpdateMessagesToReceived(userId int) error {

	//Seleziono tutti i messaggi che devono essere aggiornati
	rows, err := db.c.Query(`SELECT messageId FROM message_status_table
														WHERE userId = ? AND status = 'sent'`, userId)

	if err != nil {
		return fmt.Errorf("error in query execution: %w", err)
	}

	//rows verrà chiuso automaticamente alla fine della funzione
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			if err == nil {
				err = errClose
			}
		}
	}()

	//Creo un array con tutti i messageId che sto aggiornando,a
	//per poi controllare uno a uno se hanno tutti lo stesso stato per poter aggiornare lo stato generale del messaggio
	var messageIds []int
	for rows.Next() {
		var messageId int
		if err := rows.Scan(&messageId); err != nil {
			return fmt.Errorf("error scanning rows: %w", err)
		}
		messageIds = append(messageIds, messageId)
	}

	//Aggiorno i messaggi a 'received' nella message_status_table per ogni messaggio
	_, err = db.c.Exec(`UPDATE message_status_table 
												SET status = 'received'
												WHERE userId = ? AND status = 'sent'`, userId)
	if err != nil {
		return fmt.Errorf("error in query execution: %w", err)
	}

	//Controllo se devo aggiornare lo stato generale per ogni messaggio
	for _, msgId := range messageIds {
		err := db.UpdateMessageStatus(msgId)

		if err != nil {
			return fmt.Errorf("error updating message status: %w", err)
		}
	}

	return nil

}
