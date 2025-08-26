package auth

import (
	"errors"
	"net/http"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearerString := headers.Get("Authorization")
	if bearerString == "" {
		return "", errors.New("no authorization header found")
	}
	if len(bearerString) < 7 {
		return "", errors.New("improper authorization header found")
	}
	bearerToken := bearerString[7:]
	return bearerToken, nil
}