package api

import (
	"Progetto/service/api/reqcontext"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

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

	// Leggo l'immagine
	img, err := io.ReadAll(reader.Body)
	if err != nil {
		http.Error(writer, "Failed to read image", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Failed to read image")
		return
	}

	if len(img) == 0 || len(img) > 1_000_000 {
		http.Error(writer, "Invalid image size", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid image size")
		return
	}

	// Salvo la nuova immagine del database
	err = rt.db.UpdateGroupPhoto(conversationId, img)
	if err != nil {
		http.Error(writer, "Failed to update group photo", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Failed to updat group photo")
	}

	// Scrivo la risposta di successo
	writer.WriteHeader(http.StatusNoContent)

	ctx.Logger.Infof("group's photo is updated successfully")

}
