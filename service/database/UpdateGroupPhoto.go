package database

func (db *appdbimpl) UpdateGroupPhoto(conversationId int, image []byte) error {

	_, err := db.c.Exec(`
		UPDATE conversation_table 
		SET photo = ? 
		WHERE id = ? AND type = 'group'`, image, conversationId)

	if err != nil {
		return err

	}

	return nil
}
