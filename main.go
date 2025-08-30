package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/thulio63/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
	platform string
	secret string
	apikey string
}

//takes in a handler, creates a new handler with data from running a function on the old handler, and returns new handler
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	new := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
	return new
}

func main() {
	const fileRoot = "."
	const port = "8080"
	// use brew services start postgresql@15 to boot up database
	// and brew services stop postgresql@15 to close it
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening database:", err)
	}
	dbQueries := database.New(db)


	//create config struct to track data
	config := apiConfig{
		fileserverHits: atomic.Int32{}, 
		db: dbQueries, 
		platform: os.Getenv("PLATFORM"),
		secret: os.Getenv("SECRET"),
		apikey: os.Getenv("POLKA_KEY"),
	}

	//mux is a manager that can handle requests, custom functions, etc
	mux := http.NewServeMux()
	//on request for URL/app/, mux will serve the file index.html at . without the prefix /app
	handler := config.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(fileRoot))))
	mux.Handle("/app/", handler)
	//on request for /XXX endpoint, mux will run the corresponding handler function
	//METHOD before endpoint only allows for that method
	mux.HandleFunc("GET /api/healthz", handlerReady)
	mux.HandleFunc("POST /api/chirps", config.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", config.handlerChirpsRetrieve)
	mux.HandleFunc("GET /api/chirps/{chirpID}", config.handlerChirpRetrieve)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", config.handlerChirpDelete)
	mux.HandleFunc("POST /api/users", config.handlerUserCreate)
	mux.HandleFunc("PUT /api/users", config.handlerUserUpdate)
	mux.HandleFunc("POST /api/login", config.handlerLogin)
	mux.HandleFunc("POST /api/revoke", config.handlerRevoke)
	mux.HandleFunc("POST /api/refresh", config.handlerRefresh)
	mux.HandleFunc("POST /api/polka/webhooks", config.handlerUpgrade)
	mux.HandleFunc("GET /admin/metrics", config.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", config.handlerReset)

	server := &http.Server{Handler: mux, Addr: ":" + port}
	//boots up server and blocks main func until the server closes
	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}