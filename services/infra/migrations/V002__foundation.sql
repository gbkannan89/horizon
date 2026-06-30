-- V002__foundation.sql
-- Purpose: Shared foundation tables (outbox, event store, idempotency)

CREATE TABLE IF NOT EXISTS outbox (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type      VARCHAR(255) NOT NULL,
    event_payload   JSONB NOT NULL,
    idempotency_key VARCHAR(255) NOT NULL,
    correlation_id  VARCHAR(255) NOT NULL,
    partition_key   VARCHAR(255) NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    retry_count     INT NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_attempt_at TIMESTAMPTZ,
    UNIQUE (idempotency_key)
);

CREATE INDEX idx_outbox_status ON outbox (status, created_at) WHERE status = 'PENDING';

CREATE TABLE IF NOT EXISTS event_store (
    event_id        UUID PRIMARY KEY,
    event_type      VARCHAR(255) NOT NULL,
    aggregate_type  VARCHAR(100) NOT NULL,
    aggregate_id    VARCHAR(255) NOT NULL,
    event_data      JSONB NOT NULL,
    version         INT NOT NULL DEFAULT 1,
    correlation_id  VARCHAR(255),
    user_id         VARCHAR(255),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_event_store_aggregate ON event_store (aggregate_type, aggregate_id, version);
CREATE INDEX idx_event_store_created ON event_store (created_at DESC);
CREATE INDEX idx_event_store_type ON event_store (event_type);
