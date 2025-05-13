package api

import (
	"encoding/base64"
	"strings"
)

func validateBase64Image(photo string) bool {

	if !strings.HasPrefix(photo, "data:image/png") {
		return false
	}

	parts := strings.SplitN(photo, ",", 2)
	if len(parts) != 2 {
		return false
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	if len(data) > 500*1024 {
		return false
	}

	return true
}
