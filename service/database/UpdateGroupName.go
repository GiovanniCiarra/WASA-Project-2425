package database

func (db *appdbimpl) UpdateGroupName(conversationId int, newName string) error {

	_, err := db.c.Exec(`UPDATE conversation_table
												SET name = ? 
												WHERE conversationId = ?`, newName, conversationId)
	if err != nil {
		return err
	}
	return nil
}
