package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/troyboy95/RSS-project-test/internal/database"
)

func (apiCfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseUserToAPIUser(user))
}

func (apiCfg *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	// This handler is protected by the MiddlewareAuth, so we can assume user is valid
	respondWithJSON(w, http.StatusOK, databaseUserToAPIUser(user))
}

func (apiCfg *apiConfig) getUserPostsHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	// This handler is protected by the MiddlewareAuth, so we can assume user is valid

	posts, err := apiCfg.DB.GetPostForUser(r.Context(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error fetching user posts: %v", err))
		return
	}

	respondWithJSON(w, http.StatusOK, databasePostsToAPIPosts(posts))
}
