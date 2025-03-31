package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) putSetGroupName(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

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

	//Controllo che la conversazione sia un gruppo
	isGroup, err := rt.db.IsGroup(conversationId)
	if err != nil {
		http.Error(writer, "Database error while checking conversation type", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while checking conversation type")
		return
	}

	if !isGroup {
		http.Error(writer, "Cannot rename a privat chat", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Cannot rename a private chat")
		return
	}

	//Controllo che l'utente sia un membro del gruppo
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

	//Decodifico al richiesta
	var requestBody struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(reader.Body).Decode(&requestBody); err != nil {
		http.Error(writer, "Invalid JSON Format", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid JSON Format")
		return
	}

	//Controllo che il nuovo nome del gruppo sia valido
	if len(requestBody.Name) < 3 || len(requestBody.Name) > 16 {
		http.Error(writer, "Invalid group name length", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid group name length")
		return
	}

	//Aggiorno il nome del gruppo nel database
	err = rt.db.UpdateGroupName(conversationId, requestBody.Name)
	if err != nil {
		http.Error(writer, "Failed to update group name", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to update group name")
		return
	}

	//Scrivo la risposta di successo
	writer.WriteHeader(http.StatusOK)

	ctx.Logger.Infof("putSetGroupName completed successfully")
}
