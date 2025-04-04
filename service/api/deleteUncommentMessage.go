package api

import (
	"Progetto/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteUncommentMessage(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

	//Estraggo lo userId dall'Header
	userId, err := extractUserID(reader, ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve userId")
		http.Error(writer, "Failed to retrieve userId", http.StatusUnauthorized)
		return
	}

	//Estraggo il messageId dai parametri della richiesta
	messageId, err := strconv.Atoi(params.ByName("messageId"))
	if err != nil {
		http.Error(writer, "Invalid messageId", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid messageId")
		return
	}

	//Controllo che il commento esista e sia dell'utente
	hasCommented, err := rt.db.HasUserCommentedMessage(userId, messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Database error while checking comment existane")
		http.Error(writer, "Database error while checking comment existance", http.StatusInternalServerError)
		return
	}

	if !hasCommented {
		ctx.Logger.Error("User has no comment on this message")
		http.Error(writer, "User has no comment on this message", http.StatusNotFound)
		return
	}

	//Elimino il commento
	err = rt.db.DeleteComment(messageId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to delete comment")
		http.Error(writer, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	//Scrivo la risposta di successo
	writer.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof("comment is deleted successfully from message %d", messageId)

}
