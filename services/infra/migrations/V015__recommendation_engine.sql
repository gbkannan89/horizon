-- V015__recommendation_engine.sql
-- Purpose: Recommendation Engine storage tables
-- Engine: HZN-ENG-002

CREATE TABLE recommendations (
    recommendation_id VARCHAR(255) PRIMARY KEY,
    category         VARCHAR(50) NOT NULL,
    priority         INT NOT NULL DEFAULT 0,
    score            NUMERIC(5,2) NOT NULL DEFAULT 0,
    status           VARCHAR(50) NOT NULL DEFAULT 'Generated',
    rec_data         JSONB NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_rec_status ON recommendations (status, priority);
