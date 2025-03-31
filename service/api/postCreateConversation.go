package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postCreateConversation(writer http.ResponseWriter, reader *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

	//Estraggo lo userId dall'Header
	userId, err := extractUserID(reader, ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve userId")
		http.Error(writer, "Failed to retrieve userId", http.StatusUnauthorized)
		return
	}

	//Decodifico la richiesta JSON
	var requestBody struct {
		Name    string `json:"name,omitempty"`
		Photo   string `json:"photo,omitempty"`
		Members []int  `json:"members"`
	}

	if err := json.NewDecoder(reader.Body).Decode(&requestBody); err != nil {
		ctx.Logger.WithError(err).Error("Invalid JSON format")
		http.Error(writer, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	//Analizzo la richiesta e creo la conversazione
	conversationId, err := rt.db.CreateConversation(userId, requestBody.Name, requestBody.Photo, requestBody.Members)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to create conversation")
		http.Error(writer, "Failed to create conversation", http.StatusInternalServerError)
		return
	}

	//Scrivo la risposta, che sarà in formato JSON
	response := struct {
		ConversationId int `json:"conversationId"`
	}{ConversationId: conversationId}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(writer).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Impossible to create Response")
		http.Error(writer, "Impossible to Create Response", http.StatusInternalServerError)
		return
	}

	//Scrivo il successo nel logger
	ctx.Logger.Infof("postCreateConversation completed succesfull: %d", conversationId)

}
