package api

import (
	"regexp"
	"unicode"
)

func isValidName(name string) bool {

	// Controllo la lunghezza
	if len(name) < 3 || len(name) > 16 {
		return false
	}

	// Controllo che contenga solo caratteri alfanumerici
	regex := regexp.MustCompile(`^[a-zA-Z0-9\s]+$`)

	if !regex.MatchString(name) {
		return false
	}

	// Controllo che ci sia almeno un carattere
	for _, char := range name {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			return true
		}
	}

	return false
}
