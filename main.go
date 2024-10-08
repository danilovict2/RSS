package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danilovict2/RSS/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

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
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.getUserByApiKey))

	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.createFeed))
	mux.HandleFunc("GET /v1/feeds", cfg.getFeeds)
	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.createFeedFollow))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.deleteFeedFollow)
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.getFeedFollowsForUser))

	mux.HandleFunc("GET /v1/posts", cfg.middlewareAuth(cfg.getPostsByUser))
	
	server := http.Server{
		Handler: mux,
		Addr: ":" + os.Getenv("PORT"),
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)
	
	log.Printf("Serving on port: %s\n", os.Getenv("PORT"))
	log.Fatal(server.ListenAndServe())
}