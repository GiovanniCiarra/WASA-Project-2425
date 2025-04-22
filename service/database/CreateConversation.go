package database

import "fmt"

func (db *appdbimpl) CreateConversation(userId int, name string, photo string, members []int) (int, error) {

	var conversationId int
	conversationType := "group"

	// Controllo se si tratta di una chat privata o un gruppo
	if len(members) == 1 {
		conversationType = "chat"
		name = ""
		photo = ""
	}

	// Inserisco la nuova conversazione nel database
	err := db.c.QueryRow(`INSERT INTO conversation_table (type, name, photo) VALUES (?, ?, ?) RETURNING conversationId`,
		conversationType, name, photo).Scan(&conversationId)
	if err != nil {
		return -1, fmt.Errorf("error while creating the new conversation: %w", err)
	}

	// Aggiungo i membri alla conversazione
	members = append(members, userId) // Bisogna aggiungere anche l'utente che la sta creando
	for _, member := range members {
		_, err = db.c.Exec(`INSERT INTO conversation_members_table (conversationId, userId) VALUES (?,?)`, conversationId, member)
		if err != nil {
			return -1, fmt.Errorf("error adding user to conversation: %w", err)
		}
	}

	return conversationId, nil
}
