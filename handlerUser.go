package main

import (
	"net/http"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}
