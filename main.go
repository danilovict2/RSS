package main

import (
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		response := struct {
			Status string `json:"status"`
		} {"ok"}
		respondWithJSON(w, response, 200)
	})

	mux.HandleFunc("GET /v1/err", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, "Internal Server Error", 500)
	})
	
	server := http.Server{
		Handler: mux,
		Addr: ":" + os.Getenv("PORT"),
	}
	
	server.ListenAndServe()
}