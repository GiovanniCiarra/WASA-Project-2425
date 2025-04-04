package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postCommentMessage(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

	//Estraggo lo userId dall'Header
	userId, err := extractUserID(reader, ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve userId")
		http.Error(writer, "Failed to retrieve userId", http.StatusUnauthorized)
		return
	}

	//Estraggo il messageId dai parametri della richiesta
	messageId, err := strconv.Atoi(params.ByName("messageId"))
	if err != nil {
		http.Error(writer, "Invalid messageId", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid messageId")
		return
	}

	//Decodifica JSON della richiesta
	var requestBody struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(reader.Body).Decode(&requestBody); err != nil {
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid JSON format")
		return
	}

	if len(requestBody.Content) == 0 || len(requestBody.Content) > 100 {
		http.Error(writer, "Invalid comment content lenght", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid comment content lenght")
		return
	}

	//Controllo se il messaggio esiste e se l'utente è nella conversazione
	isMember, err := rt.db.IsUserInMessageConversation(userId, messageId)
	if err != nil {
		http.Error(writer, "Database error checking user membership", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error checking user membership")
		return
	}
	if !isMember {
		http.Error(writer, "User is not allowed to comment this message", http.StatusBadRequest)
		ctx.Logger.Error("User is not allowed to comment this message")
		return
	}

	//Controllo che non ci siano altri commenti dello stesso utente per lo stesso messaggio
	hasCommented, err := rt.db.HasUserCommentedMessage(userId, messageId)
	if err != nil {
		http.Error(writer, "Database error while cheking comment existance", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Database error while cheking comment existance")
		return
	}

	if hasCommented {
		http.Error(writer, "User has already commented this message", http.StatusConflict)
		ctx.Logger.Error("User has already commented this message")
		return
	}

	//Salvo il nuovo commento nel database
	commentId, err := rt.db.SaveComment(messageId, userId, requestBody.Content)
	if err != nil {
		http.Error(writer, "Failed to save comment", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to save comment")
		return
	}

	//Prendo lo username dell'utente tramite lo userId
	username, _, err := rt.db.GetUserInformationById(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get user information")
		http.Error(writer, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	//Scrivo la risposta
	response := struct {
		CommentId int    `json:"commentId"`
		MessageId int    `json:"messageId"`
		Sender    string `json:"sender"`
		Content   string `json:"content"`
		Timestamp string `json:"timestamp"`
	}{
		CommentId: commentId,
		MessageId: messageId,
		Sender:    username,
		Content:   requestBody.Content,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Failed to create response")
		http.Error(writer, "Failed to create response", http.StatusInternalServerError)
		return
	}

	ctx.Logger.Infof("comment %d is added successfully on message %d", commentId, messageId)
}
