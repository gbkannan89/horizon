-- V018__household.sql
-- Household domain tables

CREATE TABLE IF NOT EXISTS households (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    household_type TEXT NOT NULL DEFAULT 'Family',
    head_of_household_id TEXT NOT NULL,
    members JSONB NOT NULL DEFAULT '[]'::jsonb,
    status TEXT NOT NULL DEFAULT 'Draft',
    currency TEXT NOT NULL DEFAULT 'INR',
    country TEXT NOT NULL DEFAULT '',
    total_assets BIGINT NOT NULL DEFAULT 0,
    total_liabilities BIGINT NOT NULL DEFAULT 0,
    total_net_worth BIGINT NOT NULL DEFAULT 0,
    health TEXT NOT NULL DEFAULT 'Healthy',
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_households_status ON households(status);
CREATE INDEX idx_households_head ON households(head_of_household_id);
CREATE INDEX idx_households_members ON households USING GIN(members);
