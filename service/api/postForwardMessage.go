package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postForwardMessage(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Decodifico al richiesta
	var requestBody struct {
		TargetChatId int `json:"targetChatId"`
	}

	if err := json.NewDecoder(reader.Body).Decode(&requestBody); err != nil {
		http.Error(writer, "Invalid JSON Format", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid JSON Format")
		return
	}

	// Estraggo il messageId dai parametri della richiesta
	messageId, err := strconv.Atoi(params.ByName("messageId"))
	if err != nil {
		http.Error(writer, "Invalid messageId", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid messageId")
		return
	}

	// Estraggo il contenuto del messaggio tramite il messageId
	originalMessage, err := rt.db.GetMessageById(messageId)
	if err != nil {
		http.Error(writer, "Message not found or not accesible", http.StatusForbidden)
		ctx.Logger.WithError(err).Error("Message not found or not accesible")
		return
	}

	// Verifico che l'utente faccia parte della conversazione originale
	isMember, err := rt.db.IsUserInConversation(userId, originalMessage.ConversationId)
	if err != nil {
		http.Error(writer, "Database error while checking original conversation membership", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while checking original conversation membership")
		return
	}
	if !isMember {
		http.Error(writer, "User is not a member of the message original conversation", http.StatusForbidden)
		ctx.Logger.Error("User is not a member of the message original conversation")
		return
	}

	// Verifico che l'utente faccia parte della conversazione di destinazione
	isMember, err = rt.db.IsUserInConversation(userId, requestBody.TargetChatId)
	if err != nil {
		http.Error(writer, "Database error while checking destination membership", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while checking destination membership")
		return
	}
	if !isMember {
		http.Error(writer, "User is not a member of the destination conversation", http.StatusForbidden)
		ctx.Logger.Error("User is not a member of the destination conversation")
		return
	}

	// Copio il messaggio nella nuova conversazione
	newMessageId, err := rt.db.SaveMessage(requestBody.TargetChatId, userId, 1, originalMessage.Content, originalMessage.Type)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to forward message")
		http.Error(writer, "Failed to forward message", http.StatusInternalServerError)
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
	response := struct {
		MessageId      int    `json:"messageId"`
		ConversationId int    `json:"conversationId"`
		Sender         string `json:"sender"`
		Content        string `json:"content"`
		Type           string `json:"type"`
		Timestamp      string `json:"timestamp"`
		Status         string `json:"status"`
		IsForwarded    bool   `json:"isForwarded"`
	}{
		MessageId:      newMessageId,
		ConversationId: requestBody.TargetChatId,
		Sender:         username,
		Content:        originalMessage.Content,
		Type:           originalMessage.Type,
		Timestamp:      time.Now().UTC().Format(time.RFC3339),
		Status:         "sent",
		IsForwarded:    true,
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Failed to create response")
		http.Error(writer, "Failed to create response", http.StatusInternalServerError)
		return
	}

	// Scrivo il successo dell'operazione
	ctx.Logger.Infof("Message %d forwarded to conversation %d successfully", messageId, requestBody.TargetChatId)
}
