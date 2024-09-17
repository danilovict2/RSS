package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/danilovict2/RSS/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("CONNECTION_STRING")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	cfg := apiConfig{
		dbQueries,
	}
	
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

	mux.HandleFunc("POST /v1/users", cfg.createUser)
	mux.HandleFunc("GET /v1/users", cfg.getUserByApiKey)
	
	server := http.Server{
		Handler: mux,
		Addr: ":" + os.Getenv("PORT"),
	}
	
	server.ListenAndServe()
}