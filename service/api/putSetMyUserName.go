package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) putSetMyUserName(writer http.ResponseWriter, reader *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

	//Estraggo lo userId dall'Header
	userId, err := extractUserID(reader, ctx)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to retrieve userId")
		http.Error(writer, "Failed to retrieve userId", http.StatusUnauthorized)
		return
	}

	requestBody := struct {
		Username string `json:"username"`
	}{}

	//Decodifica della richiesta
	decoder := json.NewDecoder(reader.Body)
	err = decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid JSON Format in Request")
		http.Error(writer, "Invalid JSON Format in Request", http.StatusBadRequest)
		return
	}

	//Controllo la validità del nuovo username (compreso tra 3 e 16 caratteri e contenete caratteri alfanumerici(non solo spazi))
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
		ctx.Logger.WithError(err).Error("The username already exists")
		http.Error(writer, "The username already exists", http.StatusBadRequest)
		return
	}

	//Aggiorno il nome utente nel database
	err = rt.db.UpdateUserName(userId, requestBody.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("Impossible to update the UserName")
		http.Error(writer, "Impossible to update the UserName", http.StatusInternalServerError)
		return
	}

	response := struct {
		NewName string `json:"username"`
	}{NewName: requestBody.Username}

	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(response)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	//Scrivo il successo del Login nel logger
	ctx.Logger.Infof("Username Updated succesfully: %s", response.NewName)

}
