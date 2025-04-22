package api

import (
	"Progetto/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getConversationByUserId(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Estraggo il targetUserId dai parametri della richiesta
	targetUserId, err := strconv.Atoi(params.ByName("userId"))
	if err != nil {
		http.Error(writer, "Invalid targetUserId", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid targetUserId")
		return
	}

	// Controllo che lo userId e il targetUserId non siano lo stesso
	if userId == targetUserId {
		http.Error(writer, "targetUserId is equal to userId", http.StatusBadRequest)
		ctx.Logger.Error("targetUserId is equal to userId")
		return
	}

	// Controllo che il targetUserId sia associato ad un utente valido
	valid, err = rt.db.DoesUserExist(targetUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Impossible to validate target user")
		http.Error(writer, "Impossible to validate target user", http.StatusInternalServerError)
		return
	}

	if !valid {
		ctx.Logger.Error("targetUserId is not valid")
		http.Error(writer, "targetUserId is not valid", http.StatusBadRequest)
		return
	}

	// Chiamo la funzione del database per controllare se esiste una conversazione privata tra userId e targetuserId
	conversationId, err := rt.db.DoesConversationExists(userId, targetUserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.Logger.WithError(err).Error("the conversation does not exist")
			http.Error(writer, "the conversation does not exist", http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("Impossible to check conversation existence")
		http.Error(writer, "Impossible to check conversation existence", http.StatusInternalServerError)
		return
	}

	// Scrivo la risposta, che sarà in formato JSON
	response := struct {
		ConversationId int `json:"conversationId"`
	}{ConversationId: conversationId}

	// Scrivo la Risposta
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(writer).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Failed to Create Response")
		http.Error(writer, "Failed to Create Response", http.StatusInternalServerError)
		return
	}

	// Scrivo il successo nel logger
	ctx.Logger.Infof("getConversationByUserId completed successfully")
}
