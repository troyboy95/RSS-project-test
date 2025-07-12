package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/troyboy95/RSS-project-test/internal/database"
)

func (apiCfg *apiConfig) createFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error creating feed follow: %v", err))
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToAPIFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) getFeedFollowsHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error getting feed follows: %v", err))
		return
	}

	apiFeedFollows := databaseFeedFollowsToAPIFeedFollows(feedFollows)

	respondWithJSON(w, http.StatusOK, apiFeedFollows)
}

func (apiCfg *apiConfig) deleteFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowStr := chi.URLParam(r, "FeedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid FeedFollow ID: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error deleting feed follow: %v", err))
		return
	}
	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
