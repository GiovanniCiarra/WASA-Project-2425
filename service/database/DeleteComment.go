package database

func (db *appdbimpl) DeleteComment(messageId, senderId int) error {

	_, err := db.c.Exec(`DELETE FROM comment_table
												WHERE messageId = ? AND senderId = ?`, messageId, senderId)

	if err != nil {
		return err
	}
	return nil
}
