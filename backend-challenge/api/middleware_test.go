package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name              string
		apiKey            string
		envKey            string
		expectedStatus    int
		shouldCallHandler bool
	}{
		{
			name:              "valid key",
			apiKey:            "apitest",
			expectedStatus:    http.StatusOK,
			shouldCallHandler: true,
		},
		{
			name:              "missing key",
			apiKey:            "",
			expectedStatus:    http.StatusUnauthorized,
			shouldCallHandler: false,
		},
		{
			name:              "invalid key",
			apiKey:            "wrongkey",
			expectedStatus:    http.StatusUnauthorized,
			shouldCallHandler: false,
		},
		{
			name:              "custom env key - valid",
			apiKey:            "customkey",
			envKey:            "customkey",
			expectedStatus:    http.StatusOK,
			shouldCallHandler: true,
		},
		{
			name:              "custom env key - old key rejected",
			apiKey:            "apitest",
			envKey:            "customkey",
			expectedStatus:    http.StatusUnauthorized,
			shouldCallHandler: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envKey != "" {
				os.Setenv("API_KEY", tt.envKey)
				defer os.Unsetenv("API_KEY")
			}

			handlerCalled := false
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			})

			wrapped := AuthMiddleware(handler)

			req := httptest.NewRequest("POST", "/api/order", nil)
			if tt.apiKey != "" {
				req.Header.Set("api_key", tt.apiKey)
			}
			w := httptest.NewRecorder()

			wrapped(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if handlerCalled != tt.shouldCallHandler {
				t.Errorf("Expected handler called = %v, got %v", tt.shouldCallHandler, handlerCalled)
			}

			if tt.shouldCallHandler && w.Body.String() != "success" {
				t.Errorf("Expected 'success', got %s", w.Body.String())
			}
		})
	}
}
