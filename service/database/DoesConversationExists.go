package database

import "fmt"

func (db *appdbimpl) DoesConversationExists(userId, targetUserId int) (int, error) {
	var conversationId int

	err := db.c.QueryRow(`SELECT ct.conversationId 
												FROM conversation_table ct
												JOIN conversation_members_table cmt1 ON ct.conversationId = cmt1.conversationId
												JOIN conversation_members_table cmt2 ON ct.conversationId = cmt2.conversationId
												WHERE ct.type = 'chat' 
												AND cmt1.userId = ?
												AND cmt2.userId = ?
												GROUP BY ct.conversationId 
												HAVING COUNT(DISTINCT cmt1.userId || cmt2.userId) = 2`, userId, targetUserId).Scan(&conversationId)
	if err != nil {
		return -1, fmt.Errorf("error in query execution: %w", err)
	}

	return conversationId, nil
}
