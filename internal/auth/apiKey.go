package auth

import (
	"errors"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiStr := headers.Get("Authorization")
	if apiStr == "" {
		return "", errors.New("no authorization header found")
	}
	if len(apiStr) < 7 {
		return "", errors.New("improper authorization header found")
	}
	apiKey := apiStr[7:]
	return apiKey, nil
}