package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/iamarju/rss-aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)

	if err != nil {
		responsdWithError(w, 400, fmt.Sprintln("Error parsing json: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		responsdWithError(w, 400, fmt.Sprintf("Cannot create user: %v", err))
		return
	}

	respondWithJson(w, 201, dbFeedToFeed(feed))

}

func (apiCfg *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeed(r.Context())
	if err != nil {
		responsdWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
		return		
	}
	respondWithJson(w, 200, dbFeedsToFeeds(feeds))
}

