-- V017__simulation_engine.sql
-- Purpose: Simulation Engine storage tables
-- Engine: HZN-ENG-003

CREATE TABLE simulations (
    scenario_id   VARCHAR(255) PRIMARY KEY,
    sim_type      VARCHAR(50) NOT NULL,
    status        VARCHAR(50) NOT NULL DEFAULT 'Completed',
    sim_data      JSONB NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_simulations_type ON simulations (sim_type, created_at DESC);
