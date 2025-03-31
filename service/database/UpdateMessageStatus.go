package database

import "fmt"

func (db *appdbimpl) UpdateMessageStatus(messageId int) error {

	//Controllo se tutti i destinatari hanno lo stato "received"
	var allReceived bool
	err := db.c.QueryRow(`SELECT COUNT(*) = 0 FROM message_status_table
												WHERE messageId = ? AND status = 'sent'`, messageId).Scan(&allReceived)

	if err != nil {
		return fmt.Errorf("error checking message 'sent' status: %w", err)
	}

	//Se hanno tutti 'received' (non ci sono 'sent') aggiorno lo stato in message_table
	if allReceived {
		_, err = db.c.Exec(`UPDATE message_table
												SET status = 'received'
												WHERE messageId = ? AND status = 'sent'`, messageId)
		if err != nil {
			return fmt.Errorf("error updating message general status to 'received': %w", err)
		}
	}

	//Controllo se tutti i destinatari hanno lo stato "read"
	var allRead bool
	err = db.c.QueryRow(`SELECT COUNT(*) = 0 FROM message_status_table
												WHERE messageId = ? AND status != 'read'`, messageId).Scan(&allRead)

	if err != nil {
		return fmt.Errorf("error checking message 'received' status: %w", err)
	}

	//Se hanno tutti 'read' aggiorno lo stato in message_table
	if allRead {
		_, err = db.c.Exec(`UPDATE message_table
												SET status = 'read'
												WHERE messageId = ? AND status = 'received'`, messageId)
		if err != nil {
			return fmt.Errorf("error updating message general status to 'read': %w", err)
		}
	}

	return nil
}
