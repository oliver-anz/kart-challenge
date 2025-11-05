# Backend Challenge - Food Ordering API

A Go-based REST API server implementing the OpenAPI 3.1 specification for a food ordering system with coupon validation.

## Features

- ✅ Full OpenAPI 3.1 compliance
- ✅ Product listing and retrieval
- ✅ Order placement with validation
- ✅ Coupon code validation (8-10 characters, appearing in 2+ coupon files)
- ✅ API key authentication
- ✅ SQLite database for products and coupons
- ✅ Clean, idiomatic Go code

## Quick Start

### Prerequisites

- Go 1.21 or higher
- SQLite3 command-line tool

### 1. Setup Database

```bash
# Create database with all products and valid coupons
sqlite3 data/store.db < init.sql

# This will create:
# - 9 products from the demo API
# - 8 valid coupons: BIRTHDAY, BUYGETON, FIFTYOFF, FREEZAAA, GNULINUX, HAPPYHRS, OVER9000, SIXTYOFF
```

### 2. Build and Run

```bash
# Build the server
go build -o backend-challenge .

# Run the server
./backend-challenge

# Server will start on http://localhost:8080
```

### 3. Test the API

```bash
# List all products
curl http://localhost:8080/api/product

# Get a specific product
curl http://localhost:8080/api/product/1

# Place an order (requires api_key header)
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: apitest" \
  -d '{"items":[{"productId":"1","quantity":2}]}'

# Place an order with valid coupon
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: apitest" \
  -d '{"items":[{"productId":"1","quantity":2}],"couponCode":"HAPPYHRS"}'
```

## API Endpoints

### GET /api/product

List all available products.

**Response:** Array of Product objects

### GET /api/product/{productId}

Get details of a specific product.

**Parameters:**
- `productId` (path, required): Product ID

**Response:** Product object or 404 if not found

### POST /api/order

Place a new order.

**Headers:**
- `api_key`: apitest (required)

**Request Body:**
```json
{
  "items": [
    {"productId": "1", "quantity": 2}
  ],
  "couponCode": "HAPPYHRS" // optional
}
```

**Response:** Order object with UUID and product details

**Error Codes:**
- 400: Invalid input
- 401: Invalid or missing API key
- 422: Validation error (invalid coupon, unknown product)

## Coupon Validation

Coupons are validated based on:
1. Length: 8-10 characters
2. Must appear in at least 2 of the 3 coupon base files

**Test Coupons:**
- ✅ HAPPYHRS (valid)
- ✅ FIFTYOFF (valid)
- ❌ SUPER100 (invalid)

## Full Coupon Preprocessing (Optional)

The database already contains all 8 valid coupons. To regenerate them from scratch:

```bash
# 1. Download coupon files to coupon/ directory
cd coupon
curl -O https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase1.gz
curl -O https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase2.gz
curl -O https://orderfoodonline-files.s3.ap-southeast-2.amazonaws.com/couponbase3.gz

# 2. Run preprocessing (WARNING: Takes ~3 hours on M3 Pro, requires 15-18GB RAM)
python3 process_coupons.py > valid_coupons.txt
```

See `coupon/README.md` for details on the preprocessing algorithm.

## Project Structure

```
backend-challenge/
├── main.go                 # Entry point and server setup
├── api/
│   ├── handlers.go         # HTTP request handlers
│   ├── middleware.go       # Authentication middleware
│   └── router.go           # Route definitions
├── db/
│   ├── db.go              # Database connection and schema
│   └── queries.go         # Database queries
├── models/
│   └── models.go          # Data structures
├── coupon/
│   ├── process_coupons.py  # Coupon preprocessing script
│   ├── valid_coupons.txt   # 8 valid coupon codes
│   └── README.md           # Coupon processing documentation
├── data/
│   └── store.db           # SQLite database (pre-populated)
├── init.sql               # Database initialization SQL
└── README_SETUP.md        # This file
```

## Configuration

Command-line flags:

```bash
./backend-challenge -port 8080 -db data/store.db
```

- `-port`: Server port (default: 8080)
- `-db`: Path to SQLite database (default: data/store.db)

## Development

### Dependencies

```bash
go get github.com/mattn/go-sqlite3
go get github.com/google/uuid
```

### Rebuilding

```bash
go build -o backend-challenge .
```

## Testing

All endpoints have been tested and verified:

- ✅ Product listing
- ✅ Single product retrieval
- ✅ Order placement
- ✅ Coupon validation
- ✅ API key authentication
- ✅ Error handling (404, 400, 401, 422)

## API Documentation

Full OpenAPI specification available at:
https://orderfoodonline.deno.dev/public/openapi.yaml
