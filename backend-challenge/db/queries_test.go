package db

import (
	"context"
	"testing"
)

func setupTestDB(t *testing.T) *DB {
	// Use the committed database
	dbPath := "../data/store.db"

	db, err := New(dbPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Cleanup function (close but don't delete committed database)
	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestGetProductByID(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	// Test existing product from init.sql (Waffle with Berries, ID=1)
	product, err := db.GetProductByID(ctx, "1")
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if product == nil {
		t.Fatal("Expected product, got nil")
	}

	if product.Name != "Waffle with Berries" || product.Price != 6.5 {
		t.Errorf("Product mismatch: got %+v, want Name='Waffle with Berries' Price=6.5", product)
	}
}

func TestGetProductByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	product, err := db.GetProductByID(ctx, "999")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if product != nil {
		t.Errorf("Expected nil for non-existent product, got %+v", product)
	}
}

func TestGetAllProducts(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	all, err := db.GetAllProducts(ctx)
	if err != nil {
		t.Fatalf("Failed to get all products: %v", err)
	}

	// init.sql has 9 products
	if len(all) != 9 {
		t.Errorf("Expected 9 products, got %d", len(all))
	}
}

func TestIsCouponValid(t *testing.T) {
	db := setupTestDB(t)
	ctx := context.Background()

	tests := []struct {
		name  string
		code  string
		valid bool
	}{
		{"valid coupon from init.sql", "HAPPYHRS", true},
		{"another valid coupon", "FIFTYOFF", true},
		{"invalid coupon", "INVALID", false},
		{"empty coupon", "", true},
		{"too short", "SHORT", false},
		{"too long", "TOOLONGCODE", false},
		{"valid length but not in DB", "NOTINDB88", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, err := db.IsCouponValid(ctx, tt.code)
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
	ctx := context.Background()

	// Vanilla Panna Cotta (ID=9) has images in init.sql
	retrieved, err := db.GetProductByID(ctx, "9")
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if retrieved.Image == nil {
		t.Errorf("Expected image for Vanilla Panna Cotta, got nil")
	}
}
