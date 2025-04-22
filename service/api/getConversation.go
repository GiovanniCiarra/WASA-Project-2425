package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversation(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

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
		http.Error(writer, "error checking conversation membership", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("error checking conversation membership")
		return
	}

	if !isMember {
		http.Error(writer, "user not part of the conversation", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("user not part of the conversation")
		return
	}

	// Richiamo la funzione del database per ottenere i messaggi
	messages, err := rt.db.GetConversation(userId, conversationId)
	if err != nil {
		http.Error(writer, "Failed to retrieve messages", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to retrieve messages")
		return
	}

	// Aggiorno lo stato dei messaggi a "read" per questo utente
	err = rt.db.UpdateMessagesToRead(userId, conversationId)
	if err != nil {
		http.Error(writer, "Failed to update message status", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to update message status")
		return
	}

	// Scrivo la Risposta
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(writer).Encode(messages); err != nil {
		ctx.Logger.WithError(err).Error("Failed to Create Response")
		http.Error(writer, "Failed to Create Response", http.StatusInternalServerError)
		return
	}

	// Scrivo il successo nel logger
	ctx.Logger.Infof("GetConversation completed successfully")
}
