package api

import (
	"Progetto/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteMessage(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

	//Estraggo lo userId dall'Header
	userId, err := extractUserID(reader, ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve userId")
		http.Error(writer, "Failed to retrieve userId", http.StatusUnauthorized)
		return
	}

	//Controllo che lo userId sia associato ad un utente valido
	valid, err := rt.db.DoesUserExist(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Impossible to validate user")
		http.Error(writer, "Impossible to validate user", http.StatusInternalServerError)
		return
	}

	if !valid {
		ctx.Logger.Error("userId is not valid")
		http.Error(writer, "UserId is not valid", http.StatusUnauthorized)
		return
	}

	//Estraggo il conversationId e il messageId dai parametri della richiesta
	conversationId, err := strconv.Atoi(params.ByName("conversationId"))
	if err != nil {
		http.Error(writer, "Invalid conversationId", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid conversationId")
		return
	}

	messageId, err := strconv.Atoi(params.ByName("messageId"))
	if err != nil {
		http.Error(writer, "Invalid messageId", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid messageId")
		return
	}

	//Verifico che il messaggio appartenga all'utente e sia nella conversazione
	isOwner, err := rt.db.IsMessageOwner(userId, messageId, conversationId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Database error while checking message ownership")
		http.Error(writer, "Database error while checking message ownership", http.StatusInternalServerError)
		return
	}

	if !isOwner {
		ctx.Logger.Error("User is not the owner of this message")
		http.Error(writer, "User is not the owner of this message", http.StatusForbidden)
		return
	}

	//Elimino il messaggio dal database
	err = rt.db.DeleteMessage(messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to delete message")
		http.Error(writer, "Failed to delete message", http.StatusInternalServerError)
		return
	}

	//Controllo se la conversazione ora è rimasta senza messaggi
	hasMessages, err := rt.db.DoesConversationHaveMessages(conversationId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Database error while checking conversation messages")
		http.Error(writer, "Database error while checking conversation messages", http.StatusInternalServerError)
		return
	}

	//Se la conversazione è vuota, la elimino
	if !hasMessages {
		err = rt.db.DeleteConversation(conversationId)
		if err != nil {
			ctx.Logger.WithError(err).Error("Failed to delete conversation")
			http.Error(writer, "Failed tocreate conversation", http.StatusInternalServerError)
			return
		}
	}

	//Scrivo la risposta di successo
	writer.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof("message %d is deleted successfully", messageId)
	if !hasMessages {
		ctx.Logger.Infof("conversation %d is deleted succesfully because empty", conversationId)
	}
}
