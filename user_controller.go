package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/danilovict2/RSS/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, "Couldn't decode parameters", http.StatusInternalServerError)
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: params.Name,
	})

	if err != nil {
		respondWithError(w, "Couldn't create user", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, user, http.StatusOK)
}

func (cfg *apiConfig) getUserByApiKey(w http.ResponseWriter, r *http.Request) {
	apiKey, found := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")
	if !found {
		respondWithError(w, "Please provide a token", http.StatusUnauthorized)
		return
	}

	user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, "Couldn't find user", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, user, http.StatusOK)
}
