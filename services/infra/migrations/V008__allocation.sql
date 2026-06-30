-- V008__allocation.sql
-- Purpose: Allocation domain tables
-- Domain: HZN-DOM-004

CREATE TABLE allocations (
    allocation_id       UUID PRIMARY KEY,
    goal_id             UUID NOT NULL,
    funding_source_id   VARCHAR(255) NOT NULL,
    funding_source_type VARCHAR(50) NOT NULL DEFAULT 'Account',
    allocation_type     VARCHAR(50) NOT NULL,
    allocation_strategy VARCHAR(50) NOT NULL DEFAULT 'GoalPriority',
    allocation_health   VARCHAR(50) NOT NULL DEFAULT 'Healthy',
    allocation_confidence VARCHAR(50) NOT NULL DEFAULT 'Confirmed',
    funding_commitment  NUMERIC(19,4) NOT NULL DEFAULT 0,
    priority            INT NOT NULL,
    weight              NUMERIC(5,2),
    fixed_amount        NUMERIC(19,4),
    currency            VARCHAR(3) NOT NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'Draft',
    effective_date      TIMESTAMPTZ NOT NULL,
    expiration_date     TIMESTAMPTZ,
    reserved_amount     NUMERIC(19,4) NOT NULL DEFAULT 0,
    allocated_amount    NUMERIC(19,4) NOT NULL DEFAULT 0,
    source_of_truth     VARCHAR(50) NOT NULL DEFAULT 'User',
    created_by          VARCHAR(50) NOT NULL DEFAULT 'User',
    approved_by         VARCHAR(255),
    tags                TEXT[] DEFAULT '{}',
    metadata            JSONB DEFAULT '{}',
    notes               TEXT DEFAULT '',
    version             INT NOT NULL DEFAULT 1,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_allocation_goal FOREIGN KEY (goal_id) REFERENCES goals(goal_id),
    CONSTRAINT chk_priority CHECK (priority >= 1)
);

CREATE INDEX idx_allocations_goal ON allocations (goal_id);
CREATE INDEX idx_allocations_source ON allocations (funding_source_id);
CREATE INDEX idx_allocations_status ON allocations (status);
CREATE INDEX idx_allocations_goal_source ON allocations (goal_id, funding_source_id);
