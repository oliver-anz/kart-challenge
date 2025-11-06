package api

import (
	"backend-challenge/db/mocks"
	"backend-challenge/models"
	"backend-challenge/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestRouter(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		headers        map[string]string
		mockSetup      func(*mocks.MockDatabase)
		expectedStatus int
	}{
		{
			name:   "GET /api/product",
			method: "GET",
			path:   "/api/product",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetAllProducts(gomock.Any()).Return([]models.Product{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /api/product - wrong method",
			method:         "POST",
			path:           "/api/product",
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "GET /api/product/:id",
			method: "GET",
			path:   "/api/product/123",
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "123").Return(&models.Product{ID: "123"}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /api/product/:id - wrong method",
			method:         "POST",
			path:           "/api/product/123",
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "GET /api/product/ - trailing slash only",
			method:         "GET",
			path:           "/api/product/",
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "POST /api/order",
			method: "POST",
			path:   "/api/order",
			body: models.OrderReq{
				Items: []models.OrderItem{{ProductID: "1", Quantity: 1}},
			},
			headers: map[string]string{"api_key": "apitest"},
			mockSetup: func(m *mocks.MockDatabase) {
				m.EXPECT().GetProductByID(gomock.Any(), "1").Return(&models.Product{ID: "1", Name: "Test", Price: 10}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /api/order - wrong method",
			method:         "GET",
			path:           "/api/order",
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "POST /api/order - missing auth",
			method: "POST",
			path:   "/api/order",
			body: models.OrderReq{
				Items: []models.OrderItem{{ProductID: "1", Quantity: 1}},
			},
			mockSetup:      func(m *mocks.MockDatabase) {},
			expectedStatus: http.StatusUnauthorized,
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
			router := handler.SetupRoutes()

			var reqBody []byte
			if tt.body != nil {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			for k, v := range tt.headers {
				req.Header.Set(k, v)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}
