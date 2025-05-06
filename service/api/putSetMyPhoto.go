package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

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

	// Controllo che sia un'immagine base64
	if !strings.HasPrefix(requestBody.Photo, "data:image/png") {
		http.Error(writer, "Formato non valido, solo .png", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Invalid image Format")
		return
	}

	parts := strings.SplitN(requestBody.Photo, ",", 2)
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
		http.Error(writer, "L'immagine selezionata è troppo grande (Max 500 KB)", http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Image too large")
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
	ctx.Logger.Infof("Photo Updated succesfully")

}
