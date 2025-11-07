package service

import (
	"backend-challenge/db"
	"backend-challenge/models"
	"context"

	"github.com/google/uuid"
)

// Service provides business logic operations
type Service struct {
	db db.Database
}

// New creates a new Service
func New(database db.Database) *Service {
	return &Service{db: database}
}

// GetAllProducts retrieves all products with optional pagination
func (s *Service) GetAllProducts(ctx context.Context, limit, offset int) ([]models.Product, error) {
	return s.db.GetAllProducts(ctx, limit, offset)
}

// GetProductByID retrieves a single product by ID
func (s *Service) GetProductByID(ctx context.Context, id string) (*models.Product, error) {
	return s.db.GetProductByID(ctx, id)
}

// PlaceOrder processes an order request
func (s *Service) PlaceOrder(ctx context.Context, req models.OrderReq) (*models.Order, error) {
	// Validate coupon if provided
	if req.CouponCode != "" {
		valid, err := s.db.IsCouponValid(ctx, req.CouponCode)
		if err != nil {
			return nil, err
		}
		if !valid {
			return nil, ErrInvalidCoupon
		}
	}

	// Fetch products for order
	products := make([]models.Product, 0, len(req.Items))
	for _, item := range req.Items {
		product, err := s.db.GetProductByID(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}
		if product == nil {
			return nil, ErrProductNotFound
		}
		products = append(products, *product)
	}

	// Generate order
	// The order is not stored or will not persist anywhere for the purposes of this demo
	order := &models.Order{
		ID:         uuid.New().String(),
		Items:      req.Items,
		Products:   products,
		CouponCode: req.CouponCode,
	}

	return order, nil
}
