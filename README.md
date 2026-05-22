# Dabir

**Dabir** is an open-source letter numbering and registry system built for organizations that need a simple, auditable, and structured replacement for spreadsheet-based letter tracking.

It provides a clean backend API, a web-based admin panel, role-based access control, Excel migration support, Jalali date support, and Kubernetes-ready deployment manifests.

---

## Features

- Letter numbering and registry management
- Automatic incremental letter numbers
- Jalali / Persian calendar support
- Persian and English UI foundation
- RTL-friendly user interface
- Role-based access control
- User management
- Password change and password reset
- Excel import for migrating existing letters
- Import preview and duplicate detection
- Audit logs for important system actions
- Setup wizard for first-time initialization
- Dark mode support
- REST API with OpenAPI documentation
- Docker Compose deployment
- Kubernetes deployment with Helm
- External PostgreSQL support

---

## Roles

Dabir currently supports three access levels:

| Role | Description |
|---|---|
| `superuser` | Full system access, user management, imports, audit logs, settings |
| `editor` | Can create, edit, delete, and view letters |
| `readonly` | Can only view letters |

---

## Tech Stack

### Backend

- Go
- PostgreSQL
- Chi Router
- JWT authentication
- Clean architecture style
- Database migrations

### Frontend

- Next.js
- TypeScript
- Tailwind CSS
- shadcn/ui-style components
- Jalali date picker
- Dark mode
- Persian / English localization foundation

### Deployment

- Docker Compose
- Helm Chart
- Kubernetes
- External PostgreSQL

---

## Project Structure

```text
.
├── cmd/
│   └── api/
├── internal/
│   ├── application/
│   ├── domain/
│   ├── delivery/
│   ├── infrastructure/
│   └── shared/
├── migrations/
├── web/
├── charts/
│   └── dabir/
├── deployments/
├── docs/
│   ├── openapi.yaml
│   └── wiki/
└── README.md
```

---

## Quick Start with Docker Compose

Create environment file:

```bash
cp .env.example .env
```

Edit important values:

```env
DB_PASSWORD=change-this-db-password
JWT_SECRET=change-this-secret-in-production
```

Start services:

```bash
make compose-up
```

Run migrations:

```bash
make migrate-up
```

Open the web UI:

```text
http://localhost:3000
```

On the first run, Dabir redirects to the setup page where you can create the first superuser.

---

## Local Development

Start PostgreSQL:

```bash
make compose-up
```

Run backend:

```bash
export $(grep -v '^#' .env | xargs)
make run
```

Run frontend:

```bash
cd web
npm install
npm run dev
```

Frontend URL:

```text
http://localhost:3000
```

Backend URL:

```text
http://localhost:8080
```

---

## Kubernetes Deployment

Dabir includes a Helm chart:

```text
charts/dabir
```

Dabir does **not** deploy PostgreSQL inside Kubernetes by default.
It expects an external PostgreSQL database.

Create namespace:

```bash
kubectl create namespace dabir
```

Create secret:

```bash
kubectl -n dabir create secret generic dabir-secret \
  --from-literal=DB_PASSWORD='YOUR_DB_PASSWORD' \
  --from-literal=JWT_SECRET='YOUR_LONG_RANDOM_JWT_SECRET'
```

Install with Helm:

```bash
helm upgrade --install dabir charts/dabir \
  -n dabir \
  --create-namespace \
  -f charts/dabir/examples/values-prod.yaml
```

Run database migrations separately against the external PostgreSQL database.

---

## Excel Import

Dabir can import existing letters from Excel files.

Supported format:

```text
.xlsx
```

Supported logical columns:

| Field | Supported column names |
|---|---|
| `letter_number` | `letter_number`, `number`, `no`, `شماره نامه`, `شماره` |
| `title` | `title`, `subject`, `عنوان`, `عنوان نامه`, `موضوع` |
| `letter_date` | `letter_date`, `date`, `تاریخ`, `تاریخ نامه` |
| `sender` | `sender`, `from`, `فرستنده`, `ارسال کننده` |
| `receiver` | `receiver`, `to`, `گیرنده`, `دریافت کننده`, `مقصد` |

Import flow:

```text
Upload Excel
Preview and validate
Detect duplicates
Commit import
Continue numbering from the maximum imported number
```

---

## Jalali Date Support

Dabir supports official Iranian Jalali dates in the UI.

Accepted input examples:

```text
1405/03/01
۱۴۰۵/۰۳/۰۱
```

Dates are stored in PostgreSQL as standard `DATE` values and returned by the API in both Gregorian and Jalali formats.

Example API response:

```json
{
  "letter_date": "2026-05-22",
  "letter_date_jalali": "1405/03/01"
}
```

---

## API Documentation

OpenAPI specification:

```text
docs/openapi.yaml
```

Validate OpenAPI:

```bash
docker run --rm -v "$PWD:/work" redocly/cli lint /work/docs/openapi.yaml
```

---

## Useful Commands

Run backend tests:

```bash
go test $(go list ./... | grep -v '/web/')
```

Run frontend checks:

```bash
cd web
npm run lint
npm run build
```

Lint Helm chart:

```bash
helm lint charts/dabir
```

Render Helm templates:

```bash
helm template dabir charts/dabir -n dabir
```

---

## Documentation

More detailed documentation is available in:

```text
docs/
```

The previous full development guide is archived here:

```text
docs/wiki/full-project-guide.md
```

---

## Roadmap

- Dashboard statistics endpoint
- Better runtime configuration for web deployment
- Print / export letter confirmation
- Advanced letter search and filters
- Editable settings page
- Full i18n coverage
- Migration Job support in Helm
- GitHub Actions CI
- Release automation

---

## License

This project is licensed under the Apache License 2.0.
