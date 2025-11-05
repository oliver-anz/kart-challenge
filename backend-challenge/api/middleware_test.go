package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuthMiddleware_ValidKey(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})

	wrapped := AuthMiddleware(handler)

	req := httptest.NewRequest("POST", "/api/order", nil)
	req.Header.Set("api_key", "apitest")
	w := httptest.NewRecorder()

	wrapped(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "success" {
		t.Errorf("Expected 'success', got %s", w.Body.String())
	}
}

func TestAuthMiddleware_MissingKey(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	wrapped := AuthMiddleware(handler)

	req := httptest.NewRequest("POST", "/api/order", nil)
	w := httptest.NewRecorder()

	wrapped(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidKey(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	wrapped := AuthMiddleware(handler)

	req := httptest.NewRequest("POST", "/api/order", nil)
	req.Header.Set("api_key", "wrongkey")
	w := httptest.NewRecorder()

	wrapped(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_CustomEnvKey(t *testing.T) {
	// Set custom API key via environment
	os.Setenv("API_KEY", "customkey")
	defer os.Unsetenv("API_KEY")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	wrapped := AuthMiddleware(handler)

	req := httptest.NewRequest("POST", "/api/order", nil)
	req.Header.Set("api_key", "customkey")
	w := httptest.NewRecorder()

	wrapped(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Test that old key no longer works
	req2 := httptest.NewRequest("POST", "/api/order", nil)
	req2.Header.Set("api_key", "apitest")
	w2 := httptest.NewRecorder()

	wrapped(w2, req2)

	if w2.Code != http.StatusUnauthorized {
		t.Errorf("Expected old key to be rejected, got status %d", w2.Code)
	}
}
