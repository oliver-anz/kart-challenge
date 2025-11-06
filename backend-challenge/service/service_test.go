package service

import (
	"backend-challenge/db/mocks"
	"backend-challenge/models"
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestGetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockDB.EXPECT().GetAllProducts(gomock.Any(), 10, 0).Return([]models.Product{{ID: "1"}}, nil)

	svc := New(mockDB)
	products, err := svc.GetAllProducts(context.Background(), 10, 0)

	if err != nil || len(products) != 1 {
		t.Error("failed")
	}
}

func TestGetProductByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockDatabase(ctrl)
	mockDB.EXPECT().GetProductByID(gomock.Any(), "1").Return(&models.Product{ID: "1"}, nil)

	svc := New(mockDB)
	product, err := svc.GetProductByID(context.Background(), "1")

	if err != nil || product.ID != "1" {
		t.Error("failed")
	}
}

func TestPlaceOrder(t *testing.T) {
	tests := []struct {
		name      string
		req       models.OrderReq
		mockSetup func(*mocks.MockDatabase)
		wantErr   bool
		checkErr  func(error) bool
	}{
		{
			name: "success",
			req:  models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 2}}},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "1").Return(&models.Product{ID: "1", Price: 10.0}, nil)
			},
		},
		{
			name: "valid coupon",
			req:  models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "HAPPYHRS"},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().IsCouponValid(gomock.Any(), "HAPPYHRS").Return(true, nil)
				m.EXPECT().GetProductByID(gomock.Any(), "1").Return(&models.Product{ID: "1", Price: 10.0}, nil)
			},
		},
		{
			name: "invalid coupon",
			req:  models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "INVALID"},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().IsCouponValid(gomock.Any(), "INVALID").Return(false, nil)
			},
			wantErr:  true,
			checkErr: func(err error) bool { return errors.Is(err, ErrInvalidCoupon) },
		},
		{
			name: "coupon check error",
			req:  models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "TEST"},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().IsCouponValid(gomock.Any(), "TEST").Return(false, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "product not found",
			req:  models.OrderReq{Items: []models.OrderItem{{ProductID: "999", Quantity: 1}}},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "999").Return(nil, nil)
			},
			wantErr:  true,
			checkErr: func(err error) bool { return errors.Is(err, ErrProductNotFound) },
		},
		{
			name: "product fetch error",
			req:  models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "1").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "multiple items",
			req:  models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 2}, {ProductID: "2", Quantity: 1}}},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "1").Return(&models.Product{ID: "1", Price: 10.0}, nil)
				m.EXPECT().GetProductByID(gomock.Any(), "2").Return(&models.Product{ID: "2", Price: 20.0}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mocks.NewMockDatabase(ctrl)
			tt.mockSetup(mockDB)

			svc := New(mockDB)
			order, err := svc.PlaceOrder(context.Background(), tt.req)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				if tt.checkErr != nil && !tt.checkErr(err) {
					t.Errorf("wrong error: %v", err)
				}
				return
			}

			if err != nil || order.ID == "" || len(order.Products) != len(tt.req.Items) {
				t.Error("failed")
			}
		})
	}
}
