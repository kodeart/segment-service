# segment-service

DMP Manage Service - a CRUD HTTP API for managing audience "Segments."

## Requirements

- Go **1.25+**
- Docker + Docker Compose

## Quick start

    cp .env.example .env
    docker compose up -d
    go run cmd/server/main.go

Server runs on port 8070 by default.

---

## Configuration

Create a local `.env` file (you can start from `.env.example`).

Configure if needed:

- `SERVER_PORT` — HTTP port for the service (default: 8070)
- `POSTGRES_HOST` — Postgres host (use `127.0.0.1` for local dev)
- `POSTGRES_PORT` — Postgres port for `segment-db` (default: 5433)
- `POSTGRES_TEST_PORT` — Postgres port for `segment-test-db` (default: 5435)
- `POSTGRES_USER`
- `POSTGRES_PASSWORD`
- `POSTGRES_DB`


### Database (Docker)

This repo provides two Postgres containers:

- `segment-db` — development database
- `segment-test-db` — database used by tests

Start both:

    docker compose up -d


#### Schema initialization

The schema runs in `/docker-entrypoint-initdb.d` **only on first initialization**.
On `schema.sql` changes, you must wipe the data directories if you want Postgres to re-run initialization:

    docker compose down -v
    rm -rf ./tmp/pgdata ./tmp/pgdata_test
    docker compose up -d

## Testing

Docker containers must run:

    go test -v ./...
