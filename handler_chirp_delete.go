package main

import (
	"net/http"

	"github.com/google/uuid"
	auth "github.com/thulio63/chirpy/internal/auth"
)

func (cfg *apiConfig)handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	//make sure user has JWT
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid authorization token", err)
		return
	}
	validatedUID, err := auth.ValidateJWT(token, cfg.secret) 
	if err != nil  {
		respondWithError(w, http.StatusUnauthorized, "invalid JWT", err)
		return
	}

	//find requested chirp
	myid := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(myid)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp id", err)
		return
	}
	myChirp, err := cfg.db.RetrieveChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "couldn't find chirp", err)
		return
	}

	//check if user is author of chirp
	type response struct {}
	if myChirp.UserID != validatedUID {
		respondWithJSON(w, http.StatusForbidden, response{})
		return
	}

	//delete chirp
	err = cfg.db.DeleteChirp(r.Context(), myChirp.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to delete chirp", err)
		return
	}
	respondWithJSON(w, http.StatusNoContent, response{})
}