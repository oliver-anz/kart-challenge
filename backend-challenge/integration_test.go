package main

import (
	"backend-challenge/api"
	"backend-challenge/db"
	"backend-challenge/models"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupIntegrationTest(t *testing.T) (*httptest.Server, func()) {
	database, err := db.New("data/store.db")
	require.NoError(t, err)

	handler := api.NewHandler(database)
	router := handler.SetupRoutes()
	server := httptest.NewServer(router)

	cleanup := func() {
		server.Close()
		database.Close()
	}

	return server, cleanup
}

func TestIntegration_GetProducts(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	resp, err := http.Get(server.URL + "/api/product")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var products []models.Product
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&products))
	assert.Len(t, products, 9)

	// Verify specific products match live API
	expectedProducts := map[string]struct {
		name      string
		category  string
		price     float64
		thumbnail string
		mobile    string
		tablet    string
		desktop   string
	}{
		"1": {"Waffle with Berries", "Waffle", 6.5,
			"https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg"},
		"2": {"Vanilla Bean Crème Brûlée", "Crème Brûlée", 7,
			"https://orderfoodonline.deno.dev/public/images/image-creme-brulee-thumbnail.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-creme-brulee-mobile.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-creme-brulee-tablet.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-creme-brulee-desktop.jpg"},
		"3": {"Macaron Mix of Five", "Macaron", 8,
			"https://orderfoodonline.deno.dev/public/images/image-macaron-thumbnail.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-macaron-mobile.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-macaron-tablet.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-macaron-desktop.jpg"},
		"4": {"Classic Tiramisu", "Tiramisu", 5.5,
			"https://orderfoodonline.deno.dev/public/images/image-tiramisu-thumbnail.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-tiramisu-mobile.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-tiramisu-tablet.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-tiramisu-desktop.jpg"},
		"5": {"Pistachio Baklava", "Baklava", 4,
			"https://orderfoodonline.deno.dev/public/images/image-baklava-thumbnail.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-baklava-mobile.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-baklava-tablet.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-baklava-desktop.jpg"},
		"6": {"Lemon Meringue Pie", "Pie", 5,
			"https://orderfoodonline.deno.dev/public/images/image-meringue-thumbnail.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-meringue-mobile.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-meringue-tablet.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-meringue-desktop.jpg"},
		"7": {"Red Velvet Cake", "Cake", 4.5,
			"https://orderfoodonline.deno.dev/public/images/image-cake-thumbnail.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-cake-mobile.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-cake-tablet.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-cake-desktop.jpg"},
		"8": {"Salted Caramel Brownie", "Brownie", 4.5,
			"https://orderfoodonline.deno.dev/public/images/image-brownie-thumbnail.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-brownie-mobile.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-brownie-tablet.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-brownie-desktop.jpg"},
		"9": {"Vanilla Panna Cotta", "Panna Cotta", 6.5,
			"https://orderfoodonline.deno.dev/public/images/image-panna-cotta-thumbnail.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-panna-cotta-mobile.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-panna-cotta-tablet.jpg",
			"https://orderfoodonline.deno.dev/public/images/image-panna-cotta-desktop.jpg"},
	}

	for _, p := range products {
		expected, ok := expectedProducts[p.ID]
		assert.True(t, ok, "Unexpected product ID: %s", p.ID)
		if !ok {
			continue
		}
		assert.Equal(t, expected.name, p.Name)
		assert.Equal(t, expected.category, p.Category)
		assert.Equal(t, expected.price, p.Price)
		assert.NotNil(t, p.Image)
		if p.Image != nil {
			assert.Equal(t, expected.thumbnail, p.Image.Thumbnail)
			assert.Equal(t, expected.mobile, p.Image.Mobile)
			assert.Equal(t, expected.tablet, p.Image.Tablet)
			assert.Equal(t, expected.desktop, p.Image.Desktop)
		}
	}
}

func TestIntegration_GetProductByID(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	tests := []struct {
		name                   string
		productID              string
		expectedStatus         int
		expectedName           string
		expectedCat            string
		expectedPrice          float64
		expectedThumbnailImage string
		expectedMobileImage    string
		expectedTabletImage    string
		expectedDesktopImage   string
	}{
		{
			name:                   "valid product ID 1",
			productID:              "1",
			expectedStatus:         http.StatusOK,
			expectedName:           "Waffle with Berries",
			expectedCat:            "Waffle",
			expectedPrice:          6.5,
			expectedThumbnailImage: "https://orderfoodonline.deno.dev/public/images/image-waffle-thumbnail.jpg",
			expectedMobileImage:    "https://orderfoodonline.deno.dev/public/images/image-waffle-mobile.jpg",
			expectedTabletImage:    "https://orderfoodonline.deno.dev/public/images/image-waffle-tablet.jpg",
			expectedDesktopImage:   "https://orderfoodonline.deno.dev/public/images/image-waffle-desktop.jpg",
		},
		{
			name:           "non-existent product",
			productID:      "999",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Get(server.URL + "/api/product/" + tt.productID)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				var product models.Product
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&product))
				assert.Equal(t, tt.productID, product.ID)
				assert.Equal(t, tt.expectedName, product.Name)
				assert.Equal(t, tt.expectedCat, product.Category)
				assert.Equal(t, tt.expectedPrice, product.Price)
				assert.NotNil(t, product.Image)
				if product.Image != nil {
					assert.Equal(t, tt.expectedThumbnailImage, product.Image.Thumbnail)
					assert.Equal(t, tt.expectedMobileImage, product.Image.Mobile)
					assert.Equal(t, tt.expectedTabletImage, product.Image.Tablet)
					assert.Equal(t, tt.expectedDesktopImage, product.Image.Desktop)
				}
			}
		})
	}
}

func TestIntegration_PlaceOrder(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	tests := []struct {
		name           string
		body           models.OrderReq
		apiKey         string
		expectedStatus int
		checkOrder     func(*testing.T, *models.Order)
	}{
		{
			name:           "successful order with coupon",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 2}}, CouponCode: "HAPPYHRS"},
			apiKey:         "apitest",
			expectedStatus: http.StatusOK,
			checkOrder: func(t *testing.T, o *models.Order) {
				assert.NotEmpty(t, o.ID)
				assert.Equal(t, "HAPPYHRS", o.CouponCode)
				assert.Len(t, o.Items, 1)
				assert.Len(t, o.Products, 1)
				assert.Equal(t, "1", o.Items[0].ProductID)
				assert.Equal(t, 2, o.Items[0].Quantity)
			},
		},
		{
			name:           "valid coupon BIRTHDAY",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "BIRTHDAY"},
			apiKey:         "apitest",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid coupon BUYGETON",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "BUYGETON"},
			apiKey:         "apitest",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid coupon FIFTYOFF",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "FIFTYOFF"},
			apiKey:         "apitest",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid coupon FREEZAAA",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "FREEZAAA"},
			apiKey:         "apitest",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid coupon GNULINUX",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "GNULINUX"},
			apiKey:         "apitest",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid coupon OVER9000",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "OVER9000"},
			apiKey:         "apitest",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "valid coupon SIXTYOFF",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "SIXTYOFF"},
			apiKey:         "apitest",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid coupon",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}, CouponCode: "INVALID99"},
			apiKey:         "apitest",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "missing API key",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}},
			apiKey:         "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid API key",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 1}}},
			apiKey:         "wrongkey",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "empty items",
			body:           models.OrderReq{Items: []models.OrderItem{}},
			apiKey:         "apitest",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "zero quantity",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: 0}}},
			apiKey:         "apitest",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "negative quantity",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "1", Quantity: -5}}},
			apiKey:         "apitest",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "non-existent product",
			body:           models.OrderReq{Items: []models.OrderItem{{ProductID: "999", Quantity: 1}}},
			apiKey:         "apitest",
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, err := json.Marshal(tt.body)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", server.URL+"/api/order", bytes.NewReader(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			if tt.apiKey != "" {
				req.Header.Set("api_key", tt.apiKey)
			}

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			if !assert.Equal(t, tt.expectedStatus, resp.StatusCode) {
				body, _ := io.ReadAll(resp.Body)
				t.Logf("Response body: %s", string(body))
			}

			if tt.checkOrder != nil && resp.StatusCode == http.StatusOK {
				var order models.Order
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&order))
				tt.checkOrder(t, &order)
			}
		})
	}
}

func TestIntegration_OpenAPISpec(t *testing.T) {
	server, cleanup := setupIntegrationTest(t)
	defer cleanup()

	resp, err := http.Get(server.URL + "/public/openapi.yaml")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/yaml", resp.Header.Get("Content-Type"))

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.NotEmpty(t, body)
}
