package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responsdWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error: %s", msg)
	}

	type errResponse struct {
		Error string `json:"error"`
	}
	respondWithJson(w, code, errResponse {
		Error: msg,
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)

	if err != nil {
		w.WriteHeader(500)
		log.Println("failed to marshal JSON response %v", data)
		return
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}