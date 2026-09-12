![CI](https://github.com/hajimohammadinet/dabir/actions/workflows/ci.yml/badge.svg)

# Dabir

**Dabir** is an open-source letter numbering and registry system for organizations that need a simple, auditable, and structured replacement for spreadsheet-based letter tracking.

It provides a clean backend API, a web-based admin panel, incoming and outgoing letter registries, role-based access control, Jalali date support, Excel migration, manual numbering, scanned letter attachments, S3-compatible object storage integration, and Kubernetes-ready deployment manifests.

---

## Features

- Incoming letter registry
- Outgoing letter registry
- Create, edit, preview, soft-delete, and search letters
- Multiple numbering modes:
  - Fixed prefix numbering, such as `DABIR-000001`
  - Jalali yearly numbering, such as `405-0001`
  - Manual numbering, such as `405-158`, `405-ق-103`, or any custom structure
- Independent direction-aware numbering and smart manual number suggestions
- Support for Persian and Arabic digits in date and number inputs
- Jalali / Persian calendar support
- Persian and English UI foundation
- RTL-friendly user interface
- Role-based access control
- User management
- Password change and password reset
- Excel import for migrating existing letters
- Import preview and duplicate detection
- Preserving original imported letter numbers
- Audit logs for important system actions
- Setup wizard for first-time initialization
- Optional scanned letter / attachment upload
- S3-compatible object storage support for attachments
- Private attachment storage with presigned download/view URLs
- Attachment management: upload, view/download, and delete
- Dark mode support
- REST API with OpenAPI documentation
- Docker Compose deployment
- Kubernetes deployment with Helm
- External PostgreSQL support
- GitHub Actions CI

---

## Roles

Dabir currently supports three access levels:

| Role | Description |
|---|---|
| `superuser` | Full system access, user management, imports, audit logs, settings, and attachment management |
| `editor` | Can create, edit, delete, and view letters; can manage attachments |
| `readonly` | Can view letters and attachments |

---

## Tech Stack

### Backend

- Go
- PostgreSQL
- Chi Router
- JWT authentication
- Clean architecture style
- Database migrations
- S3-compatible object storage client

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
- S3-compatible object storage, such as MinIO

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

S3_ENDPOINT=http://localhost:9000
S3_PUBLIC_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
S3_BUCKET=dabir-attachments
S3_ACCESS_KEY=dabir
S3_SECRET_KEY=dabir_minio_secret
S3_USE_SSL=false
S3_FORCE_PATH_STYLE=true
S3_PRESIGNED_URL_TTL_MINUTES=15
S3_MAX_UPLOAD_SIZE_MB=20
```

Start services:

```bash
make compose-up
```

Compose connects the API to MinIO at `http://minio:9000` and uses
`S3_PUBLIC_ENDPOINT` when signing browser-facing attachment URLs. It also waits
for MinIO, creates `S3_BUCKET` if needed, and explicitly keeps the bucket
private. No manual bucket setup is required.

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

Start PostgreSQL, MinIO, the API, and the web application:

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

MinIO Console URL:

```text
http://localhost:9001
```

Default local MinIO credentials, unless changed in `.env`:

```text
username: dabir
password: dabir_minio_secret
```

---

## Numbering Modes

Dabir supports multiple numbering modes to fit different organizational workflows.

Incoming and outgoing letters use independent fixed and yearly sequence namespaces. Manual number duplicate checks and smart suggestions are also scoped to direction, so the same manual number can exist once in each registry without one direction influencing the other. Existing installations automatically classify pre-feature rows as incoming letters.

### Fixed Prefix

Example:

```text
DABIR-000001
DABIR-000002
```

### Jalali Yearly

Example:

```text
405-0001
405-0002
406-0001
```

### Manual Numbering

Manual numbering stores the letter number exactly as entered by the user.

Examples:

```text
405-158
405-ق-103
1401-MM-DM-95
HR-2026-0042
```

Manual numbering also includes smart suggestions. When a user starts typing a prefix, Dabir searches for the latest similar number and suggests the next one.

Examples:

```text
Input prefix: 405-
Latest similar number: 405-158
Suggested number: 405-159
```

```text
Input prefix: 405-ق-
Latest similar number: 405-ق-102
Suggested number: 405-ق-103
```

Persian and Arabic digits are normalized for suggestion lookup, so inputs such as `۴۰۵-ق-` are also supported.

---

## Letter Attachments / Scanned Letters

Dabir supports optional scanned letter attachments.

Supported file types:

```text
PDF
JPG / JPEG
PNG
```

Attachment flow:

```text
Create or edit a letter
Upload one or more optional files
Store files in S3-compatible object storage
Store only metadata in PostgreSQL
View/download files through short-lived presigned URLs
Delete attachments when needed
```

Attachments are stored privately in object storage. Dabir does not make the bucket public. The backend generates presigned URLs for viewing or downloading files.

---

## Object Storage / S3 Configuration

Dabir uses S3-compatible object storage for attachments. MinIO is recommended for local development and small to medium on-prem deployments.

Required environment variables:

```env
S3_ENDPOINT=http://minio:9000
S3_PUBLIC_ENDPOINT=http://localhost:9000
S3_REGION=us-east-1
S3_BUCKET=dabir-attachments
S3_ACCESS_KEY=dabir
S3_SECRET_KEY=dabir_minio_secret
S3_USE_SSL=false
S3_FORCE_PATH_STYLE=true
S3_PRESIGNED_URL_TTL_MINUTES=15
S3_MAX_UPLOAD_SIZE_MB=20
```

`S3_ENDPOINT` is the backend's network path to object storage.
`S3_PUBLIC_ENDPOINT` is optional and defaults to `S3_ENDPOINT`; set it whenever
the browser reaches the same S3 service through a different host or scheme. In
Docker Compose, those values are `http://minio:9000` and
`http://localhost:9000` respectively. The backend signs directly against the
public endpoint, so the signature remains valid without exposing credentials
or rewriting URLs in the frontend.

The Compose environment uses the official Quay-hosted MinIO server
`RELEASE.2025-09-07T16-13-09Z` and MinIO client
`RELEASE.2025-08-13T08-35-41Z`. The one-shot `minio-init` service waits for a
healthy server, creates the configured bucket idempotently, and removes any
anonymous access policy.

For production, use a private bucket and a dedicated access key with limited permissions to only the required bucket.

Recommended production notes:

```text
- Keep the bucket private
- Use TLS for the S3 endpoint
- Use a dedicated access key for Dabir
- Configure object storage backup or replication
- Keep S3_MAX_UPLOAD_SIZE_MB aligned with the ingress upload limit
```

For NGINX Ingress, make sure the upload size is higher than `S3_MAX_UPLOAD_SIZE_MB`:

```yaml
nginx.ingress.kubernetes.io/proxy-body-size: "25m"
```

---

## Kubernetes Deployment

Dabir includes a Helm chart:

```text
charts/dabir
```

Dabir does **not** deploy PostgreSQL inside Kubernetes by default.
It expects an external PostgreSQL database.

Dabir also expects an S3-compatible object storage endpoint for attachments, such as MinIO, Ceph RGW, or a cloud S3-compatible service.

Create namespace:

```bash
kubectl create namespace dabir
```

Create secret:

```bash
kubectl -n dabir create secret generic dabir-secret \
  --from-literal=DB_PASSWORD='YOUR_DB_PASSWORD' \
  --from-literal=JWT_SECRET='YOUR_LONG_RANDOM_JWT_SECRET' \
  --from-literal=S3_ACCESS_KEY='YOUR_S3_ACCESS_KEY' \
  --from-literal=S3_SECRET_KEY='YOUR_S3_SECRET_KEY'
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

## Helm S3 Values

The Helm chart supports S3 configuration through API environment variables and secrets.

Example values:

```yaml
api:
  env:
    S3_ENDPOINT: "https://s3.example.com"
    S3_PUBLIC_ENDPOINT: "https://s3.example.com"
    S3_REGION: "us-east-1"
    S3_BUCKET: "dabir-attachments"
    S3_USE_SSL: "true"
    S3_FORCE_PATH_STYLE: "true"
    S3_PRESIGNED_URL_TTL_MINUTES: "15"
    S3_MAX_UPLOAD_SIZE_MB: "20"

  secrets:
    create: true
    existingSecret: ""
    DB_PASSWORD: "change-this-db-password"
    JWT_SECRET: "change-this-jwt-secret"
    S3_ACCESS_KEY: "change-this-s3-access-key"
    S3_SECRET_KEY: "change-this-s3-secret-key"
```

When using an existing secret, it must contain:

```text
DB_PASSWORD
JWT_SECRET
S3_ACCESS_KEY
S3_SECRET_KEY
```

---

## Excel Import

Dabir can import existing incoming letters from Excel files.

Excel import remains intentionally incoming-only in this release. Outgoing letters can be created through the API or web UI; the existing import behavior and data mapping are unchanged.

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
Preserve original letter numbers
Continue numbering from the imported data
```

The importer supports real-world legacy files and can preserve original manual letter numbers instead of forcing them into a generated numeric format.

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
  "direction": "incoming",
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

Check rendered S3 environment variables:

```bash
helm template dabir charts/dabir -n dabir > /tmp/dabir-rendered.yaml
grep -n "S3_" /tmp/dabir-rendered.yaml
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
- Full i18n coverage
- Migration Job support in Helm
- Release automation
- Attachment thumbnails and richer preview experience
- Export letters to Excel / PDF

---

## License

This project is licensed under the Apache License 2.0.
