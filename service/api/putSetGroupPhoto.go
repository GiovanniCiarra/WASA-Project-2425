package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// putSetGroupPhoto updates the photo of a group conversation.
// It validates the request and ensures the image is in base64 PNG format.

func (rt *_router) putSetGroupPhoto(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

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
		http.Error(writer, "Cannot change photo of a private chat", http.StatusBadRequest)
		ctx.Logger.Error("Cannot change photo of a private chat")
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

	// Decodifico la richiesta
	var requestBody struct {
		Photo string `json:"photo"`
	}

	if err := json.NewDecoder(reader.Body).Decode(&requestBody); err != nil {
		http.Error(writer, "Invalid JSON Format", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid JSON Format")
		return
	}

	// Controllo che sia un'immagine base64 e non troppo grande
	validImage := validateBase64Image(requestBody.Photo)

	if !validImage {
		http.Error(writer, "Invalid Image", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid Image")
		return
	}

	// Salvo la nuova immagine del database
	err = rt.db.UpdateGroupPhoto(conversationId, requestBody.Photo)
	if err != nil {
		http.Error(writer, "Failed to update group photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Failed to update group photo")
	}

	// Scrivo la risposta di successo
	writer.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof("group's photo is updated successfully")

}
