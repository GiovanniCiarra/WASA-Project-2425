package database

func (db *appdbimpl) IsUserInMessageConversation(userId, messageId int) (bool, error) {
	var count int

	err := db.c.QueryRow(`SELECT COUNT(*)
												FROM conversation_members_table cm
												JOIN message_table mt ON cm.conversationId = mt.conversationId
												WHERE mt.messageId = ? AND cm.userId = ?`,
		messageId, userId).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
