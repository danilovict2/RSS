package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/danilovict2/RSS/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, "Couldn't decode parameters", http.StatusInternalServerError)
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: params.Name,
		Url: params.Url,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, "Couldn't create feed", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, feed, http.StatusOK)
}

func (cfg *apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, "Couldn't get all feeds", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, feeds, http.StatusOK)
}