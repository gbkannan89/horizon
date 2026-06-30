# Development Setup

## Prerequisites

- Go 1.24+
- Docker and Docker Compose
- Git

## Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd horizon

# Run setup
make setup

# Start development environment
make dev

# Run tests
make test
```

## Development Environment

The development environment runs in Docker containers:

- **PostgreSQL 17** — Database (port 5432)
- **Redis 7** — Cache (port 6379)
- **NATS 2** — Event Bus (port 4222)
- **MinIO** — Object Storage (port 9000)
- **Prometheus** — Metrics (port 9090)
- **Grafana** — Dashboards (port 3000)

## Commands

| Command | Description |
|---|---|
| `make setup` | One-time developer setup |
| `make dev` | Start development environment |
| `make build` | Compile all Go packages |
| `make test` | Run all tests |
| `make lint` | Run all linters |
| `make proto` | Compile protobuf definitions |
