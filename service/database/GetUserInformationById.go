package database

import "fmt"

func (db *appdbimpl) GetUserInformationById(userId int) (string, string, error) {

	var name, photo string
	err := db.c.QueryRow(`SELECT name, photo FROM user_table WHERE userId = ?`, userId).Scan(&name, &photo)

	if err != nil {
		return "", "", fmt.Errorf("error in query execution: %w", err)
	}

	return name, photo, nil
}
