# Backend Challenge - Food Ordering API

Go-based REST API implementing OpenAPI 3.1 spec for food ordering with coupon validation.

## Quick Start

### Prerequisites

- Go 1.21+
- SQLite3

### Setup & Run

```bash
# Initialize database (pre-populated with products and coupons)
make init

# Build and run
make run
```

Server starts at `http://localhost:8080`

### Test API

```bash
# List products
curl http://localhost:8080/api/product

# Get product
curl http://localhost:8080/api/product/1

# Place order (requires api_key header)
curl -X POST http://localhost:8080/api/order \
  -H "Content-Type: application/json" \
  -H "api_key: apitest" \
  -d '{"items":[{"productId":"1","quantity":2}],"couponCode":"HAPPYHRS"}'
```

## API Endpoints

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| `/api/product` | GET | No | List all products (supports `?limit=N&offset=N`) |
| `/api/product/{id}` | GET | No | Get product by ID |
| `/api/order` | POST | Yes | Place order with optional coupon |
| `/health` | GET | No | Health check endpoint |
| `/public/openapi.yaml` | GET | No | OpenAPI specification |

### Request/Response Examples

**POST /api/order**
```json
{
  "items": [{"productId": "1", "quantity": 2}],
  "couponCode": "HAPPYHRS"  // optional
}
```

**Response**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "items": [{"productId": "1", "quantity": 2}],
  "products": [{...}],
  "couponCode": "HAPPYHRS"
}
```

**Error Response**
```json
{
  "code": 422,
  "type": "validation_error",
  "message": "Invalid coupon code"
}
```

## Configuration

### Command-line Flags

```bash
./backend-challenge -port 8080
```

- `-port`: Server port (default: `8080`)

### Environment Variables

- `API_KEY`: Authentication key (default: `apitest`)

```bash
API_KEY=custom_key ./backend-challenge
```

## Coupon Validation

Valid coupons must:
1. Be 8-10 characters long
2. Appear in ≥2 of the 3 coupon files

**Valid coupons:** `BIRTHDAY`, `BUYGETON`, `FIFTYOFF`, `FREEZAAA`, `GNULINUX`, `HAPPYHRS`, `OVER9000`, `SIXTYOFF`

**Case sensitivity:** Coupons are case-sensitive. All valid coupons in the database are uppercase.

Database is pre-populated with these coupons. See `coupon/README.md` to reproduce the preprocessing.

## Development

### Dependencies

```bash
go mod download
```

Core dependencies:
- `github.com/mattn/go-sqlite3` - SQLite driver
- `github.com/google/uuid` - UUID generation
- `github.com/stretchr/testify` - Testing assertions
- `go.uber.org/mock` - Mock generation

### Building

```bash
make build           # Build binary
make run             # Build and run
make clean           # Clean artifacts
```

### Testing

```bash
make test            # Run all tests
make test-coverage   # Generate coverage report (80%)
make generate        # Regenerate mocks
```

**Coverage:** 80.0% total (API: 96.5%, Service: 100%, DB: 85.7%)

### Project Structure

```
backend-challenge/
├── main.go              # Entry point, server lifecycle
├── api/                 # HTTP layer
│   ├── handlers.go      # Request handlers
│   ├── middleware.go    # Auth, CORS, request ID
│   └── router.go        # Route definitions
├── service/             # Business logic
│   ├── service.go       # Order processing
│   └── errors.go        # Domain errors
├── db/                  # Data layer
│   ├── db.go            # Connection management
│   ├── queries.go       # SQL queries
│   ├── interface.go     # Database interface
│   └── mocks/           # Generated mocks
├── models/              # Data structures
│   └── models.go        # Product, Order, etc.
├── coupon/              # Coupon preprocessing (standalone)
└── data/
    ├── init.sql         # Schema and seed data
    └── store.db         # SQLite database
```

## Changes from demo API

The [demo API](https://orderfoodonline.deno.dev/api) intentionally omits edge case handling. This implementation addresses:

**Edge cases handled:**
- **Coupon validation**: Demo accepts any coupon code (including invalid ones). This implementation validates against the preprocessed coupon database.
- **Negative quantities**: Demo accepts negative values (e.g., -5). Returns `400` for quantities ≤ 0.
- **Empty product ID**: Demo returns product list. Returns `400` for missing/empty ID.

**HTTP status code semantics:**
- `400` - Malformed request (invalid JSON, empty items, missing productId, non-positive quantity, empty product ID)
- `404` - Product not found (GET endpoint only)
- `422` - Validation error (invalid coupon, product doesn't exist in order)
- `500` - Server error (database failures)

**Implementation notes:**
- **Product IDs**: OpenAPI spec defines `productId` as `integer/int64`, but demo API uses strings (e.g., `"1"`). We follow the demo's string implementation for consistency with existing data.
- **Coupon case sensitivity**: Not specified in requirements. Implementation treats coupons as case-sensitive (all valid coupons are uppercase).

## Design Decisions

### Architecture

**3-Layer Clean Architecture**
- **API Layer:** HTTP handling, validation, serialization
- **Service Layer:** Business logic, orchestration
- **DB Layer:** Data access, queries

**Why:** Clear separation of concerns, testable, maintainable. Each layer has single responsibility and can be tested in isolation, though perhaps excessive for a small web server.

### Database: SQLite

**Why SQLite:**
- Embedded, zero-config
- Perfect for local/demo use
- Simple file-based deployment
- Can commit database for even easier setup

### Order Handling

Orders are not persisted to the database. They are validated, assigned a UUID, and returned immediately.

**Why:** OpenAPI spec has no order retrieval endpoints (no `GET /order` or `GET /order/{id}`). For the purposes of this demo, there's no need to store them.

### Coupon Preprocessing: Standalone Python Script

**Why Python:**
- Simpler set operations (`set` type built-in)
- Faster prototyping for one-time processing
- Result stored in DB—preprocessing doesn't run per request

### Graceful Shutdown: Signal-Based Context

Uses `signal.NotifyContext` for clean shutdown on SIGINT/SIGTERM. Prevents request interruption.

### Testing Strategy

**~80% coverage achieved with:**
- **Unit tests:** Handlers, service, DB queries (mocked)
- **Integration tests:** Full stack with test DB
- **Table-driven tests:** Concise, parameterized test cases

**Uncovered paths:** Mostly error branches (DB connection failures, shutdown errors) - low ROI to test.

### Error Handling

**Typed domain errors** (`service.ErrInvalidCoupon`, `service.ErrProductNotFound`)
- Enables error type checking with `errors.Is()`
- Separates business logic errors from infrastructure errors
- Allows precise HTTP status mapping (400 vs 422 vs 500)

### Middleware Stack

**Order:** MaxBodySize → CORS → RequestID → Auth (per-route)

**Why:**
- Body limit first (DoS protection)
- CORS early (preflight support)
- Request ID for traceability
- Auth only on the POST order endpoint

### Pagination

**Pagination:** Not in spec but added query params (`?limit=100&offset=0`) with 100-item cap for pagination.

---

