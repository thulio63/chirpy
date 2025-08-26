package main

import (
	"encoding/json"
	"net/http"
	"time"

	auth "github.com/thulio63/chirpy/internal/auth"
	"github.com/thulio63/chirpy/internal/database"
)

func (cfg *apiConfig)handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}  
	type response struct {
		User
	}

	//prepare space for incoming data
	userLogin := parameters{}
	//decode json data
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userLogin)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	//Check for existing token

	data, err := cfg.db.Login(r.Context(), userLogin.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(userLogin.Password, data.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	usr_token, err := auth.MakeJWT(data.ID, cfg.secret, time.Second * time.Duration(3600))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to generate token", err)
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to generate refresh token", err)
		return
	}

	ref_token, err := cfg.db.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
		Token: refresh_token,
		UserID: data.ID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to store refresh token in database", err)
		return
	}

	user := User{
		ID: ref_token.UserID,
		Created_At: data.CreatedAt,
		Updated_At: data.UpdatedAt,
		Email: userLogin.Email,
		Token: usr_token,
		RefreshToken: ref_token.Token,
	}

	respondWithJSON(w, http.StatusOK, response{
		User: user,
	})
}