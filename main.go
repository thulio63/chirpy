package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thulio63/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	queries *database.Queries
}

func main() {
	// use brew services start postgresql@15 to boot up database
	// and brew services stop postgresql@15 to close it
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error:", err)
	}
	dbQueries := database.New(db)

	const fileRoot = "."
	const port = "8080"

	//create config struct to track data
	config := apiConfig{fileserverHits: atomic.Int32{}, queries: dbQueries}

	//mux is a manager that can handle requests, custom functions, etc
	mux := http.NewServeMux()
	//on request for URL/app/, mux will serve the file index.html at . without the prefix /app
	handler := http.StripPrefix("/app",http.FileServer(http.Dir(fileRoot)))
	mux.Handle("/app/", config.middlewareMetricsInc(handler))
	//on request for /XXX endpoint, mux will run the corresponding handler function
	mux.HandleFunc("/api/healthz", handlerReady)
	mux.HandleFunc("/admin/metrics", config.handlerMetrics)
	mux.HandleFunc("/admin/reset", config.handlerReset)
	mux.HandleFunc("/api/validate_chirp", handlerValidate)

	server := &http.Server{Handler: mux, Addr: ":" + port}
	//boots up server and blocks main func until the server closes
	server.ListenAndServe()
}

//takes in a handler, creates a new handler with data from running a function on the old handler, and returns new handler
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	new := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
	return new
}

//creates a header with basic info, returns 200 response, writes 200 response text (OK) to body of page, and serves page
func handlerReady(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	} else {
		w.WriteHeader(405)
	}
}

func (cfg *apiConfig)handlerMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		//str := strconv.FormatInt(int64(cfg.fileserverHits.Load()), 10)
		text := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())
		w.Write([]byte(text))
	} else {
		w.WriteHeader(405)
	}
}

func (cfg *apiConfig)handlerReset(w http.ResponseWriter, r *http.Request) {
	//this will be a POST request, and as such doesn't need to add the html headers
	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusOK)
		cfg.fileserverHits.Store(0)
		w.Write([]byte("Reset"))
	} else {
		w.WriteHeader(405)
	}
}

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	//incoming chirp to validate
	type chirp struct {
		Body string `json:"body"`
		
	}
	//response sent by server
	type response struct {
		Cleaned_Body string `json:"cleaned_body"`
		//Valid bool `json:"valid"`
		Error string `json:"error"`
	}
	//decode incoming chirp and validate
	decoder := json.NewDecoder(r.Body)
	chirpInc := chirp{}
	err := decoder.Decode(&chirpInc)
	if err != nil {
		fmt.Println("Error decoding chirp:", err)
		w.WriteHeader(500)
		return
	}
	//encode response
	resp := response{Cleaned_Body: "", Error: ""}

	if len(chirpInc.Body) > 140 {
		//resp.Valid = false
		resp.Error = "Chirp is too long"
	}

	//check for profane
	var bad_words []string
	bad_words = append(bad_words, "kerfuffle", "sharbert", "fornax")
	words := strings.Split(chirpInc.Body, " ")
	for i, word := range words {
		for _, bad := range bad_words {
			if strings.ToLower(word) == bad {
				words[i] = "****"
			}
		}
	}
	combined_words := strings.Join(words, " ")
	resp.Cleaned_Body = combined_words

	dat, err := json.Marshal(resp)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		w.WriteHeader(500)
	}

	//send response
	w.Header().Set("Content-Type", "application/json")
	if len(chirpInc.Body) > 140 {
		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)	
	}
	w.Write(dat)

}