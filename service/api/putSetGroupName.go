package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) putSetGroupName(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Controllo che la conversazione sia un gruppo
	isGroup, err := rt.db.IsGroup(conversationId)
	if err != nil {
		http.Error(writer, "Database error while checking conversation type", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while checking conversation type")
		return
	}

	if !isGroup {
		http.Error(writer, "Cannot rename a privat chat", http.StatusBadRequest)
		ctx.Logger.Error("Cannot rename a private chat")
		return
	}

	// Controllo che l'utente sia un membro del gruppo
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

	// Decodifico al richiesta
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(reader.Body).Decode(&requestBody); err != nil {
		http.Error(writer, "Invalid JSON Format", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid JSON Format")
		return
	}

	// Controllo la validità del nuovo nome (compreso tra 3 e 16 caratteri e contenete caratteri alfanumerici(non solo spazi))
	if !isValidName(requestBody.Name) {
		ctx.Logger.WithError(err).Error("Invalid name")
		http.Error(writer, "Invalid name", http.StatusBadRequest)
		return
	}

	// Aggiorno il nome del gruppo nel database
	err = rt.db.UpdateGroupName(conversationId, requestBody.Name)
	if err != nil {
		http.Error(writer, "Failed to update group name", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to update group name")
		return
	}

	writer.WriteHeader(http.StatusOK)

	// Scrivo la risposta di successo
	ctx.Logger.Infof("group's name is updated succesfully: %s", requestBody.Name)
}
