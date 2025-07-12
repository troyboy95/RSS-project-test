package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("no authorization found, API key is required")
	}

	vals := strings.Split(apiKey, " ")

	if len(vals) != 2 {
		return "", errors.New("invalid API key format'")
	}

	if vals[0] != "ApiKey" {
		return "", errors.New("invalid API key format in the first part, expected 'ApiKey'")
	}

	return vals[1], nil

}
