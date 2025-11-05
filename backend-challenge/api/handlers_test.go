package api

import (
	"backend-challenge/db"
	"backend-challenge/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupTestHandler(t *testing.T) (*Handler, func()) {
	dbPath := "test_handlers.db"
	database, err := db.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Insert test products
	products := []*models.Product{
		{ID: "1", Name: "Product 1", Category: "Cat1", Price: 10.0, Image: &models.ProductImage{Thumbnail: "t1.jpg"}},
		{ID: "2", Name: "Product 2", Category: "Cat2", Price: 20.0, Image: &models.ProductImage{Thumbnail: "t2.jpg"}},
	}
	for _, p := range products {
		database.InsertProduct(p)
	}

	// Insert test coupons
	database.InsertCoupon("VALIDCODE", 2)

	handler := NewHandler(database)

	cleanup := func() {
		database.Close()
		os.Remove(dbPath)
	}

	return handler, cleanup
}

func TestListProducts(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/product", nil)
	w := httptest.NewRecorder()

	handler.ListProducts(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var products []models.Product
	if err := json.NewDecoder(w.Body).Decode(&products); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %d", len(products))
	}
}

func TestGetProduct_Success(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/product/1", nil)
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var product models.Product
	if err := json.NewDecoder(w.Body).Decode(&product); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if product.ID != "1" || product.Name != "Product 1" {
		t.Errorf("Unexpected product: %+v", product)
	}
}

func TestGetProduct_NotFound(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/product/999", nil)
	w := httptest.NewRecorder()

	handler.GetProduct(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestPlaceOrder_Success(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	orderReq := models.OrderReq{
		Items: []models.OrderItem{
			{ProductID: "1", Quantity: 2},
		},
	}

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest("POST", "/api/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PlaceOrder(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var order models.Order
	if err := json.NewDecoder(w.Body).Decode(&order); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if order.ID == "" {
		t.Error("Expected order ID to be generated")
	}

	if len(order.Products) != 1 {
		t.Errorf("Expected 1 product in order, got %d", len(order.Products))
	}
}

func TestPlaceOrder_WithValidCoupon(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	orderReq := models.OrderReq{
		Items:      []models.OrderItem{{ProductID: "1", Quantity: 1}},
		CouponCode: "VALIDCODE",
	}

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest("POST", "/api/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PlaceOrder(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var order models.Order
	json.NewDecoder(w.Body).Decode(&order)

	if order.CouponCode != "VALIDCODE" {
		t.Errorf("Expected coupon code VALIDCODE, got %s", order.CouponCode)
	}
}

func TestPlaceOrder_InvalidCoupon(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	orderReq := models.OrderReq{
		Items:      []models.OrderItem{{ProductID: "1", Quantity: 1}},
		CouponCode: "INVALID",
	}

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest("POST", "/api/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PlaceOrder(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", w.Code)
	}
}

func TestPlaceOrder_ProductNotFound(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	orderReq := models.OrderReq{
		Items: []models.OrderItem{{ProductID: "999", Quantity: 1}},
	}

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest("POST", "/api/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.PlaceOrder(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", w.Code)
	}
}

func TestPlaceOrder_InvalidInput(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	tests := []struct {
		name    string
		reqBody string
		status  int
	}{
		{"empty items", `{"items":[]}`, http.StatusBadRequest},
		{"zero quantity", `{"items":[{"productId":"1","quantity":0}]}`, http.StatusBadRequest},
		{"negative quantity", `{"items":[{"productId":"1","quantity":-1}]}`, http.StatusBadRequest},
		{"missing product ID", `{"items":[{"quantity":1}]}`, http.StatusBadRequest},
		{"invalid JSON", `{invalid}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/order", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.PlaceOrder(w, req)

			if w.Code != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, w.Code)
			}
		})
	}
}
