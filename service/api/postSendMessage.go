package api

import (
	"Progetto/service/api/reqcontext"
	"Progetto/service/database"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postSendMessage(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Decodifico il corpo JSON della richiesta
	var requestBody struct {
		Content string `json:"content"`
		Type    string `json:"type"`
	}

	if err := json.NewDecoder(reader.Body).Decode(&requestBody); err != nil {
		ctx.Logger.WithError(err).Error("Invalid JSON format")
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Controllo i dati del requestBody
	if requestBody.Type != "text" && requestBody.Type != "photo" {
		ctx.Logger.WithError(err).Error("Invalid message type")
		http.Error(writer, "Invalid message type", http.StatusBadRequest)
		return
	}

	//Se un immagine, controllo che sia in base64 e non troppo pesante
	if requestBody.Type == "photo" {
		if !strings.HasPrefix(requestBody.Content, "data:image/png") {
			http.Error(writer, "Invalid image Format", http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("Invalid image Format")
			return
		}

		parts := strings.SplitN(requestBody.Content, ",", 2)
		if len(parts) != 2 {
			http.Error(writer, "invalid base64 format", http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("Invalid base64 Format")
			return
		}

		data, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			http.Error(writer, "invalid base64 encoding", http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("Invalid base64 encoding")
			return
		}

		if len(data) > 500*1024 {
			http.Error(writer, "image too large", http.StatusBadRequest)
			ctx.Logger.WithError(err).Error("Image too large")
			return
		}
	}
	// Controllo che l'utente partecipi alla conversazione
	isMember, err := rt.db.IsUserInConversation(userId, conversationId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Database error while checking membership")
		http.Error(writer, "Database error while checking membership", http.StatusInternalServerError)
		return
	}

	if !isMember {
		ctx.Logger.WithError(err).Error("User is not a member of the conversation")
		http.Error(writer, "User is not a member of the conversation", http.StatusForbidden)
		return
	}

	// Inserisco il messaggio nel database
	messageId, err := rt.db.SaveMessage(conversationId, userId, 0, -1, requestBody.Content, requestBody.Type)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to save message")
		http.Error(writer, "Failed to save message", http.StatusInternalServerError)
		return
	}

	// Prendo lo username dell'utente tramite lo userId
	username, _, err := rt.db.GetUserInformationById(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get user information")
		http.Error(writer, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Creo la risposta
	var response database.Message
	response.MessageId = messageId
	response.ConversationId = conversationId
	response.Type = requestBody.Type
	response.Content = requestBody.Content
	response.SenderId = userId
	response.SenderUsername = username
	response.Status = "sent"
	response.Timestamp = time.Now().UTC().Format(time.RFC3339)
	response.IsForwarded = 0

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Failed to create response")
		http.Error(writer, "Failed to create response", http.StatusInternalServerError)
		return
	}

	// Scrivo il successo nel Logger
	ctx.Logger.Infof("postSendMessage completed succesfully: %d", messageId)
}
