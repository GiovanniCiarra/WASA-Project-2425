package database

func (db *appdbimpl) DeleteMessage(messageId int) error {

	_, err := db.c.Exec(`DELETE FROM message_table WHERE messageId = ?`, messageId)
	return err
}
