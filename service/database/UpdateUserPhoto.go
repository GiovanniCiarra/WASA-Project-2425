package database

import "fmt"

func (db *appdbimpl) UpdateUserPhoto(userId int, photo string) error {

	_, err := db.c.Exec(`UPDATE user_table 
												SET photo = ? 
												WHERE userId = ?`, photo, userId)

	return fmt.Errorf("error updating photo: %w", err)
}
