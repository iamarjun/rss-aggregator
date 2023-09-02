package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/iamarju/rss-aggregator/internal/auth"
	"github.com/iamarju/rss-aggregator/internal/database"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)

	if err != nil {
		responsdWithError(w, 400, fmt.Sprintln("Error parsing json: %v", err))
		return
	}

	usr, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		responsdWithError(w, 400, fmt.Sprintf("Cannot create user: %v", err))
		return
	}

	respondWithJson(w, 201, dbUserToUser(usr))

}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)

	if err != nil {
		responsdWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
		return
	}

	usr, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)

	if err != nil {
		responsdWithError(w, 400, fmt.Sprintf("Could not get the user %v", err))
		return
	}

	respondWithJson(w, 200, dbUserToUser(usr))
}
