# Implementation Plan: Food Ordering API

## Overview
Build a Go-based API server implementing the OpenAPI 3.1 spec for a food ordering system with promo code validation.

## Technology Stack
- **Language**: Go 1.21+
- **Database**: SQLite3
- **Router**: `chi` or `gorilla/mux`
- **Database Driver**: `mattn/go-sqlite3`
- **UUID Generation**: `google/uuid`

## Database Schema

### Products Table
```sql
CREATE TABLE products (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    price REAL NOT NULL,
    image_thumbnail TEXT,
    image_mobile TEXT,
    image_tablet TEXT,
    image_desktop TEXT
);
```

### Valid Coupons Table
```sql
CREATE TABLE valid_coupons (
    code TEXT PRIMARY KEY
);
```

## API Endpoints

### GET /api/product
List all products

### GET /api/product/{productId}
Get single product by ID
- 404 if not found
- 400 for invalid ID format

### POST /api/order
Create order with optional coupon validation
- Request body: OrderReq (items array, optional couponCode)
- Response: Order object with UUID, items, and product details
- Validation:
  - All productIds must exist
  - Quantities must be positive
  - Coupon code must be valid (if provided)
- Error codes: 400 (bad request), 422 (validation error)

## Coupon Processing

### What We Actually Did
Created standalone Python solution in `coupon/` directory:

1. **Script**: `process_coupons.py`
   - Streams through 3 gzipped files (~600MB each)
   - Filters codes 8-10 characters long
   - Uses Python sets for memory-efficient storage
   - Performs set intersections to find codes in 2+ files

2. **Performance**:
   - Runtime: ~3 hours on M3 Pro
   - Memory: 15-18GB peak
   - Processing rate: ~5M lines/min

3. **Results**: 8 valid codes found:
   - BIRTHDAY, BUYGETON, FIFTYOFF, FREEZAAA
   - GNULINUX, HAPPYHRS, OVER9000, SIXTYOFF

4. **Integration**: Import `coupon/valid_coupons.txt` into SQLite at startup

See `coupon/README.md` for reproduction steps.

## Project Structure
```
backend-challenge/
├── main.go                 # Entry point
├── api/                    # HTTP handlers, routes
├── db/                     # Database setup, queries
├── models/                 # Data structures
├── coupon/                 # Coupon processing (standalone)
│   ├── process_coupons.py  # Main script
│   ├── valid_coupons.txt   # Results (8 codes)
│   └── README.md           # Reproduction steps
└── data/
    └── store.db            # SQLite database
```

## Testing
1. Test product endpoints
2. Verify valid coupons (HAPPYHRS, FIFTYOFF, etc.)
3. Verify invalid coupons rejected
4. Test order creation with/without coupons

## Deliverables
- Working Go server on port 8080
- SQLite database with products and valid coupons
- All endpoints matching OpenAPI spec
- Clean, maintainable code
