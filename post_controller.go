package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/danilovict2/RSS/internal/database"
)

const DEFAULT_LIMIT = 10

func (cfg *apiConfig) getPostsByUser(w http.ResponseWriter, r *http.Request, user database.User) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		log.Println("Couldn't get limit: %v", err)
		limit = DEFAULT_LIMIT
	}

	posts, err := cfg.DB.GetPostsByUser(r.Context(), database.GetPostsByUserParams{
		UserID: user.ID,
		Limit: int32(limit),
	})

	if err != nil {
		respondWithError(w, "Couldn't fetch posts", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, posts, http.StatusOK)
}