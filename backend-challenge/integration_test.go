package main

import (
	"backend-challenge/api"
	"backend-challenge/db"
	"backend-challenge/models"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

// Integration tests that test the full stack including routing

func setupIntegrationTest(t *testing.T) (*httptest.Server, func()) {
	dbPath := "test_integration.db"

	// Initialize database from init.sql
	cmd := exec.Command("sqlite3", dbPath, ".read init.sql")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	database, err := db.New(dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// init.sql already contains all necessary test data (9 products and 8 coupons)

	handler := api.NewHandler(database)
	router := handler.SetupRoutes()
	server := httptest.NewServer(router)

	cleanup := func() {
		server.Close()
		database.Close()
		os.Remove(dbPath)
	}

	return server, cleanup
}

func TestIntegration_FullOrderFlow(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// 1. List all products
	resp, err := http.Get(server.URL + "/api/product")
	if err != nil {
		t.Fatalf("Failed to get products: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var products []models.Product
	json.NewDecoder(resp.Body).Decode(&products)
	if len(products) != 9 {
		t.Errorf("Expected 9 products, got %d", len(products))
	}

	// 2. Get specific product
	resp, err = http.Get(server.URL + "/api/product/1")
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}
	defer resp.Body.Close()

	var product models.Product
	json.NewDecoder(resp.Body).Decode(&product)
	if product.Name != "Waffle with Berries" {
		t.Errorf("Expected 'Waffle with Berries', got %s", product.Name)
	}

	// 3. Place order with coupon
	orderReq := models.OrderReq{
		Items:      []models.OrderItem{{ProductID: "1", Quantity: 2}},
		CouponCode: "HAPPYHRS",
	}
	body, _ := json.Marshal(orderReq)

	req, _ := http.NewRequest("POST", server.URL+"/api/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api_key", "apitest")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to place order: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var order models.Order
	json.NewDecoder(resp.Body).Decode(&order)

	if order.ID == "" {
		t.Error("Expected order ID")
	}

	if order.CouponCode != "HAPPYHRS" {
		t.Errorf("Expected coupon HAPPYHRS, got %s", order.CouponCode)
	}

	if len(order.Products) != 1 {
		t.Errorf("Expected 1 product, got %d", len(order.Products))
	}
}

func TestIntegration_AuthenticationRequired(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	orderReq := models.OrderReq{
		Items: []models.OrderItem{{ProductID: "1", Quantity: 1}},
	}
	body, _ := json.Marshal(orderReq)

	// Request without API key
	req, _ := http.NewRequest("POST", server.URL+"/api/order", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", resp.StatusCode)
	}
}

func TestIntegration_AllValidCoupons(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	validCoupons := []string{"HAPPYHRS", "FIFTYOFF"}

	for _, coupon := range validCoupons {
		t.Run(coupon, func(t *testing.T) {
			orderReq := models.OrderReq{
				Items:      []models.OrderItem{{ProductID: "1", Quantity: 1}},
				CouponCode: coupon,
			}
			body, _ := json.Marshal(orderReq)

			req, _ := http.NewRequest("POST", server.URL+"/api/order", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("api_key", "apitest")

			resp, _ := http.DefaultClient.Do(req)
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Coupon %s should be valid, got status %d", coupon, resp.StatusCode)
			}
		})
	}
}
