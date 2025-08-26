package main

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/thulio63/chirpy/internal/auth"
)

func (cfg *apiConfig)handlerRefresh(w http.ResponseWriter, r *http.Request) {
	headString := r.Header.Get("Authorization")
	if headString == "" {
		respondWithError(w, http.StatusNotFound, "No refresh token provided", errors.New("no refresh token found"))
		return
	}
	
	refreshToken := headString[7:]

	info, err := cfg.db.RetrieveUserByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Provided refresh token is not valid", err)
		return
	}

	now := time.Now()

	if now.After(info.ExpiresAt.Time) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has expired", errors.New("expired refresh token"))
		return
	}

	if info.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token has been revoked", errors.New("revoked refresh token"))
		return
	}

	//after passing test, provide new token to user

	tok, err := auth.MakeJWT(info.UserID, os.Getenv("SECRET"), time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: tok,
	})
	
}

func (cfg *apiConfig)handlerRevoke(w http.ResponseWriter, r *http.Request) { 
	headString := r.Header.Get("Authorization")
	if headString == "" {
		respondWithError(w, http.StatusNotFound, "No refresh token provided", errors.New("no refresh token found"))
		return
	}
	
	refreshToken := headString[7:]

	_, err := cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token", errors.New("failed to revoke token"))
	}

	type response struct {

	}
	respondWithJSON(w, http.StatusNoContent, response{})
}