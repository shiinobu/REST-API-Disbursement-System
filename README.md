# Disbursement API

A backend REST API for managing fund disbursement requests, built with Go. The system supports authentication, disbursement creation, filtering, approval/rejection workflow, deletion of pending requests, pagination, and CSV export.

## Highlights

- JWT authentication with bcrypt password hashing
- Role-based authorization for approval/rejection actions
- Disbursement lifecycle: `PENDING` → `APPROVED` / `REJECTED`
- Pagination, search, and status filtering
- Business-rule validation for disbursement amounts
- CSV export
- Layered architecture: Handler → Service → Repository → MySQL
- Docker Compose for local development
- Automated tests and `go vet` through GitHub Actions

## Tech Stack

- Go 1.24+
- Gin
- GORM
- MySQL 8.4
- JWT (`golang-jwt/jwt`)
- bcrypt
- Docker & Docker Compose
- GitHub Actions

## Architecture

```text
HTTP Request
    ↓
Handler
    ↓
Service (business rules)
    ↓
Repository
    ↓
GORM
    ↓
MySQL
```

## Project Structure

```text
.
├── cmd/
│   └── main.go
├── internal/
│   ├── config/
│   │   ├── config.go
│   │   └── seed.go
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── disbursement.go
│   │   ├── response.go
│   │   └── validation.go
│   ├── middleware/
│   │   └── auth.go
│   ├── models/
│   │   └── disbursement.go
│   ├── repository/
│   │   ├── disbursement.go
│   │   └── user.go
│   └── services/
│       ├── auth.go
│       ├── disbursement.go
│       └── disbursement_test.go
├── .env.example
├── .gitignore
├── Dockerfile
├── docker-compose.yaml
├── go.mod
└── README.md
```

## Requirements

For local development:

- Go 1.24+
- MySQL 8+

For Docker development:

- Docker Desktop
- Docker Compose

## Installation

### Option 1 — Docker Compose

This is the recommended way to run the project because MySQL is included in the Compose stack.

1. Clone the repository and enter the directory.

2. Create `.env` from the example:

```bash
cp .env.example .env
```

PowerShell:

```powershell
Copy-Item .env.example .env
```

3. Open `.env` and set a strong random value for `JWT_SECRET`.

4. Build and start the application:

```bash
docker compose build
docker compose up -d
```

Or in one command:

```bash
docker compose up -d --build
```

5. Check running containers:

```bash
docker compose ps
```

The API will be available at:

```text
http://localhost:8080
```

Stop the stack:

```bash
docker compose down
```

To remove the database volume as well:

```bash
docker compose down -v
```

### Option 2 — Run Locally

1. Start MySQL and create the database:

```sql
CREATE DATABASE disbursement_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. Create `.env`:

```bash
cp .env.example .env
```

PowerShell:

```powershell
Copy-Item .env.example .env
```

3. Configure the database and JWT secret in `.env`.

4. Download dependencies:

```bash
go mod download
```

5. Run the application:

```bash
go run ./cmd
```

The application automatically runs database migrations and seeds the demo users on startup.

## Environment Variables

Example configuration:

```env
APP_PORT=8080

DB_HOST=mysql
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=disbursement_db

JWT_SECRET=replace_with_a_long_random_secret
JWT_EXPIRES_IN_HOURS=24

GIN_MODE=debug
```

`JWT_SECRET` is required. Do not commit `.env` or real credentials to the repository.

## Authentication

Login:

```http
POST /api/auth/login
Content-Type: application/json
```

Request:

```json
{
  "username": "admin",
  "password": "admin123"
}
```

The response contains a JWT token. Send it using the standard Bearer authentication scheme:

```http
Authorization: Bearer <token>
```

### Demo Users

The application seeds these demo accounts for local development:

| Username | Password | Role |
|---|---|---|
| `admin` | `admin123` | `admin` |
| `operator` | `operator123` | `operator` |

These credentials are for local/demo use only.

## API Endpoints

All disbursement endpoints require a valid JWT unless stated otherwise.

### Authentication

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/api/auth/login` | No | Authenticate user and issue JWT |

### Disbursement

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/disbursements` | List with pagination/filtering |
| POST | `/api/disbursements` | Create a new pending disbursement |
| GET | `/api/disbursements/:id` | Get disbursement detail |
| PATCH | `/api/disbursements/:id/status` | Approve or reject; admin only |
| DELETE | `/api/disbursements/:id` | Delete a pending disbursement |
| GET | `/api/disbursements/export` | Export disbursements as CSV |

### List & Filtering

Default pagination:

```http
GET /api/disbursements?page=1&limit=10
```

Search recipient name:

```http
GET /api/disbursements?search=budi
```

Filter status:

```http
GET /api/disbursements?status=PENDING
```

Combine filters:

```http
GET /api/disbursements?page=1&limit=10&search=budi&status=PENDING
```

Supported statuses:

```text
PENDING
APPROVED
REJECTED
```

### Create Disbursement

```http
POST /api/disbursements
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "recipient_name": "Budi",
  "bank_code": "BCA",
  "account_number": "1234567890",
  "amount": 150000,
  "note": "Pembayaran reimbursement"
}
```

The minimum disbursement amount is greater than `10,000`.

### Approve / Reject

Only users with the `admin` role can change the status.

Approve:

```http
PATCH /api/disbursements/1/status
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "status": "APPROVED",
  "note": "Approved"
}
```

Reject:

```json
{
  "status": "REJECTED",
  "note": "Data rekening tidak valid"
}
```

Only `PENDING` disbursements can be processed or deleted.

### CSV Export

Export all disbursements with a selected status:

```http
GET /api/disbursements/export?status=PENDING
Authorization: Bearer <token>
```

The API returns a CSV attachment.

## Response Format

Success response:

```json
{
  "success": true,
  "message": "Disbursement berhasil dibuat",
  "data": {}
}
```

List response:

```json
{
  "data": [],
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

Error response:

```json
{
  "success": false,
  "message": "Request tidak valid",
  "errors": {
    "field": "pesan validasi"
  }
}
```

## Development

Format code:

```bash
gofmt -w .
```

Run tests:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

## CI

GitHub Actions runs automatically for pushes and pull requests targeting `main`.

The workflow checks:

- `go mod download`
- `go test ./...`
- `go vet ./...`

## Security Notes

- JWT secrets are loaded from environment variables and are not stored in source code.
- Passwords are stored using bcrypt hashes.
- Protected endpoints require a valid JWT.
- Approval/rejection actions require the `admin` role.
- User-controlled SQL values are passed through GORM parameter binding.
- `.env` is excluded from Git through `.gitignore`.

For security-related reports, see `SECURITY.md`.

## Current Scope

This project is intentionally focused on demonstrating backend API development, authentication, authorization, business-rule handling, persistence, pagination, filtering, and export functionality.

Production-scale concerns such as refresh-token rotation, rate limiting, distributed locking, audit logging, observability, and asynchronous payment processing are outside the current scope.

## License

MIT License. See `LICENSE` for details.
