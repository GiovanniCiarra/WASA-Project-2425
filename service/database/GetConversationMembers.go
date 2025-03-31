package database

import "fmt"

func (db *appdbimpl) GetConversationMembers(conversationId int) ([]User, error) {

	//Eseguo la query per recuperare i membri della conversazione
	rows, err := db.c.Query(`SELECT u.userId, u.name, u.photo
													FROM user_table u
													JOIN conversation_members_table cmt ON u.userId = cmt.userId
													WHERE cmt.conversationId = ?`, conversationId)

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

	//Costruisco l'array degli utenti da restituire
	var members []User
	for rows.Next() {
		var member User
		if err := rows.Scan(&member.UserId, &member.Name, &member.Photo); err != nil {
			return nil, fmt.Errorf("error scanning rows: %w", err)
		}

		members = append(members, member)
	}

	return members, nil
}
