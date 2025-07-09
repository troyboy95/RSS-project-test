package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", payload)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)

}

func respondWithError(w http.ResponseWriter, status int, message string) {

	if status > 499 {
		log.Println("Server error:", message)
	}

	type ErrorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, status, ErrorResponse{Error: message})
}
