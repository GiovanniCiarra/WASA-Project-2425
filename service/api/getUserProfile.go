package api

import (
	"Progetto/service/api/reqcontext"
	"Progetto/service/database"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// getUserProfile returns the authenticated user's profile information (userId, username, and photo).

func (rt *_router) getUserProfile(writer http.ResponseWriter, reader *http.Request, params httprouter.Params, ctx reqcontext.RequestContext) {

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

	// Richiamo la funzione del database per ricevere le informazioni dell'utente
	var username, photo string
	username, photo, err = rt.db.GetUserInformationById(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to get user information")
		http.Error(writer, "Failed to get user information", http.StatusInternalServerError)
		return
	}

	// Creo la risposta
	var response database.User
	response.UserId = userId
	response.Name = username
	response.Photo = photo

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(writer).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Failed to create response")
		http.Error(writer, "Failed to create response", http.StatusInternalServerError)
		return
	}

	// Scrivo il successo nel logger
	ctx.Logger.Infof("getUserProfile completed successfully")
}
