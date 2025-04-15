package database

import "fmt"

func (db *appdbimpl) DoesConversationExists(userId, targetUserId int) (int, error) {
	var conversationId int

	err := db.c.QueryRow(`SELECT ct.conversationId
												FROM conversation_table ct
												JOIN conversation_members_table cmt ON ct.conversationId = cmt.conversationId
												WHERE ct.type = 'chat' 
												GROUP BY ct.conversationId 
												HAVING COUNT(*) = 2 
												AND SUM(cmt.userId = ?) > 0
												AND SUM(cmt.userId = ?) > 0 `, userId, targetUserId).Scan(&conversationId)
	if err != nil {
		return -1, fmt.Errorf("error in query execution: %w", err)
	}

	return conversationId, nil
}
