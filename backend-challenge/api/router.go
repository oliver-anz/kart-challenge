package api

import (
	"net/http"
	"strings"
)

func (h *Handler) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/product", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.ListProducts(w, r)
	})

	mux.HandleFunc("/api/product/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if strings.HasPrefix(r.URL.Path, "/api/product/") && r.URL.Path != "/api/product/" {
			h.GetProduct(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/api/order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		AuthMiddleware(h.PlaceOrder)(w, r)
	})

	return mux
}
