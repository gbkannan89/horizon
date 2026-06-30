# Architecture Overview

Horizon is a Goal-Centric Personal Financial Operating System.

## Architecture

```
Flutter (Mobile/Desktop)
    ↓
Go API Gateway
    ↓
Go Application Services
    ↓
Go Domains (10)  │  Go Engines (6)
    ↓
PostgreSQL  │  Redis  │  NATS JetStream
    ↓
MinIO / S3
```

## Layers

- **Experience Layer** — Flutter applications
- **Application Layer** — Go services, CQRS orchestration
- **Domain Layer** — Go domain services, business logic
- **Engine Layer** — Go computation services, deterministic
- **Infrastructure Layer** — PostgreSQL, Redis, NATS, MinIO

## Key Principles

- Event-driven architecture
- CQRS — Commands and queries separated
- Deterministic engines
- AI never computes financial values

## Documentation

See specifications in `../horizon-spec/` for detailed documentation.
