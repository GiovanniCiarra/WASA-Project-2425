package database

func (db *appdbimpl) UpdateGroupPhoto(conversationId int, image string) error {

	_, err := db.c.Exec(`
		UPDATE conversation_table 
		SET photo = ? 
		WHERE conversationId = ? AND type = 'group'`, image, conversationId)

	if err != nil {
		return err

	}

	return nil
}
