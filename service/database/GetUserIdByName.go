package database

import "fmt"

func (db *appdbimpl) GetUserIdByName(name string) (int, error) {
	var userId int
	err := db.c.QueryRow("SELECT userId FROM user_table WHERE name = ?", name).Scan(&userId)
	if err != nil {
		return -1, fmt.Errorf("error retreaving userId: %w", err)
	}
	return userId, nil
}
