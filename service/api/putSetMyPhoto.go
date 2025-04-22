package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) putSetMyPhoto(writer http.ResponseWriter, reader *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

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
		Photo string `json:"photo"`
	}{}

	// Decodifica della richiesta
	decoder := json.NewDecoder(reader.Body)
	err = decoder.Decode(&requestBody)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid Json format in request")
		http.Error(writer, "Invalid Json format in request", http.StatusBadRequest)
		return
	}

	// controllo la validità della nuova foto
	_, err = base64.StdEncoding.DecodeString(requestBody.Photo)
	if err != nil {
		http.Error(writer, "Invalid Photo Format", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid Photo Format")
		return
	}

	// Aggiorno la photo nel database
	err = rt.db.UpdateUserPhoto(userId, requestBody.Photo)
	if err != nil {
		http.Error(writer, "Impossible to update photo", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Impossible to update photo")
		return
	}

	// Scrivo la risposta
	response := struct {
		NewPhoto string `json:"photo"`
	}{NewPhoto: requestBody.Photo}

	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(response)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Failed to encode response")
		return
	}

	// Scrivo il successo del Login nel logger
	ctx.Logger.Infof("Photo Updated succesfully: %s", response.NewPhoto)

}
