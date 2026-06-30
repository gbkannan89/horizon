-- V014__health_score_engine.sql
-- Purpose: Health Score Engine storage tables
-- Engine: HZN-ENG-004

CREATE TABLE health_scores (
    score_id        VARCHAR(255) PRIMARY KEY,
    overall_score   INT NOT NULL DEFAULT 0,
    score_grade     VARCHAR(50) NOT NULL DEFAULT 'Fair',
    score_data      JSONB NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_health_scores_created ON health_scores (created_at DESC);
