package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getMyConversations retrieves all conversations for the user.
// It updates the status of messages to 'received' and returns the list as JSON.

func (rt *_router) getMyConversations(writer http.ResponseWriter, reader *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Richiamo la funzione del database per ottenere la lista delle chat
	chatList, err := rt.db.GetMyConversations(userId)

	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve conversations")
		http.Error(writer, "Failed to retrieve conversations", http.StatusInternalServerError)
		return
	}

	// Aggiorno lo stato dei messaggi a "received" per questo utente
	err = rt.db.UpdateMessagesToReceived(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to update message status")
		http.Error(writer, "Failed to update message status", http.StatusInternalServerError)
		return
	}

	// Scrivo la Risposta
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(writer).Encode(chatList); err != nil {
		ctx.Logger.WithError(err).Error("Impossible to create response")
		http.Error(writer, "Impossible to create Response", http.StatusInternalServerError)
		return
	}

	// Scrivo il successo nel logger
	ctx.Logger.Infof("GetMyConversations completed successfully")

}
