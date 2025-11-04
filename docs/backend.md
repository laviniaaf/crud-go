# Backend Documentation

## Overview

A small Go backend that provides a CRUD API for `bills` (utility bills). The project uses:

- Go (module inside `backend/`)
- Chi router (github.com/go-chi/chi)
- MySQL / MariaDB as the database (driver: github.com/go-sql-driver/mysql)
- godotenv for local `.env` loading

---

## Project layout

- `backend/`
  - `main.go` — server setup, routes and file serving
  - `db.go` — database initialization and table creation
  - `models.go` — data models (Bill)
  - `handlers.go` — HTTP handlers for create/read/update/delete and filtering by date
  - `go.mod`, `go.sum` — Go module files
- `Dockerfile` — multi-stage build for the backend
- `docker-compose.yml` — defines `mariadb` and `backend` services
- `frontend/` — frontend files 
- `.env` — local environment variables for DB connection

---

## Environment variables

The backend reads the following environment variables:

- `DB_HOST` — database hostname (e.g. `mariadb` when running in docker-compose)
- `DB_USER` — database user
- `DB_PASS` — database password
- `DB_NAME` — database name
- `DB_PORT` — database port (default example: `3306`)

---

## Database schema

The backend creates a table called `bills` (if not exists) with the following simplified schema:

- `id` BINARY(16) PRIMARY KEY — UUID stored as raw bytes (BINARY(16)); SQL queries use `UNHEX(REPLACE(uuid, '-', ''))` to convert UUID string to raw bytes
- `embasa` DECIMAL(10,2)
- `coelba` DECIMAL(10,2)
- `created_at` TIMESTAMP (default CURRENT_TIMESTAMP)
- `updated_at` TIMESTAMP (default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP)

Note: When returned via the API, `id` is converted back to the usual UUID string format.

---

## API Endpoints

Base path: `/bills`

- POST `/bills/`
  - Creates a bill
  - Request body: JSON matching `Bill` (except `id` is generated if missing)
  - Response: created `Bill` JSON with `id`

- GET `/bills/` (or `/bills?start=YYYY-MM-DD&end=YYYY-MM-DD`)
  - Returns bills in a date range. If no range provided, returns all bills.
  - Query params: `start` (inclusive, YYYY-MM-DD), `end` (inclusive, YYYY-MM-DD)
  - Response: JSON array of bills

- GET `/bills/all`
  - Returns all bills (alias for read with no date-based filtering)

- PUT `/bills/{id}`
  - Updates the bill with the given UUID
  - Request body: JSON with `embasa` and/or `coelba` (updated_at updated automatically)
  - Response: updated bill JSON

- DELETE `/bills/{id}`
  - Deletes the bill with the given UUID
  - Response: 200 OK

All endpoints use `application/json` for request and response bodies.

---


## Running with Docker Compose 

1. Ensure Docker and Docker Compose are installed.
2. From the project root run:

```
docker compose up --build
```

---

## Tests

There are unit tests for handlers in `backend/handlers_test.go`. Run them with:

```
cd backend
go test ./...
```

---

## Swagger UI

There is a static Swagger UI bundled under `docs/swagger/index.html` that loads the OpenAPI spec at `docs/openapi.yaml`.

To view it locally run:

```bash
cd docs
python3 -m http.server 8000

# then open http://localhost:8000/swagger/index.html
```
---

