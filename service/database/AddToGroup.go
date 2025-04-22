package database

func (db *appdbimpl) AddToGroup(userId, conversationId int) error {

	_, err := db.c.Exec(`INSERT INTO conversation_members_table (conversationId, userId) 
												VALUES (?, ?)`, conversationId, userId)

	if err != nil {
		return err
	}
	return nil
}
