package db

import "backend-challenge/models"

//go:generate mockgen -source=interface.go -destination=mocks/mock_db.go -package=mocks

// Database defines the interface for database operations
type Database interface {
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id string) (*models.Product, error)
	InsertProduct(p *models.Product) error
	IsCouponValid(code string) (bool, error)
	InsertCoupon(code string, count int) error
	Close() error
}
