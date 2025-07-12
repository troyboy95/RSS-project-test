package main

import (
	"fmt"
	"net/http"

	"github.com/troyboy95/RSS-project-test/internal/auth"
	"github.com/troyboy95/RSS-project-test/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) MiddlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)

		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Unauthorized: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, 404, fmt.Sprintf("User not found: %v", err))
			return
		}

		handler(w, r, user)
	}
}
