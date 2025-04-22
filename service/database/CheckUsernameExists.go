package database

import (
	"database/sql"
	"errors"
)

// Funzione che controlla che uno username non sia già associato ad un utente
func (db *appdbimpl) CheckUsernameExists(username string) (bool, error) {

	// Eseguo la Query
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM user_table WHERE name = ?", username).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	// Se count > 0 allora lo username è già associato ad un utente
	if count > 0 {
		return true, nil
	}

	// Non è stato trovato un match
	return false, nil
}
