package main

import (
	"net/http"
	"strings"
)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, found := strings.CutPrefix(r.Header.Get("Authorization"), "ApiKey ")
		if !found {
			respondWithError(w, "Please provide an api key", http.StatusUnauthorized)
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, "Invalid api key", http.StatusBadRequest)
			return
		}

		handler(w, r, user)
	})
}