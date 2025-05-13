package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Register routes

	// Login e Profilo utente
	rt.router.POST("/session", rt.wrap(rt.postDoLogin))
	rt.router.GET("/users", rt.wrap(rt.getUsers))
	rt.router.GET("/profile", rt.wrap(rt.getUserProfile))
	rt.router.PUT("/profile/username", rt.wrap(rt.putSetMyUserName))
	rt.router.PUT("/profile/photo", rt.wrap(rt.putSetMyPhoto))

	// Conversazioni
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))
	rt.router.POST("/conversations/create", rt.wrap(rt.postCreateConversation))
	rt.router.GET("/privateChat/:userId", rt.wrap(rt.getConversationByUserId))

	// Operazioni su conversazioni
	rt.router.GET("/conversation-info/:conversationId", rt.wrap(rt.getConversation))
	rt.router.GET("/conversation-members/:conversationId", rt.wrap(rt.getConversationMembers))
	rt.router.DELETE("/conversation-delete/:conversationId", rt.wrap(rt.deleteConversation))

	// Messaggi
	rt.router.POST("/conversation-messages/:conversationId", rt.wrap(rt.postSendMessage))
	rt.router.DELETE("/conversation-messages/:conversationId/:messageId", rt.wrap(rt.deleteMessage))
	rt.router.POST("/messages/:messageId/forwards", rt.wrap(rt.postForwardMessage))
	rt.router.POST("/messages/:messageId/replies", rt.wrap(rt.postReplyMessage))

	// Commenti
	rt.router.POST("/messages/:messageId/comments", rt.wrap(rt.postCommentMessage))
	rt.router.DELETE("/messages/:messageId/comments", rt.wrap(rt.deleteUncommentMessage))

	// Gruppi
	rt.router.POST("/conversation-add/:conversationId", rt.wrap(rt.postAddToGroup))
	rt.router.DELETE("/conversation-members/:conversationId/me", rt.wrap(rt.deleteLeaveGroup))
	rt.router.PUT("/conversation-name/:conversationId", rt.wrap(rt.putSetGroupName))
	rt.router.PUT("/conversation-photo/:conversationId", rt.wrap(rt.putSetGroupPhoto))

	// Special routes
	return rt.router
}
