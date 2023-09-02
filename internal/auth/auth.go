package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Get API Key
func GetApiKey(headers http.Header)(string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("no authentication found")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2 {
		return "", errors.New("malformed authentication header")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of the authentication header")
	}

	return vals[1], nil
}

