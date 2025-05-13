package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// getConversationMembers returns the list of members for a specific conversation,
// only if the requesting user is a member of that conversation.

func (rt *_router) getConversationMembers(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

	// Estraggo lo userId dall'Header
	userId, err := extractUserID(reader, ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve userId")
		http.Error(writer, "Failed to retrieve userId", http.StatusUnauthorized)
		return
	}

	// Controllo che lo userId sia associato ad un utente valido
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

	// Estraggo il conversationId dai parametri della richiesta
	conversationId, err := strconv.Atoi(params.ByName("conversationId"))
	if err != nil {
		http.Error(writer, "Invalid conversationId", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid conversationId")
		return
	}

	// Controllo che l'utente faccia parte della conversazione
	isMember, err := rt.db.IsUserInConversation(userId, conversationId)
	if err != nil {
		http.Error(writer, "error checking conversation membership", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("error checking conversation membership")
		return
	}

	if !isMember {
		http.Error(writer, "user is not a member of the conversation", http.StatusForbidden)
		ctx.Logger.Error("user is not a member of the conversation")
		return
	}

	// Richiamo la funzione del database per ottenere i membri della conversazione
	members, err := rt.db.GetConversationMembers(conversationId)
	if err != nil {
		http.Error(writer, "Failed to retrieve conversation members", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to retrieve conversation members")
		return
	}

	// Scrivo la Risposta
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(writer).Encode(members); err != nil {
		ctx.Logger.WithError(err).Error("Failed to create response")
		http.Error(writer, "Failed to create response", http.StatusInternalServerError)
		return
	}

	// Scrivo il successo nel logger
	ctx.Logger.Infof("GetConversationMembers completed successfully")
}
