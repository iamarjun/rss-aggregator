package main

import (
	"fmt"
	"net/http"

	"github.com/iamarju/rss-aggregator/internal/auth"
	"github.com/iamarju/rss-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)

		if err != nil {
			responsdWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		usr, err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)

		if err != nil {
			responsdWithError(w, 400, fmt.Sprintf("Could not get the user %v", err))
			return
		}

		handler(w, r, usr)
	}
}
