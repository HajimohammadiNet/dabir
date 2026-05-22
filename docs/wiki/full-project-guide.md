# Dabir

**Dabir** is an open-source letter numbering and registry backend built with **Go** and **PostgreSQL**.

It helps organizations replace spreadsheet-based letter tracking with a clean, auditable, role-based backend system.

Dabir is designed for teams that need a simple but reliable internal system for registering official letters, generating incremental letter numbers, managing users, and tracking all important actions through audit logs.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Roles and Permissions](#roles-and-permissions)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Getting Started](#getting-started)
- [Environment Variables](#environment-variables)
- [Database Migrations](#database-migrations)
- [Initial Setup Wizard](#initial-setup-wizard)
- [Authentication](#authentication)
- [User Management](#user-management)
- [Letter Management](#letter-management)
- [Excel Import](#excel-import)
- [Audit Logs](#audit-logs)
- [Health Checks](#health-checks)
- [Development Commands](#development-commands)
- [API Response Format](#api-response-format)
- [API Documentation](#api-documentation)
- [Letter Numbering Policy](#letter-numbering-policy)
- [Security Notes](#security-notes)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [License](#license)

---

## Overview

Dabir provides a backend API for managing organization letter numbers.

In many organizations, letter numbers are still managed manually in spreadsheets. This approach can become risky over time because it has limited access control, limited auditability, and poor concurrency safety.

Dabir solves this by providing:

- A PostgreSQL-backed letter registry
- Safe incremental numbering
- Role-based access control
- Initial setup wizard
- User management
- Soft delete for letters
- Audit logs for important actions
- Docker-based local development

The current version focuses only on the backend API. A UI/admin panel can be added later.

---

## Features

- Incremental letter numbering
- PostgreSQL-backed storage
- Initial setup wizard
- JWT authentication
- Role-based access control
- Superuser user management
- Letter create, read, update, and soft delete
- Readonly access for users who only need visibility
- Public settings endpoint
- Audit logging
- Health and readiness endpoints
- Clean Architecture-style project structure
- Docker Compose for local development
- Database migrations using `golang-migrate`

---

## Roles and Permissions

Dabir currently supports three roles.

| Role | Description |
|---|---|
| `superuser` | Full access to the system. Can manage users, letters, settings, and audit logs. |
| `editor` | Can create, update, delete, and view letters. Cannot manage users. |
| `readonly` | Can only view letters. Cannot create, update, or delete letters. |

### Permission Matrix

| Action | superuser | editor | readonly |
|---|---:|---:|---:|
| Login | Yes | Yes | Yes |
| View own profile | Yes | Yes | Yes |
| View public settings | Yes | Yes | Yes |
| Create users | Yes | No | No |
| List users | Yes | No | No |
| View user details | Yes | No | No |
| Update users | Yes | No | No |
| Activate/deactivate users | Yes | No | No |
| Create letters | Yes | Yes | No |
| List letters | Yes | Yes | Yes |
| View letter details | Yes | Yes | Yes |
| Update letters | Yes | Yes | No |
| Delete letters | Yes | Yes | No |
| View audit logs | Yes | No | No |

---

## Tech Stack

| Component | Technology |
|---|---|
| Language | Go |
| Database | PostgreSQL |
| Router | Chi |
| Database Driver | pgx / pgxpool |
| Authentication | JWT |
| Password Hashing | bcrypt |
| Migrations | golang-migrate |
| Configuration | Environment variables |
| Containerization | Docker / Docker Compose |

---

## Architecture

Dabir follows a clean and layered architecture.

```text
HTTP Delivery Layer
        ↓
Application Use Cases
        ↓
Domain Layer
        ↑
Infrastructure Layer
```

### Main Layers

| Layer | Responsibility |
|---|---|
| `delivery/http` | HTTP handlers, routes, middleware, response formatting |
| `application` | Use cases and business workflows |
| `domain` | Core entities and repository interfaces |
| `infrastructure` | PostgreSQL repositories, JWT, password hashing |
| `config` | Environment-based configuration |
| `bootstrap` | Application wiring and startup |

---

## Project Structure

```text
dabir/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── application/
│   │   ├── audit/
│   │   ├── auth/
│   │   ├── letters/
│   │   ├── settings/
│   │   ├── setup/
│   │   └── users/
│   │
│   ├── bootstrap/
│   │   └── app.go
│   │
│   ├── config/
│   │   └── config.go
│   │
│   ├── delivery/
│   │   └── http/
│   │       ├── handlers/
│   │       ├── middleware/
│   │       ├── response/
│   │       └── router.go
│   │
│   ├── domain/
│   │   ├── audit/
│   │   ├── letter/
│   │   ├── settings/
│   │   └── user/
│   │
│   ├── infrastructure/
│   │   ├── auth/
│   │   ├── postgres/
│   │   └── security/
│   │
│   └── shared/
│       └── logger/
│
├── migrations/
├── deployments/
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── .dockerignore
├── .env.example
├── .gitignore
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

---

## Requirements

For local development, you need:

- Go
- Docker
- Docker Compose
- PostgreSQL client tools, optional
- `jq`, optional but useful for testing API responses
- `golang-migrate` CLI

Install migrate CLI:

```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Make sure your Go bin path is available:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Check installation:

```bash
migrate -version
```

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/hajimohammadinet/dabir.git
cd dabir
```

### 2. Create environment file

```bash
cp .env.example .env
```

Edit `.env` and set a secure `JWT_SECRET`.

Example:

```env
JWT_SECRET=local-dev-secret-change-in-production
```

### 3. Start services with Docker Compose

```bash
make compose-up
```

This starts:

- PostgreSQL
- Dabir API

### 4. Run database migrations

```bash
make migrate-up
```

### 5. Check health

```bash
curl http://localhost:8080/healthz | jq
```

Expected response:

```json
{
  "success": true,
  "data": {
    "status": "ok"
  }
}
```

Check readiness:

```bash
curl http://localhost:8080/readyz | jq
```

Expected response:

```json
{
  "success": true,
  "data": {
    "database": "ok",
    "status": "ready"
  }
}
```

---

## Environment Variables

Example `.env.example`:

```env
APP_NAME=dabir
APP_ENV=development
APP_HOST=0.0.0.0
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=dabir
DB_PASSWORD=dabir_secret
DB_NAME=dabir
DB_SSLMODE=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=10

JWT_SECRET=change-this-secret
JWT_ACCESS_TOKEN_TTL_MINUTES=60
```

### Application Variables

| Variable | Description | Default |
|---|---|---|
| `APP_NAME` | Application name | `dabir` |
| `APP_ENV` | Application environment | `development` |
| `APP_HOST` | HTTP bind host | `0.0.0.0` |
| `APP_PORT` | HTTP port | `8080` |

### Database Variables

| Variable | Description | Default |
|---|---|---|
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL username | `dabir` |
| `DB_PASSWORD` | PostgreSQL password | `dabir_secret` |
| `DB_NAME` | PostgreSQL database name | `dabir` |
| `DB_SSLMODE` | PostgreSQL SSL mode | `disable` |
| `DB_MAX_OPEN_CONNS` | Max open database connections | `25` |
| `DB_MAX_IDLE_CONNS` | Min/idle database connections | `10` |

### Auth Variables

| Variable | Description | Default |
|---|---|---|
| `JWT_SECRET` | Secret used to sign JWT tokens | Required |
| `JWT_ACCESS_TOKEN_TTL_MINUTES` | Access token lifetime in minutes | `60` |

---

## Database Migrations

Dabir uses `golang-migrate`.

### Run migrations

```bash
make migrate-up
```

### Rollback last migration

```bash
make migrate-down
```

### Check migration version

```bash
make migrate-version
```

### Create a new migration

```bash
make migrate-create NAME=create_some_table
```

### Force migration version

Use this only when you know what you are doing.

```bash
make migrate-force VERSION=1
```

---

## Initial Setup Wizard

Before using the system, Dabir must be initialized.

The setup wizard creates the first `superuser` and stores basic application settings.

### Check setup status

```bash
curl http://localhost:8080/api/v1/setup/status | jq
```

If the app is not initialized, response will look like:

```json
{
  "success": true,
  "data": {
    "initialized": false,
    "setup_needed": true
  }
}
```

### Initialize application

```bash
curl -X POST http://localhost:8080/api/v1/setup/initialize \
  -H "Content-Type: application/json" \
  -d '{
    "organization_name": "Dabir",
    "superuser": {
      "username": "admin",
      "full_name": "System Administrator",
      "password": "Admin123456!"
    },
    "letter_config": {
      "number_prefix": "DABIR",
      "number_padding": 6
    }
  }' | jq
```

Expected response:

```json
{
  "success": true,
  "data": {
    "initialized": true,
    "user_id": "USER_ID",
    "username": "admin"
  }
}
```

After initialization, the setup endpoint cannot be executed again.

---

## Public Settings

Dabir exposes a public settings endpoint.

```bash
curl http://localhost:8080/api/v1/settings/public | jq
```

Example response:

```json
{
  "success": true,
  "data": {
    "organization_name": "Dabir",
    "letter_config": {
      "number_prefix": "DABIR",
      "number_padding": 6
    }
  }
}
```

---

## Authentication

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123456!"
  }' | jq
```

Example response:

```json
{
  "success": true,
  "data": {
    "access_token": "JWT_TOKEN",
    "token_type": "Bearer",
    "expires_in": 3600,
    "user": {
      "id": "USER_ID",
      "username": "admin",
      "full_name": "System Administrator",
      "role": "superuser"
    }
  }
}
```

### Store token

```bash
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "Admin123456!"
  }' | jq -r '.data.access_token')
```

### Get current user

```bash
curl http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq
```

Example response:

```json
{
  "success": true,
  "data": {
    "id": "USER_ID",
    "username": "admin",
    "full_name": "System Administrator",
    "role": "superuser",
    "is_active": true
  }
}
```

---

## User Management

User management endpoints require `superuser` role.

### Create editor user

```bash
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "editor1",
    "full_name": "Editor User",
    "password": "Editor123456!",
    "role": "editor"
  }' | jq
```

### Create readonly user

```bash
curl -X POST http://localhost:8080/api/v1/users/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "username": "readonly1",
    "full_name": "Readonly User",
    "password": "Readonly123456!",
    "role": "readonly"
  }' | jq
```

### List users

```bash
curl http://localhost:8080/api/v1/users/ \
  -H "Authorization: Bearer $TOKEN" | jq
```

### List users with filters

Filter by role:

```bash
curl "http://localhost:8080/api/v1/users/?role=editor" \
  -H "Authorization: Bearer $TOKEN" | jq
```

Search by username or full name:

```bash
curl "http://localhost:8080/api/v1/users/?search=editor" \
  -H "Authorization: Bearer $TOKEN" | jq
```

Filter by active status:

```bash
curl "http://localhost:8080/api/v1/users/?is_active=true" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Get user by ID

```bash
curl "http://localhost:8080/api/v1/users/USER_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Update user

```bash
curl -X PATCH "http://localhost:8080/api/v1/users/USER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "full_name": "Updated User",
    "role": "editor",
    "is_active": true
  }' | jq
```

### Deactivate user

```bash
curl -X PATCH "http://localhost:8080/api/v1/users/USER_ID/deactivate" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Activate user

```bash
curl -X PATCH "http://localhost:8080/api/v1/users/USER_ID/activate" \
  -H "Authorization: Bearer $TOKEN" | jq
```

---

## Letter Management

Letter endpoints require authentication.

### Create a letter

Allowed roles:

- `superuser`
- `editor`

`registrar_name` is automatically set from the authenticated user.  
Do not send `registrar_name` in the request body.

```bash
curl -X POST http://localhost:8080/api/v1/letters/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Contract Review Request",
    "letter_date": "2026-05-19",
    "sender": "Finance Department",
    "receiver": "Legal Department",
    "description": "Request for contract review"
  }' | jq
```

Example response:

```json
{
  "success": true,
  "data": {
    "id": "LETTER_ID",
    "letter_number": 1,
    "formatted_letter_number": "DABIR-000001",
    "title": "Contract Review Request",
    "letter_date": "2026-05-19",
    "registrar_name": "admin",
    "sender": "Finance Department",
    "receiver": "Legal Department",
    "description": "Request for contract review",
    "created_by": "USER_ID",
    "is_deleted": false,
    "created_at": "2026-05-19T10:00:00Z",
    "updated_at": "2026-05-19T10:00:00Z"
  }
}
```

### List letters

Allowed roles:

- `superuser`
- `editor`
- `readonly`

```bash
curl http://localhost:8080/api/v1/letters/ \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Pagination

```bash
curl "http://localhost:8080/api/v1/letters/?page=1&page_size=20" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Search letters

Search matches:

- title
- sender
- receiver
- registrar name
- letter number

```bash
curl "http://localhost:8080/api/v1/letters/?search=contract" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Filter by sender

```bash
curl "http://localhost:8080/api/v1/letters/?sender=Finance" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Filter by receiver

```bash
curl "http://localhost:8080/api/v1/letters/?receiver=Legal" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Filter by registrar name

```bash
curl "http://localhost:8080/api/v1/letters/?registrar_name=Admin" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Filter by letter date

```bash
curl "http://localhost:8080/api/v1/letters/?from_date=2026-01-01&to_date=2026-12-31" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Include deleted letters

By default, deleted letters are hidden.

```bash
curl "http://localhost:8080/api/v1/letters/?include_deleted=true" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Get letter by ID

```bash
curl "http://localhost:8080/api/v1/letters/LETTER_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Update letter

Allowed roles:

- `superuser`
- `editor`

```bash
curl -X PATCH "http://localhost:8080/api/v1/letters/LETTER_ID" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "title": "Updated Contract Review Request",
    "letter_date": "2026-05-19",
    "sender": "Finance Department",
    "receiver": "Legal Department",
    "description": "Updated description"
  }' | jq
```

### Delete letter

Dabir uses soft delete for letters.

Allowed roles:

- `superuser`
- `editor`

```bash
curl -X DELETE "http://localhost:8080/api/v1/letters/LETTER_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

Example response:

```json
{
  "success": true,
  "data": {
    "deleted": true
  }
}
```

---

## Excel Import

Dabir supports importing existing letters from Excel files.

This is useful when migrating from spreadsheet-based letter tracking.

Currently supported format:

- `.xlsx`

Import is available only for `superuser`.

### Required columns

The Excel file must contain these logical fields:

| Field | Supported column names |
|---|---|
| `letter_number` | `letter_number`, `number`, `no`, `شماره نامه`, `شماره` |
| `title` | `title`, `subject`, `عنوان`, `عنوان نامه`, `موضوع` |
| `letter_date` | `letter_date`, `date`, `تاریخ`, `تاریخ نامه` |
| `sender` | `sender`, `from`, `فرستنده`, `ارسال کننده` |
| `receiver` | `receiver`, `to`, `گیرنده`, `دریافت کننده`, `مقصد` |

`registrar_name` is not read from Excel.  
It is automatically set to the username of the user who commits the import.

### Import flow

The import process has two steps:

```text
Preview Excel file
    ↓
Review detected columns, valid rows, errors, and duplicates
    ↓
Commit import
    ↓
Letters are inserted into database
    ↓
Letter number sequence continues from the maximum imported number
```

### Preview import

```bash
curl -X POST http://localhost:8080/api/v1/imports/letters/preview \
  -H "Authorization: Bearer $TOKEN" \
  -F "file=@letters_import.xlsx" | jq
```

Example response:

```json
{
  "success": true,
  "data": {
    "id": "IMPORT_ID",
    "type": "letters",
    "status": "previewed",
    "file_name": "letters_import.xlsx",
    "total_rows": 2,
    "valid_rows": 2,
    "invalid_rows": 0,
    "max_letter_number": 1002,
    "detected_columns": {
      "letter_number": "شماره نامه",
      "title": "عنوان نامه",
      "letter_date": "تاریخ نامه",
      "sender": "فرستنده",
      "receiver": "گیرنده"
    },
    "preview_data": [
      {
        "row_number": 2,
        "letter_number": 1001,
        "title": "نامه تست ۱",
        "letter_date": "2026-05-19",
        "sender": "شرکت الف",
        "receiver": "شرکت ب"
      }
    ],
    "errors": []
  }
}
```

### Get import job

```bash
curl "http://localhost:8080/api/v1/imports/IMPORT_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Commit import

```bash
curl -X POST "http://localhost:8080/api/v1/imports/letters/IMPORT_ID/commit" \
  -H "Authorization: Bearer $TOKEN" | jq
```

Example response:
```json
{
  "success": true,
  "data": {
    "import_id": "IMPORT_ID",
    "imported_rows": 2,
    "skipped_rows": 0,
    "next_letter_number": 1003
  }
}
```

### Duplicate detection

During preview, Dabir detects:

- Duplicate letter numbers inside the Excel file
- Letter numbers that already exist in the database

If invalid rows exist, commit is blocked.

Example error:

```json
{
  "row": 3,
  "field": "letter_number",
  "message": "duplicate letter number in file, first seen at row 2"
}
```

---

## Audit Logs

Audit log access requires `superuser` role.

### List audit logs

```bash
curl http://localhost:8080/api/v1/audit-logs/ \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Filter audit logs by action

```bash
curl "http://localhost:8080/api/v1/audit-logs/?action=letter.created" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Filter audit logs by entity type

```bash
curl "http://localhost:8080/api/v1/audit-logs/?entity_type=letter" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Filter audit logs by actor user ID

```bash
curl "http://localhost:8080/api/v1/audit-logs/?actor_user_id=USER_ID" \
  -H "Authorization: Bearer $TOKEN" | jq
```

### Current audit actions

```text
setup.initialized
auth.login_success
auth.login_failed
user.created
user.updated
user.activated
user.deactivated
letter.created
letter.updated
letter.deleted
letters.import_previewed
letters.import_committed
```

---

## Health Checks

### Liveness

```bash
curl http://localhost:8080/healthz | jq
```

### Readiness

```bash
curl http://localhost:8080/readyz | jq
```

The readiness endpoint checks the database connection.

---

## Development Commands

### Run locally

If PostgreSQL is already running:

```bash
export $(grep -v '^#' .env | xargs)
make run
```

### Run tests

```bash
make test
```

### Format code

```bash
make fmt
```

### Tidy modules

```bash
make tidy
```

### Build binary

```bash
make build
```

### Start PostgreSQL only

```bash
make postgres-up
```

### Stop PostgreSQL

```bash
make postgres-down
```

### Start full stack

```bash
make compose-up
```

### Stop full stack

```bash
make compose-down
```

### View logs

```bash
make compose-logs
```

### View API logs

```bash
make api-logs
```

### Reset development database

Warning: this removes the local PostgreSQL volume.

```bash
make dev-reset
```

---

## API Response Format

Dabir uses a consistent JSON response format.

### Success response

```json
{
  "success": true,
  "data": {}
}
```

### Error response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message"
  }
}
```

---

## API Documentation

The OpenAPI specification is available at:

```text
docs/openapi.yaml
```

You can preview it using Swagger Editor, Redoc, or any OpenAPI-compatible tool.

### Validate OpenAPI file

Using Docker:
```bash
docker run --rm -v "$PWD:/work" redocly/cli lint /work/docs/openapi.yaml
```

Using npm:
```bash
npx @redocly/cli lint docs/openapi.yaml
```

---

## Letter Numbering Policy

Dabir currently uses a PostgreSQL sequence for letter numbers.

This is concurrency-safe and production-friendly.

Example:

```text
DABIR-000001
DABIR-000002
DABIR-000003
```

The raw number is stored in the database as `letter_number`.

The formatted number is generated by the application using:

- `number_prefix`
- `number_padding`

### Important note about gaps

PostgreSQL sequences are not gapless.

If a request receives a number from the sequence but fails before the letter is saved, that number may be skipped.

This is normal behavior for sequence-based systems and is useful for high-concurrency safety.

If strict gapless numbering is required in the future, Dabir can implement a locked counter table using database transactions and `SELECT ... FOR UPDATE`.

---

## Security Notes

For production usage:

- Change `JWT_SECRET`
- Use a strong database password
- Use HTTPS behind a reverse proxy
- Do not expose PostgreSQL to the public internet
- Keep `.env` out of Git
- Rotate secrets periodically
- Use proper backup for PostgreSQL
- Review audit logs regularly
- Consider enabling rate limiting on login endpoints
- Consider adding refresh tokens and token revocation for advanced deployments

---

## Roadmap

Planned improvements:

- OpenAPI / Swagger documentation
- Request validation improvements
- Request logging middleware
- CORS configuration for future UI
- Better error codes
- Audit logging refactor into application layer
- Refresh token support
- Password change endpoint
- Forgot password flow
- Organization settings management
- Letter numbering strategies
- Export letters to CSV/XLSX
- Import initial letters from Excel/CSV
- Admin UI
- Helm chart
- CI/CD with GitHub Actions

---

## Contributing

Contributions are welcome.

Recommended workflow:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests
5. Open a pull request

```bash
make fmt
make tidy
make test
```

Please keep the project structure clean and follow the existing architecture.

---

## License

Dabir is released under the Apache License 2.0.

See the `LICENSE` file for details.
