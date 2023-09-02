package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/iamarju/rss-aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)

	if err != nil {
		responsdWithError(w, 400, fmt.Sprintln("Error parsing json %v", err))
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
		responsdWithError(w, 400, fmt.Sprintf("Cannot feed follow %v", err))
		return
	}

	respondWithJson(w, 201, dbFeedFollowToFeedFollow(feedFollow))

}

func (apiCfg *apiConfig) handlerGetFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := apiCfg.DB.GetFeedFollow(r.Context(), user.ID)
	if err != nil {
		responsdWithError(w, 400, fmt.Sprintf("Couldn't get feed follows %v", err))
		return
	}
	respondWithJson(w, 200, dbFeedFollowsToFeedFollows(feedFollows))
}

func (apiCcfg *apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowId, err := uuid.Parse(chi.URLParam(r, "feedFollowId"))

	if err != nil {
		responsdWithError(w, 400, fmt.Sprintf("could not parse feed follow id %v", err))
		return
	}

	err = apiCcfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})

	if err != nil {
		responsdWithError(w, 400, fmt.Sprintf("could not delete feed follow %v", err))
		return
	}
	respondWithJson(w, 200, struct{}{})
}
