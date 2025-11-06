package api

import (
	"backend-challenge/db/mocks"
	"backend-challenge/models"
	"backend-challenge/service"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestListProducts(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*mocks.MockDatabase)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetAllProducts(gomock.Any(), 0, 0).Return([]models.Product{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "with limit > 100 capped to 100",
			queryParams: "?limit=200",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetAllProducts(gomock.Any(), 100, 0).Return([]models.Product{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "with offset",
			queryParams: "?offset=10",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetAllProducts(gomock.Any(), 0, 10).Return([]models.Product{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "with limit and offset",
			queryParams: "?limit=20&offset=5",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetAllProducts(gomock.Any(), 20, 5).Return([]models.Product{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "database error",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetAllProducts(gomock.Any(), 0, 0).Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mocks.NewMockDatabase(ctrl)
			tt.mockSetup(mockDB)
			svc := service.New(mockDB)
			handler := NewHandler(svc)

			url := "/api/product" + tt.queryParams
			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			handler.ListProducts(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mockSetup      func(*mocks.MockDatabase)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name:      "success",
			productID: "1",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "1").Return(&models.Product{
					ID:       "1",
					Name:     "Test Product",
					Category: "Test",
					Price:    9.99,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var product models.Product
				if err := json.NewDecoder(w.Body).Decode(&product); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if product.Name != "Test Product" {
					t.Errorf("Expected 'Test Product', got %s", product.Name)
				}
			},
		},
		{
			name:           "empty product ID",
			productID:      "",
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "not found",
			productID: "999",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "999").Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "database error",
			productID: "1",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "1").Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mocks.NewMockDatabase(ctrl)
			tt.mockSetup(mockDB)
			svc := service.New(mockDB)
			handler := NewHandler(svc)

			req := httptest.NewRequest("GET", "/api/product/"+tt.productID, nil)
			w := httptest.NewRecorder()

			handler.GetProduct(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestPlaceOrder(t *testing.T) {
	tests := []struct {
		name           string
		orderReq       interface{}
		mockSetup      func(*mocks.MockDatabase)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			orderReq: models.OrderReq{
				Items: []models.OrderItem{{ProductID: "1", Quantity: 2}},
			},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "1").Return(&models.Product{
					ID:       "1",
					Name:     "Waffle",
					Category: "Breakfast",
					Price:    6.5,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
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
			},
		},
		{
			name: "valid coupon",
			orderReq: models.OrderReq{
				Items:      []models.OrderItem{{ProductID: "1", Quantity: 1}},
				CouponCode: "HAPPYHRS",
			},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().IsCouponValid(gomock.Any(), "HAPPYHRS").Return(true, nil)
				m.EXPECT().GetProductByID(gomock.Any(), "1").Return(&models.Product{
					ID: "1", Name: "Waffle", Category: "Breakfast", Price: 6.5,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var order models.Order
				json.NewDecoder(w.Body).Decode(&order)
				if order.CouponCode != "HAPPYHRS" {
					t.Errorf("Expected coupon HAPPYHRS, got %s", order.CouponCode)
				}
			},
		},
		{
			name: "invalid coupon",
			orderReq: models.OrderReq{
				Items:      []models.OrderItem{{ProductID: "1", Quantity: 1}},
				CouponCode: "INVALID",
			},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().IsCouponValid(gomock.Any(), "INVALID").Return(false, nil)
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "product not found",
			orderReq: models.OrderReq{
				Items: []models.OrderItem{{ProductID: "999", Quantity: 1}},
			},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "999").Return(nil, nil)
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "empty items",
			orderReq:       `{"items":[]}`,
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "zero quantity",
			orderReq:       `{"items":[{"productId":"1","quantity":0}]}`,
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "negative quantity",
			orderReq:       `{"items":[{"productId":"1","quantity":-1}]}`,
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing product ID",
			orderReq:       `{"items":[{"quantity":1}]}`,
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			orderReq:       `{invalid}`,
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mocks.NewMockDatabase(ctrl)
			tt.mockSetup(mockDB)
			svc := service.New(mockDB)
			handler := NewHandler(svc)

			var body []byte
			switch v := tt.orderReq.(type) {
			case string:
				body = []byte(v)
			case models.OrderReq:
				body, _ = json.Marshal(v)
			}

			req := httptest.NewRequest("POST", "/api/order", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.PlaceOrder(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockDatabase)
		expectedStatus int
		checkResponse  func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "healthy",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetAllProducts(gomock.Any(), 1, 0).Return([]models.Product{}, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]string
				json.NewDecoder(w.Body).Decode(&resp)
				if resp["status"] != "healthy" {
					t.Errorf("Expected healthy status, got %s", resp["status"])
				}
			},
		},
		{
			name: "unhealthy",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetAllProducts(gomock.Any(), 1, 0).Return(nil, errors.New("db error"))
			},
			expectedStatus: http.StatusServiceUnavailable,
			checkResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				var resp map[string]string
				json.NewDecoder(w.Body).Decode(&resp)
				if resp["status"] != "unhealthy" {
					t.Errorf("Expected unhealthy status, got %s", resp["status"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mocks.NewMockDatabase(ctrl)
			tt.mockSetup(mockDB)
			svc := service.New(mockDB)
			handler := NewHandler(svc)

			req := httptest.NewRequest("GET", "/health", nil)
			w := httptest.NewRecorder()

			handler.HealthCheck(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, w)
			}
		})
	}
}
