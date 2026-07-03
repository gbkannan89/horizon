-- V019__budget.sql
-- Budget domain tables

CREATE TABLE IF NOT EXISTS budgets (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    household_id TEXT NOT NULL DEFAULT '',
    name TEXT NOT NULL,
    period TEXT NOT NULL DEFAULT 'Monthly',
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status TEXT NOT NULL DEFAULT 'Draft',
    total_budgeted BIGINT NOT NULL DEFAULT 0,
    total_spent BIGINT NOT NULL DEFAULT 0,
    total_remaining BIGINT NOT NULL DEFAULT 0,
    currency TEXT NOT NULL DEFAULT 'INR',
    categories JSONB NOT NULL DEFAULT '[]'::jsonb,
    tags JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_budgets_user ON budgets(user_id);
CREATE INDEX idx_budgets_status ON budgets(status);
CREATE INDEX idx_budgets_period ON budgets(start_date, end_date);
