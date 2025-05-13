package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// postCreateConversation creates a new group or private conversation.
// It validates the current user, checks the members, and optionally verifies the group photo.

func (rt *_router) postCreateConversation(writer http.ResponseWriter, reader *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Decodifico la richiesta JSON
	var requestBody struct {
		Name    string `json:"name,omitempty"`
		Photo   string `json:"photo,omitempty"`
		Members []int  `json:"members"`
	}

	if err := json.NewDecoder(reader.Body).Decode(&requestBody); err != nil {
		ctx.Logger.WithError(err).Error("Invalid JSON format in request")
		http.Error(writer, "Invalid JSON format in request", http.StatusBadRequest)
		return
	}

	// Controllo che i memebri siano utenti validi
	for _, member := range requestBody.Members {
		ok, err := rt.db.DoesUserExist(member)
		if err != nil {
			ctx.Logger.WithError(err).Error("Impossible to validate members")
			http.Error(writer, "Impossible to validate members", http.StatusInternalServerError)
			return
		}

		if !ok {
			ctx.Logger.Error("member userId is not valid")
			http.Error(writer, "member userId is not valid", http.StatusBadRequest)
			return
		}
	}

	// Se c'è la foto, controllo che sia valida in base 64 e non troppo grande
	if requestBody.Photo != "" {
		validImage := validateBase64Image(requestBody.Photo)

		if !validImage {
			http.Error(writer, "Invalid Image", http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("Invalid Image")
			return
		}
	}

	// Analizzo la richiesta e creo la conversazione.
	conversationId, err := rt.db.CreateConversation(userId, requestBody.Name, requestBody.Photo, requestBody.Members)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create conversation")
		http.Error(writer, "Failed to create conversation", http.StatusInternalServerError)
		return
	}

	// Scrivo la risposta, che sarà in formato JSON
	response := struct {
		ConversationId int `json:"conversationId"`
	}{ConversationId: conversationId}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(writer).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Failed to create Response")
		http.Error(writer, "Failed to Create Response", http.StatusInternalServerError)
		return
	}

	// Scrivo il successo nel logger
	ctx.Logger.Infof("postCreateConversation completed successfully: %d", conversationId)

}
