package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Register routes

	rt.router.POST("/session", rt.wrap(rt.postDoLogin))                                    //Controllato
	rt.router.GET("/users", rt.wrap(rt.getUsers))                                          //Controllato
	rt.router.PUT("/users/me/username", rt.wrap(rt.putSetMyUserName))                      //Controllato
	rt.router.PUT("/users/me/photo", rt.wrap(rt.putSetMyPhoto))                            //Controllato
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversations))                        //Controllato
	rt.router.GET("/conversations/:conversationId", rt.wrap(rt.getConversation))           //Controllato
	rt.router.POST("/conversations", rt.wrap(rt.postCreateConversation))                   //Controllato
	rt.router.POST("/conversations/:conversationId/messages", rt.wrap(rt.postSendMessage)) //Controllato
	rt.router.GET("/conversations/:conversationId/members", rt.wrap(rt.getConversationMembers))
	rt.router.POST("/messages/:messageId/forward", rt.wrap(rt.postForwardMessage))                    //Controllato
	rt.router.DELETE("/conversations/:conversationId/messages/:messageId", rt.wrap(rt.deleteMessage)) //Controllato
	rt.router.POST("/messages/:messageId/comments", rt.wrap(rt.postCommentMessage))
	rt.router.DELETE("/messages/:messageId/comments", rt.wrap(rt.deleteUncommentMessage))
	rt.router.POST("/conversations/:conversationId/add", rt.wrap(rt.postAddToGroup))
	rt.router.DELETE("/conversations/:conversationId/leave", rt.wrap(rt.deleteLeaveGroup))
	rt.router.PUT("/conversations/:conversationId/name", rt.wrap(rt.putSetGroupName))
	rt.router.PUT("/conversations/:conversationId/photo", rt.wrap(rt.putSetGroupPhoto))

	// Special routes

	return rt.router
}
