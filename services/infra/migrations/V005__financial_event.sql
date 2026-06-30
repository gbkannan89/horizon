-- V005__financial_event.sql
-- Purpose: Financial Event domain tables (root event stream)
-- Domain: HZN-DOM-002

CREATE TABLE financial_events (
    event_id            UUID PRIMARY KEY,
    user_id             UUID NOT NULL,
    household_id        UUID,
    event_type          VARCHAR(50) NOT NULL,
    event_sub_type      VARCHAR(100),
    event_date          TIMESTAMPTZ NOT NULL,
    effective_date      TIMESTAMPTZ NOT NULL,
    currency            VARCHAR(3) NOT NULL,
    amount              NUMERIC(19,4) NOT NULL,
    source              TEXT DEFAULT '',
    destination         TEXT DEFAULT '',
    description         TEXT NOT NULL DEFAULT '',
    notes               TEXT DEFAULT '',
    reference           TEXT DEFAULT '',
    attachments         TEXT[] DEFAULT '{}',
    origin              VARCHAR(50) NOT NULL,
    confidence          VARCHAR(50) NOT NULL DEFAULT 'Confirmed',
    source_of_truth     VARCHAR(50),
    created_by          VARCHAR(50) NOT NULL DEFAULT 'User',
    imported_from       VARCHAR(255),
    correlation_id      VARCHAR(255),
    reversal_of_event_id UUID,
    tags                TEXT[] DEFAULT '{}',
    state               VARCHAR(50) NOT NULL DEFAULT 'DRAFT',
    sequence_number     BIGINT NOT NULL DEFAULT 0,
    order_index         VARCHAR(100),
    version             INT NOT NULL DEFAULT 1,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_reversal FOREIGN KEY (reversal_of_event_id) REFERENCES financial_events(event_id)
);

CREATE INDEX idx_fe_user_date ON financial_events (user_id, effective_date DESC, event_date DESC, sequence_number DESC);
CREATE INDEX idx_fe_user_type ON financial_events (user_id, event_type);
CREATE INDEX idx_fe_state ON financial_events (state);
CREATE INDEX idx_fe_source ON financial_events (source);
CREATE INDEX idx_fe_destination ON financial_events (destination);
CREATE INDEX idx_fe_correlation ON financial_events (correlation_id);
CREATE INDEX idx_fe_reference ON financial_events (source_of_truth, reference) WHERE source_of_truth IS NOT NULL AND reference != '';
CREATE INDEX idx_fe_imported ON financial_events (imported_from, reference) WHERE imported_from IS NOT NULL AND reference != '';
