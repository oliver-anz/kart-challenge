package api

import (
	"backend-challenge/models"
	"backend-challenge/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Parse pagination query parameters
	limit := 0 // 0 means no limit
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > 100 {
				limit = 100 // Max limit of 100
			}
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	products, err := h.svc.GetAllProducts(r.Context(), limit, offset)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		h.sendError(w, http.StatusInternalServerError, "error", "Failed to fetch products")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/product/")
	productID := path

	if productID == "" {
		h.sendError(w, http.StatusBadRequest, "error", "Invalid product ID")
		return
	}

	product, err := h.svc.GetProductByID(r.Context(), productID)
	if err != nil {
		log.Printf("Error fetching product %s: %v", productID, err)
		h.sendError(w, http.StatusInternalServerError, "error", "Failed to fetch product")
		return
	}

	if product == nil {
		h.sendError(w, http.StatusNotFound, "error", "Product not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req models.OrderReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "error", "Invalid input")
		return
	}

	if len(req.Items) == 0 {
		h.sendError(w, http.StatusBadRequest, "error", "Order must contain at least one item")
		return
	}

	for _, item := range req.Items {
		if item.ProductID == "" {
			h.sendError(w, http.StatusBadRequest, "error", "Product ID is required")
			return
		}
		if item.Quantity <= 0 {
			h.sendError(w, http.StatusBadRequest, "error", "Quantity must be positive")
			return
		}
	}

	order, err := h.svc.PlaceOrder(r.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCoupon) {
			h.sendError(w, http.StatusUnprocessableEntity, "error", "Invalid coupon code")
			return
		}
		if errors.Is(err, service.ErrProductNotFound) {
			h.sendError(w, http.StatusUnprocessableEntity, "error", "Product not found")
			return
		}
		log.Printf("Error placing order: %v", err)
		h.sendError(w, http.StatusInternalServerError, "error", "Failed to place order")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Check database connectivity
	ctx := r.Context()
	// Just fetch one product to test connectivity
	_, err := h.svc.GetAllProducts(ctx, 1, 0)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"reason": "database unavailable",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func (h *Handler) sendError(w http.ResponseWriter, statusCode int, errType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Code:    statusCode,
		Type:    errType,
		Message: message,
	})
}
