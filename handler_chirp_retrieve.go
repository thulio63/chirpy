package main

import (
	"net/http"

	"github.com/google/uuid"
)

	func (cfg *apiConfig)handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {	

	allChirps := []Chirp{}
	chirps, err := cfg.db.RetrieveChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
	}
	for _, chirp := range chirps {
		newChirp := Chirp{
			ID: chirp.ID,
			Created_At: chirp.CreatedAt,
			Updated_At: chirp.UpdatedAt,
			Body: chirp.Body,
			User_ID: chirp.UserID,
		}
		allChirps = append(allChirps, newChirp)
	}

	respondWithJSON(w, 200, allChirps)
}

func (cfg *apiConfig)handlerChirpRetrieve(w http.ResponseWriter, r *http.Request) {	
	myid := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(myid)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse uuid", err)
	}
	myChirp, err := cfg.db.RetrieveChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
	}
	resp := Chirp{
		ID: myChirp.ID,
		Created_At: myChirp.CreatedAt,
		Updated_At: myChirp.UpdatedAt,
		Body: myChirp.Body,
		User_ID: myChirp.UserID,
	}

	respondWithJSON(w, 200, resp)
}