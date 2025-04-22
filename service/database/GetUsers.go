package database

// Funzione che ritorna la lista di tutti gli utenti
func (db *appdbimpl) GetUsers(userId int) ([]User, error) {

	rows, err := db.c.Query(`SELECT * FROM user_table`)

	if err != nil {
		return nil, err
	}

	// rows verrà chiuso automaticamente alla fine della funzione
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			if err == nil {
				err = errClose
			}
		}
	}()

	// Creo la UsersList da restituire, controllando di non restituire l'utente chiamante
	var UsersList []User
	for rows.Next() {
		var User User
		err = rows.Scan(&User.UserId, &User.Name, &User.Photo)

		if err != nil {
			return nil, err
		}

		if User.UserId != userId {
			UsersList = append(UsersList, User)
		}
	}

	// Controllo se ci sono stati errori durante l'iterazione
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return UsersList, err
}
