-- V016__optimization_engine.sql
-- Purpose: Optimization Engine storage tables
-- Engine: HZN-ENG-005

CREATE TABLE optimizations (
    optimization_id     VARCHAR(255) PRIMARY KEY,
    status              VARCHAR(50) NOT NULL DEFAULT 'Completed',
    candidates_evaluated INT NOT NULL DEFAULT 0,
    candidates_valid    INT NOT NULL DEFAULT 0,
    opt_data            JSONB NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_optimizations_created ON optimizations (created_at DESC);
