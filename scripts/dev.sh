#!/usr/bin/env bash
set -euo pipefail

echo "=== Starting Horizon Development Environment ==="

# Start infrastructure
docker compose -f deploy/docker/docker-compose.yml up -d

# Wait for PostgreSQL
echo "Waiting for PostgreSQL..."
until docker compose exec postgres pg_isready -U horizon 2>/dev/null; do
  sleep 2
done
echo "PostgreSQL ready."

# Wait for Redis
echo "Waiting for Redis..."
until docker compose exec redis redis-cli ping 2>/dev/null; do
  sleep 1
done
echo "Redis ready."

# Wait for NATS
echo "Waiting for NATS..."
until curl -sf http://localhost:8222/healthz 2>/dev/null; do
  sleep 2
done
echo "NATS ready."

echo ""
echo "=== Development Environment Ready ==="
echo "PostgreSQL: localhost:5432"
echo "Redis:      localhost:6379"
echo "NATS:       localhost:4222"
echo "MinIO:      localhost:9000 (console: localhost:9001)"
echo "Prometheus: localhost:9090"
echo "Grafana:    localhost:3000"
echo ""
echo "Run 'make test' to run tests."
