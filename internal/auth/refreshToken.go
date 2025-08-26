package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		errorText := fmt.Errorf("error with rand.Read(): %w", err)
		return "", errorText
	}
	randStr := hex.EncodeToString(key)
	return randStr, nil
}