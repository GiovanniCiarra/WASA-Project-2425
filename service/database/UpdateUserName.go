package database

import "fmt"

func (db *appdbimpl) UpdateUserName(userId int, username string) error {

	_, err := db.c.Exec(`UPDATE user_table 
												SET name = ? 
												WHERE userId = ?`, username, userId)

	return fmt.Errorf("error updating username: %w", err)
}
