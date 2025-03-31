package database

import "fmt"

func (db *appdbimpl) GetOtherUserInformation(userId int, conversationId int) (string, string, error) {
	var name string
	var photo string

	err := db.c.QueryRow(`SELECT u.name, u.photo 
												FROM user_table u JOIN conversation_members_table cm ON u.userId = cm.userId
												WHERE cm.conversationId = ? AND cm.userId != ?`, conversationId, userId).Scan(&name, &photo)

	if err != nil {
		return "", "", fmt.Errorf("impossible to find user information: %w", err)
	}

	return name, photo, nil
}
