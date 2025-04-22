package api

import (
	"Progetto/service/api/reqcontext"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) postDoLogin(writer http.ResponseWriter, reader *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

	requestBody := struct {
		Name string `json:"name"`
	}{}

	// Decodifica della richiesta
	decoder := json.NewDecoder(reader.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid Json Format in request")
		http.Error(writer, "Invalid Json Format in request", http.StatusBadRequest)
		return
	}

	// Controllo la validità dello username (compreso tra 3 e 16 caratteri e contenete caratteri alfanumerici(non solo spazi))
	if !isValidName(requestBody.Name) {
		ctx.Logger.WithError(err).Error("Invalid name")
		http.Error(writer, "Invalid name", http.StatusBadRequest)
		return
	}

	// Controllo se l'utente esiste per ricevere il suo identifier e restituirlo come risposta http
	var userId int
	httpStatusResult := http.StatusOK
	userId, err = rt.db.GetUserIdByName(requestBody.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { //se l'utente non esiste, lo creo
			userId, err = rt.db.CreateUser(requestBody.Name) //quando creo l'utente la funzione ritorna anche l'id
			if err != nil {
				ctx.Logger.WithError(err).Error("User creation failed")
				http.Error(writer, "User Creation Failed", http.StatusInternalServerError)
				return
			}
			httpStatusResult = http.StatusCreated
		} else {
			ctx.Logger.WithError(err).Error("User creation failed")
			http.Error(writer, "User Creation Failed", http.StatusInternalServerError)
			return
		}
	}

	// Scrivo la risposta, che sarà in formato JSON
	response := struct {
		Identifier int `json:"identifier"`
	}{Identifier: userId}

	writer.Header().Set("content-type", "application/json")
	writer.Header().Set("Autenthication", fmt.Sprintf("bearerAuth: %d", userId))
	writer.WriteHeader(httpStatusResult)

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(response)
	if err != nil {
		http.Error(writer, "Failed to encode response", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	// Scrivo il successo del Login nel logger
	ctx.Logger.Infof("Login Succesfull: %s", requestBody.Name)

}
