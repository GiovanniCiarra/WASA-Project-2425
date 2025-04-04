package database

import "fmt"

func (db *appdbimpl) DoesUserExist(userId int) (bool, error) {
	var count int
	err := db.c.QueryRow(`SELECT COUNT(*) FROM user_table 
																	WHERE userId = ?`, userId).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("error in query execution: %w", err)
	}

	return count > 0, nil
}
