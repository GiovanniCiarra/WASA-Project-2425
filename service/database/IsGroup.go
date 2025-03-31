package database

func (db *appdbimpl) IsGroup(conversationId int) (bool, error) {

	var conversationType string
	err := db.c.QueryRow(`SELECT type FROM conversation_table
												WHERE conversationId = ?`, conversationId).Scan(&conversationType)

	if err != nil {
		return false, err
	}

	return conversationType == "group", nil
}
