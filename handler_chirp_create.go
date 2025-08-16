package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thulio63/chirpy/internal/database"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	Created_At time.Time `json:"created_at"`
	Updated_At time.Time `json:"updated_at"`
	Body string `json:"body"`
	User_ID uuid.UUID `json:"user_id"`
}

//address how errors are handled

func (cfg *apiConfig)handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	//incoming data
	params := parameters{}
	
	//decode incoming chirp and validate
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	validated, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
	}
	
	//connect with db to add chirp, chirps.sql fills in data
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: validated,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Chirp", err)
		return
	}

	//in custom function, pass data created from chirp
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID: chirp.ID, 
		Created_At: chirp.CreatedAt, 
		Updated_At: chirp.UpdatedAt, 
		Body: chirp.Body, 
		User_ID: chirp.UserID,
	})
}

func validateChirp(body string) (string, error) {
	const maxLength = 140
	if len(body) > maxLength {
		return "", errors.New("Chirp is too long")
	}

	//check/replace profane
	var bad_words []string
	bad_words = append(bad_words, "kerfuffle", "sharbert", "fornax")
	words := strings.Split(body, " ")
	for i, word := range words {
		for _, bad := range bad_words {
			if strings.ToLower(word) == bad {
				words[i] = "****"
			}
		}
	}
	clean := strings.Join(words, " ")

	return clean, nil
}