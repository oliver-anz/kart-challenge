package db

import (
	"backend-challenge/models"
	"os"
	"testing"
)

func setupTestDB(t *testing.T) *DB {
	// Create temporary database
	dbPath := "test.db"
	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	// Cleanup function
	t.Cleanup(func() {
		db.Close()
		os.Remove(dbPath)
	})

	return db
}

func TestInsertAndGetProduct(t *testing.T) {
	db := setupTestDB(t)

	product := &models.Product{
		ID:       "1",
		Name:     "Test Product",
		Category: "Test",
		Price:    9.99,
		Image: &models.ProductImage{
			Thumbnail: "thumb.jpg",
			Mobile:    "mobile.jpg",
		},
	}

	if err := db.InsertProduct(product); err != nil {
		t.Fatalf("Failed to insert product: %v", err)
	}

	retrieved, err := db.GetProductByID("1")
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Expected product, got nil")
	}

	if retrieved.Name != product.Name || retrieved.Price != product.Price {
		t.Errorf("Product mismatch: got %+v, want %+v", retrieved, product)
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	db := setupTestDB(t)

	product, err := db.GetProductByID("999")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if product != nil {
		t.Errorf("Expected nil for non-existent product, got %+v", product)
	}
}

func TestGetAllProducts(t *testing.T) {
	db := setupTestDB(t)

	products := []*models.Product{
		{ID: "1", Name: "Product 1", Category: "Cat1", Price: 10.0},
		{ID: "2", Name: "Product 2", Category: "Cat2", Price: 20.0},
	}

	for _, p := range products {
		if err := db.InsertProduct(p); err != nil {
			t.Fatalf("Failed to insert product: %v", err)
		}
	}

	all, err := db.GetAllProducts()
	if err != nil {
		t.Fatalf("Failed to get all products: %v", err)
	}

	if len(all) != 2 {
		t.Errorf("Expected 2 products, got %d", len(all))
	}
}

func TestIsCouponValid(t *testing.T) {
	db := setupTestDB(t)

	// Insert valid coupon
	if err := db.InsertCoupon("TESTCODE", 2); err != nil {
		t.Fatalf("Failed to insert coupon: %v", err)
	}

	tests := []struct {
		name  string
		code  string
		valid bool
	}{
		{"valid coupon", "TESTCODE", true},
		{"invalid coupon", "INVALID", false},
		{"empty coupon", "", true},
		{"too short", "SHORT", false},
		{"too long", "TOOLONGCODE", false},
		{"valid length but not in DB", "NOTINDB88", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := db.IsCouponValid(tt.code)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if valid != tt.valid {
				t.Errorf("IsCouponValid(%q) = %v, want %v", tt.code, valid, tt.valid)
			}
		})
	}
}

func TestProductImage_NilHandling(t *testing.T) {
	db := setupTestDB(t)

	// Product without image
	product := &models.Product{
		ID:       "1",
		Name:     "No Image Product",
		Category: "Test",
		Price:    5.0,
	}

	if err := db.InsertProduct(product); err != nil {
		t.Fatalf("Failed to insert product: %v", err)
	}

	retrieved, err := db.GetProductByID("1")
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if retrieved.Image != nil {
		t.Errorf("Expected nil image, got %+v", retrieved.Image)
	}
}
