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
			// Save original env value to avoid test pollution
			originalKey := os.Getenv("API_KEY")
			if tt.envKey != "" {
				os.Setenv("API_KEY", tt.envKey)
			} else {
				os.Unsetenv("API_KEY")
			}
			defer func() {
				if originalKey != "" {
					os.Setenv("API_KEY", originalKey)
				} else {
					os.Unsetenv("API_KEY")
				}
			}()

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

func TestRequestIDMiddleware(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("generates ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Header().Get("X-Request-ID") == "" {
			t.Error("X-Request-ID header not set")
		}
	})

	t.Run("uses existing ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-Request-ID", "existing-id")
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Header().Get("X-Request-ID") != "existing-id" {
			t.Errorf("Expected existing-id, got %s", w.Header().Get("X-Request-ID"))
		}
	})
}

func TestCORSMiddleware(t *testing.T) {
	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("OPTIONS request", func(t *testing.T) {
		req := httptest.NewRequest("OPTIONS", "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected 204, got %d", w.Code)
		}

		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("CORS headers not set")
		}
	})

	t.Run("regular request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		if w.Header().Get("Access-Control-Allow-Origin") == "" {
			t.Error("CORS headers not set on regular request")
		}
	})
}

func TestMaxBodySizeMiddleware(t *testing.T) {
	handler := MaxBodySizeMiddleware(10)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("within limit", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})
}
