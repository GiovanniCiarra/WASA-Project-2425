package database

func (db *appdbimpl) RemoveFromGroup(userId, conversationId int) error {

	_, err := db.c.Exec(`DELETE FROM conversation_members_table
												WHERE conversationId = ? AND userId = ?`, conversationId, userId)

	return err
}
