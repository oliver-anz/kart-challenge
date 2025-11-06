package api

import (
	"backend-challenge/models"
	"encoding/json"
	"net/http"
	"os"
)

func getValidAPIKey() string {
	// Allow API key to be set via environment variable, default to "apitest"
	if key := os.Getenv("API_KEY"); key != "" {
		return key
	}
	return "apitest"
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	validAPIKey := getValidAPIKey()
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("api_key")
		if apiKey != validAPIKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Type:    "error",
				Message: "Invalid or missing API key",
			})
			return
		}
		next(w, r)
	}
}
