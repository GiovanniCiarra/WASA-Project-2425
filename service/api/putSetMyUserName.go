package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// putSetMyUserName updates the username of the currently user.
// It validates the new name and ensures it is not already taken by another user.

func (rt *_router) putSetMyUserName(writer http.ResponseWriter, reader *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

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

	requestBody := struct {
		Username string `json:"username"`
	}{}

	// Decodifica della richiesta
	decoder := json.NewDecoder(reader.Body)
	err = decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid JSON Format in Request")
		http.Error(writer, "Invalid JSON Format in Request", http.StatusBadRequest)
		return
	}

	// Controllo la validità del nuovo username (compreso tra 3 e 16 caratteri e contenente caratteri alfanumerici(non solo spazi))
	if !isValidName(requestBody.Username) {
		ctx.Logger.WithError(err).Error("Invalid name")
		http.Error(writer, "Invalid name", http.StatusBadRequest)
		return
	}

	// controllo che il nuovo username non sia già associato ad un altro Utente
	check, err := rt.db.CheckUsernameExists(requestBody.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error Checking the username")
		http.Error(writer, "Error Checking the username", http.StatusInternalServerError)
		return
	}
	if check {
		ctx.Logger.WithError(err).Error("The selected username already exists")
		http.Error(writer, "The selected username already exists", http.StatusBadRequest)
		return
	}

	// Aggiorno il nome utente nel database
	err = rt.db.UpdateUserName(userId, requestBody.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to update the Username")
		http.Error(writer, "Failed to update the Username", http.StatusInternalServerError)
		return
	}

	// Scrivo la risposta di successo
	writer.WriteHeader(http.StatusNoContent)

	// Scrivo il successo del Login nel logger
	ctx.Logger.Infof("Username Updated successfully: %s", requestBody.Username)

}
