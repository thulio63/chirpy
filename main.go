package main

import (
	"net/http"
)

func main() {
	const fileRoot = "."
	const port = "8080"

	//mux is a manager that can handle requests, custom functions, etc
	mux := http.NewServeMux()
	//on request for URL/app/, mux will serve the file index.html at . without the prefix /app
	mux.Handle("/app/", http.StripPrefix("/app",http.FileServer(http.Dir(fileRoot))))
	//on request for any /healthz endpoint, mux will run the handlerReady function
	mux.HandleFunc("/healthz", handlerReady)

	server := http.Server{Handler: mux, Addr: ":" + port}
	//boots up server and blocks main func until the server closes
	server.ListenAndServe()
}

//creates a header with basic info, returns 200 response, writes 200 response text (OK) to body of page, and serves page
func handlerReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}