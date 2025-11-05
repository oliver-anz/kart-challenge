package api

import (
	"backend-challenge/db"
	"backend-challenge/models"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type Handler struct {
	DB *db.DB
}

func NewHandler(database *db.DB) *Handler {
	return &Handler{DB: database}
}

func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.DB.GetAllProducts()
	if err != nil {
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

	product, err := h.DB.GetProductByID(productID)
	if err != nil {
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

	products := make([]models.Product, 0)
	for _, item := range req.Items {
		product, err := h.DB.GetProductByID(item.ProductID)
		if err != nil {
			h.sendError(w, http.StatusInternalServerError, "error", "Failed to fetch product")
			return
		}
		if product == nil {
			h.sendError(w, http.StatusUnprocessableEntity, "error", "Product not found: "+item.ProductID)
			return
		}
		products = append(products, *product)
	}

	if req.CouponCode != "" {
		valid, err := h.DB.IsCouponValid(req.CouponCode)
		if err != nil {
			h.sendError(w, http.StatusInternalServerError, "error", "Failed to validate coupon")
			return
		}
		if !valid {
			h.sendError(w, http.StatusUnprocessableEntity, "error", "Invalid coupon code")
			return
		}
	}

	order := models.Order{
		ID:         uuid.New().String(),
		Items:      req.Items,
		Products:   products,
		CouponCode: req.CouponCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
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
