package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/thulio63/chirpy/internal/auth"
)

func (cfg *apiConfig)handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	type data struct {
		User_ID uuid.UUID `json:"user_id"`
	}
	type parameters struct{
		Event string `json:"event"`
		Data data `json:"data"`
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid api key", err)
		return
	}
	if key != cfg.apikey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//incoming data
	params := parameters{}
	
	//decode incoming chirp and validate
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	//check event
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	//upgrade user
	err = cfg.db.Upgrade(r.Context(), params.Data.User_ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't find user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}