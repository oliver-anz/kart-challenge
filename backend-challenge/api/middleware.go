package api

import (
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
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    401,
				"type":    "error",
				"message": "Invalid or missing API key",
			})
			return
		}
		next(w, r)
	}
}
