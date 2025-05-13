package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// putSetMyPhoto updates the profile photo of the user.
// The image must be in base64-encoded PNG format and no larger than 500 KB.

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

	// Controllo che sia un immagine in base64 e non troppo pesante
	validImage := validateBase64Image(requestBody.Photo)

	if !validImage {
		http.Error(writer, "Invalid Image", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid Image")
		return
	}

	// Aggiorno la photo nel database
	err = rt.db.UpdateUserPhoto(userId, requestBody.Photo)
	if err != nil {
		http.Error(writer, "Impossible to update photo", http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Impossible to update photo")
		return
	}

	// Scrivo la risposta di successo
	writer.WriteHeader(http.StatusNoContent)

	// Scrivo il successo del Login nel logger
	ctx.Logger.Infof("Photo Updated successfully")

}
