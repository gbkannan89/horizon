# Phase 2 — Milestone 1: Docker & Infrastructure

**Context Reminder:** Architecture frozen. Specs frozen. Stack frozen. Monorepo. No Redis/K8s/Grafana/testing.

## Objective

Add the Go monolith server to docker-compose and automate database migrations on startup.

## Backend Audit

Read these files to understand current state:
- `deploy/docker/docker-compose.yml` — currently has postgres only
- `Dockerfile` — builds the monolith binary
- `cmd/migrate/main.go` — migration runner
- `cmd/server/main.go` — server entry point (listens on `PORT` env var, default 8080)
- `.env` — DATABASE_URL=postgres://horizon:horizon@localhost:5433/horizon

## Implementation

1. **Add `horizon-api` service to docker-compose.yml** — build from project root Dockerfile, depends_on postgres (healthy), maps 8081:8080, passes DATABASE_URL pointing to the postgres container

2. **Add `migrate` init container to docker-compose** — runs the migration tool before the API starts, then exits. Can be a separate service with `depends_on: postgres` and `restart: "no"` that runs `go run ./cmd/migrate/`

   OR simpler: add a startup script that runs migrations then starts the server

3. **Restart verification**: `docker-compose up` should start postgres, run migrations, then start the API server. `curl http://localhost:8081/health/live` should respond.

## Files to Modify

| File | Change |
|------|--------|
| `deploy/docker/docker-compose.yml` | Add horizon-api service + migrate runner |

## Verification

- [ ] `docker-compose down -v && docker-compose up` — all services start
- [ ] postgres health check passes
- [ ] Migrations run (check logs for "all migrations applied")
- [ ] API starts and responds on port 8081
- [ ] `curl http://localhost:8081/health/live` returns 200
- [ ] Flutter app can connect to backend via docker network

## Deliverables

Files modified, docker-compose structure, startup sequence.
