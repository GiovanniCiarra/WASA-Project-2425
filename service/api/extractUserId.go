package api

import (
	"Progetto/service/api/reqcontext"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

// extractUserID parses the Authorization header in the form "Bearer <userId>"
// and returns the userId as an integer.

func extractUserID(reader *http.Request, ctx reqcontext.RequestContext) (int, error) {

	// Estraggo il token dall'header Authorization
	authHeader := reader.Header.Get("Authorization")
	if authHeader == "" {
		ctx.Logger.Error("Empty Token")
		return -1, errors.New("empty Token")
	}

	// Rimuovo il prefisso "Bearer " dal token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Converto il token in intero, così da ottenere lo userId
	userId, err := strconv.Atoi(tokenString)
	if err != nil {
		ctx.Logger.WithError(err).Error("Invalid Token Format")
		return -1, errors.New("invalid Token Format")
	}

	return userId, nil
}
