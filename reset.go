package main

import "net/http"

func (cfg *apiConfig)handlerReset(w http.ResponseWriter, r *http.Request) {
	//ensure this is only accessible via local development
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment"))
		return
	}

	cfg.fileserverHits.Store(0)
	err := cfg.db.Reset(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset the database: " + err.Error()))
		return
	}
	
	//this will be a POST request, and as such doesn't need to add the html headers
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database reset to initial state, and hits reset to 0."))
}