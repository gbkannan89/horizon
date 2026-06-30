-- V013__risk_engine.sql
-- Purpose: Risk Engine storage tables
-- Engine: HZN-ENG-006

CREATE TABLE risk_assessments (
    assessment_id   VARCHAR(255) PRIMARY KEY,
    composite_score INT NOT NULL DEFAULT 0,
    risk_level      VARCHAR(50) NOT NULL DEFAULT 'Minimal',
    output_data     JSONB NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_risk_created ON risk_assessments (created_at DESC);
