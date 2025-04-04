package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postAddToGroup(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

	//Estraggo lo userId dall'Header
	userId, err := extractUserID(reader, ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve userId")
		http.Error(writer, "Failed to retrieve userId", http.StatusUnauthorized)
		return
	}

	//Estraggo il conversationId dai parametri della richiesta
	conversationId, err := strconv.Atoi(params.ByName("conversationId"))
	if err != nil {
		http.Error(writer, "Invalid conversationId", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid conversationId")
		return
	}

	//Decodifico la richiesta JSON
	var requestBody struct {
		NewUserId int `json:"userId"`
	}

	if err := json.NewDecoder(reader.Body).Decode(&requestBody); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid JSON format")
		return
	}

	// Controllo che il nuovo utente esista
	exists, err := rt.db.DoesUserExist(requestBody.NewUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Database error while checking if user exists")
		http.Error(writer, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(writer, "User does not exist", http.StatusNotFound)
		ctx.Logger.Error("User does not exist")
		return
	}

	//Controllo che la conversazione sia un gruppo
	isGroup, err := rt.db.IsGroup(conversationId)
	if err != nil {
		http.Error(writer, "Database error while checking conversation type", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while checking conversation type")
		return
	}

	if !isGroup {
		http.Error(writer, "Cannot add user to a private conversation", http.StatusForbidden)
		ctx.Logger.Error("Cannot add user to a private conversation")
		return
	}

	//Controllo che l'utente chiamante si trovi nel gruppo
	isMember, err := rt.db.IsUserInConversation(userId, conversationId)
	if err != nil {
		http.Error(writer, "Database error while checking membership", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while checking membership")
		return
	}

	if !isMember {
		http.Error(writer, "User is not a member of this group", http.StatusForbidden)
		ctx.Logger.Error("User is not a member of this group")
		return
	}

	//Controllo che l'utente da aggiungere non si trovi già nel gruppo
	isMember, err = rt.db.IsUserInConversation(requestBody.NewUserId, conversationId)
	if err != nil {
		http.Error(writer, "Database error while checking new user membership", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while checking new user membership")
		return
	}

	if isMember {
		http.Error(writer, "User is already in the group", http.StatusConflict)
		ctx.Logger.Error("User is already in the group")
		return
	}

	//Aggiungo il nuovo utente al gruppo
	err = rt.db.AddToGroup(requestBody.NewUserId, conversationId)
	if err != nil {
		http.Error(writer, "Failed to add user to group", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to add user to group")
		return
	}

	//Scrivo al risposta
	writer.WriteHeader(http.StatusNoContent)

	//Scrivo il successo nel logger
	ctx.Logger.Infof("user %d is added successfully to coversation %d", requestBody.NewUserId, conversationId)
}
