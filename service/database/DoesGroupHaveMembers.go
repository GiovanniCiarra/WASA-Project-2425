package database

func (db *appdbimpl) DoesGroupHaveMembers(conversationId int) (bool, error) {
	var count int

	err := db.c.QueryRow(`SELECT COUNT(*) FROM conversation_members_table
												WHERE conversationId = ?`, conversationId).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
