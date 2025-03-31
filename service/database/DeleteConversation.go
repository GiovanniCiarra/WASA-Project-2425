package database

func (db *appdbimpl) DeleteConversation(conversationId int) error {
	//Cancello i membri della conversazione
	_, err := db.c.Exec(`DELETE FROM conversation_members_table
												WHERE conversationId = ?`, conversationId)

	if err != nil {
		return err
	}

	//Cancello la conversazione
	_, err = db.c.Exec(`DELETE FROM conversation_table WHERE conversationId = ?`, conversationId)
	return err
}
