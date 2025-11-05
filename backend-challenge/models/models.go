package models

type Product struct {
	ID       string        `json:"id"`
	Name     string        `json:"name"`
	Category string        `json:"category"`
	Price    float64       `json:"price"`
	Image    *ProductImage `json:"image,omitempty"`
}

type ProductImage struct {
	Thumbnail string `json:"thumbnail,omitempty"`
	Mobile    string `json:"mobile,omitempty"`
	Tablet    string `json:"tablet,omitempty"`
	Desktop   string `json:"desktop,omitempty"`
}

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderReq struct {
	Items      []OrderItem `json:"items"`
	CouponCode string      `json:"couponCode,omitempty"`
}

type Order struct {
	Items      []OrderItem `json:"items"`
	CouponCode string      `json:"couponCode,omitempty"`
	ID         string      `json:"id"`
	Products   []Product   `json:"products"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}
