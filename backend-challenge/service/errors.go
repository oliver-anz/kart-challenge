package service

import "errors"

var (
	// ErrInvalidCoupon is returned when a coupon code is invalid
	ErrInvalidCoupon = errors.New("invalid coupon code")

	// ErrProductNotFound is returned when a product does not exist
	ErrProductNotFound = errors.New("product not found")
)
