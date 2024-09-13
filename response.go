package main

import (
	"encoding/json"
	"net/http"
)

func respondWithError(w http.ResponseWriter, msg string, code int) {
	type returnError struct {
		Error string `json:"error"`
	}

	dat, err := json.Marshal(returnError{msg})

	w.WriteHeader(code)

	if err != nil {
		w.Write([]byte("Error - something went wrong"))
		return
	}

	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, payload interface{}, code int) {
	dat, err := json.Marshal(payload)

	if err != nil {
		respondWithError(w, "There was an error with your request", 400)
		return
	}

	w.WriteHeader(code)
	w.Write(dat)
}