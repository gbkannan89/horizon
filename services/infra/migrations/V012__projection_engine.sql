-- V012__projection_engine.sql
-- Purpose: Projection Engine storage tables
-- Engine: HZN-ENG-001

CREATE TABLE projection_outputs (
    output_id       VARCHAR(255) PRIMARY KEY,
    projection_type VARCHAR(50) NOT NULL,
    version         INT NOT NULL DEFAULT 1,
    status          VARCHAR(50) NOT NULL DEFAULT 'Completed',
    output_data     JSONB NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projection_type ON projection_outputs (projection_type, created_at DESC);
