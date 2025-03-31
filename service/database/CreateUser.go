package database

import "fmt"

func (db *appdbimpl) CreateUser(name string) (int, error) {

	result, err := db.c.Exec("INSERT INTO user_table (name, photo) VALUES (?,?)", name, "")
	if err != nil {
		return -1, fmt.Errorf("error creating user: %w", err)
	}

	var userId int64
	userId, err = result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("error retreaving last userId: %w", err)
	}
	return int(userId), nil
}
