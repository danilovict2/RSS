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

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		respondWithError(w, "Couldn't create feed follow", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, struct{
		Feed database.Feed `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{feed, feedFollow}, http.StatusOK)
}

func (cfg *apiConfig) getFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := cfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, "Couldn't get all feeds", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, feeds, http.StatusOK)
}

func (cfg *apiConfig) createFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, "Couldn't decode parameters", http.StatusInternalServerError)
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID: user.ID,
		FeedID: params.FeedID,
	})

	if err != nil {
		respondWithError(w, "Couldn't create feed follow", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, feedFollow, http.StatusOK)
}

func (cfg *apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	feedFollowID, err := uuid.Parse(r.PathValue("feedFollowID"))
	if err != nil {
		respondWithError(w, "Couldn't parse UUID", http.StatusInternalServerError)
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), feedFollowID)
	if err != nil {
		respondWithError(w, "Couldn't delete feed follow", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, struct{}{}, http.StatusNoContent)
}

func (cfg *apiConfig) getFeedFollowsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollowsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, "Couldn't get feed follows", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, feedFollows, http.StatusOK)
}