package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/troyboy95/RSS-project-test/internal/database"
)

func (apiCfg *apiConfig) createFeedHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		UserID:    user.ID,
		Url:       params.URL, // Assuming the URL is the same as the name for simplicity
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating feed: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedToAPIFeed(feed))
}

func (apiCfg *apiConfig) getFeedsHandler(w http.ResponseWriter, r *http.Request) {

	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error fetching feeds: %v", err))
		return
	}

	respondWithJSON(w, http.StatusFound, databaseFeedsToAPIFeeds(feeds))
}
