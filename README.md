# REST-API-Disbursement-System

REST API sederhana untuk menampilkan daftar pengajuan transfer dana, membuat pengajuan baru, dan menyetujui atau menolaknya.

## Tech Stack

- Go
- Gin Framework
- MySQL
- GORM
- JWT Authentication
- bcrypt password hashing

## Struktur Project

```text
cmd/
  main.go
internal/
  config/
    config.go
    seed.go
  handlers/
    auth.go
    disbursement.go
    response.go
    validation.go
  middleware/
    auth.go
  models/
    disbursement.go
  repository/
    disbursement.go
    user.go
  services/
    auth.go
    disbursement.go
.env.example
go.mod
```

## Setup

1. Buat database MySQL:

```sql
CREATE DATABASE disbursement_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. Salin konfigurasi env:

```bash
cp .env.example .env
```

Jika memakai PowerShell:

```powershell
Copy-Item .env.example .env
```

3. Sesuaikan isi `.env`, lalu install dependency:

```bash
go mod tidy
```

4. Jalankan aplikasi:

```bash
go run cmd/main.go
```

Server berjalan di `http://localhost:8080`.

## Format Response

Success:

```json
{
  "success": true,
  "message": "Disbursement berhasil dibuat",
  "data": {}
}
```

Error:

```json
{
  "success": false,
  "message": "Deskripsi error yang jelas",
  "errors": {
    "field": "pesan validasi"
  }
}
```

## Endpoint

### Auth

Default seeded users:

```text
username: admin
password: admin123
role: admin
```

```text
username: operator
password: operator123
role: operator
```

Login:

```http
POST /api/auth/login
Content-Type: application/json
```

```json
{
  "username": "admin",
  "password": "admin123"
}
```

Gunakan token dari response login untuk endpoint disbursement:

```http
Authorization: Bearer <token>
```

### Disbursement

List:

```http
GET /api/disbursements
```

Filter by status:

```http
GET /api/disbursements?status=PENDING
```

Create:

```http
POST /api/disbursements
Content-Type: application/json
Authorization: Bearer <token>
```

```json
{
  "beneficiary_name": "Budi",
  "bank_name": "BCA",
  "account_number": "1234567890",
  "amount": 150000,
  "description": "Pembayaran reimbursement"
}
```

Detail:

```http
GET /api/disbursements/:id
```

Approve:

```http
PATCH /api/disbursements/:id/approve
```

Reject:

```http
PATCH /api/disbursements/:id/reject
Content-Type: application/json
```

```json
{
  "reason": "Data rekening tidak valid"
}
```
