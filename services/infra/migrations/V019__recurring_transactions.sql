CREATE TABLE recurring_transactions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    household_id TEXT,
    name TEXT NOT NULL,
    description TEXT,
    amount BIGINT NOT NULL,
    currency TEXT NOT NULL,
    frequency TEXT NOT NULL,
    interval INTEGER NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE,
    next_occurrence TIMESTAMP WITH TIME ZONE,
    last_occurrence TIMESTAMP WITH TIME ZONE,
    status TEXT NOT NULL,
    event_template JSONB NOT NULL,
    skip_holidays BOOLEAN NOT NULL DEFAULT false,
    skip_weekends BOOLEAN NOT NULL DEFAULT false,
    tags JSONB,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_recurring_user ON recurring_transactions(user_id);
CREATE INDEX idx_recurring_due ON recurring_transactions(next_occurrence) WHERE status = 'Active';
