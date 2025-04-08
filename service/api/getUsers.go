package api

import (
	"Progetto/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUsers(writer http.ResponseWriter, reader *http.Request, _ httprouter.Params, ctx reqcontext.RequestContext) {

	//Estraggo lo userId dall'Header
	userId, err := extractUserID(reader, ctx)
	if err != nil {
		http.Error(writer, "Failed to retrieve userId", http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("Failed to retrieve userId")
		return
	}

	//Controllo che lo userId sia associato ad un utente valido
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

	//Richiamo la funzione del database per ottenere la lista degli utenti
	UsersList, err := rt.db.GetUsers(userId)

	if err != nil {
		ctx.Logger.WithError(err).Error("Impossible to get Users")
		http.Error(writer, "Impossible to get Users", http.StatusInternalServerError)
		return
	}

	//Creo la Risposta
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	encoder := json.NewEncoder(writer)
	err = encoder.Encode(UsersList)

	if err != nil {
		ctx.Logger.WithError(err).Error("Impossible to Create Response")
		http.Error(writer, "Impossible to Create Response", http.StatusInternalServerError)
		return
	}

	//Scrivo il successo nel logger
	ctx.Logger.Infof("getUsers completed succesfully")
}
