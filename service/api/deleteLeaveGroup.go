package api

import (
	"Progetto/service/api/reqcontext"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteLeaveGroup(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

	//Estraggo lo userId dall'Header
	userId, err := extractUserID(reader, ctx)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	}

	//Estraggo il conversationId dai parametri della richiesta
	conversationId, err := strconv.Atoi(params.ByName("conversationId"))
	if err != nil {
		http.Error(writer, "Invalid conversationId", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid conversationId")
		return
	}

	//Controllo se la conversazione è un gruppo
	isGroup, err := rt.db.IsGroup(conversationId)
	if err != nil {
		http.Error(writer, "Database error while checking conversation type", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while checking converation type")
		return
	}

	if !isGroup {
		http.Error(writer, "Cannot leave a private chat", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Cannot leave a private chat")
		return
	}

	//Controllo che l'utente sia membro del gruppo
	isMember, err := rt.db.IsUserInConversation(userId, conversationId)
	if err != nil {
		http.Error(writer, "Database error while checking membership", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while checking membership")
		return
	}

	if !isMember {
		http.Error(writer, "User is not a member of this group", http.StatusForbidden)
		ctx.Logger.WithError(err).Error("User is not a member of this group")
		return
	}

	//Rimuovo l'utente dal gruppo
	err = rt.db.RemoveFromGroup(userId, conversationId)
	if err != nil {
		http.Error(writer, "Failed to remove user from group", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to remove user from group")
		return
	}

	//Controllo se il gruppo è rimasto senza membri
	hasMembers, err := rt.db.DoesGroupHaveMembers(conversationId)
	if err != nil {
		http.Error(writer, "Database error while checking group members", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while checking group members")
		return
	}

	//Se il gruppo è vuoto, lo elimino
	if !hasMembers {
		err = rt.db.DeleteConversation(conversationId)
		if err != nil {
			http.Error(writer, "Failed to delete group", http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Failed to delete group")
			return
		}
	}

	//Scrivo la risposta di successo
	writer.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof("deleteLeaveGroup completed successfully")

}
