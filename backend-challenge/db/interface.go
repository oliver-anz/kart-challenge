package db

import (
	"backend-challenge/models"
	"context"
)

//go:generate mockgen -source=interface.go -destination=mocks/mock_db.go -package=mocks

// Database defines the interface for database operations
type Database interface {
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	GetProductByID(ctx context.Context, id string) (*models.Product, error)
	IsCouponValid(ctx context.Context, code string) (bool, error)
	Close() error
}
