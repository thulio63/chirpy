package main

import (
	"encoding/json"
	"net/http"

	auth "github.com/thulio63/chirpy/internal/auth"
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

	user := User{
		ID: data.ID,
		Created_At: data.CreatedAt,
		Updated_At: data.UpdatedAt,
		Email: userLogin.Email,
	}

	respondWithJSON(w, http.StatusOK, response{
		User: user,
	})
}