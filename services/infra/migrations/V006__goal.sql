-- V006__goal.sql
-- Purpose: Goal domain tables
-- Domain: HZN-DOM-001

CREATE TABLE goals (
    goal_id             UUID PRIMARY KEY,
    user_id             UUID NOT NULL,
    name                VARCHAR(200) NOT NULL,
    importance          VARCHAR(50) NOT NULL,
    type                VARCHAR(50) NOT NULL,
    subtype             VARCHAR(50),
    success_criteria    JSONB NOT NULL,
    target_date         TIMESTAMPTZ,
    priority            INT NOT NULL,
    contribution_schedule JSONB,
    risk_tolerance      VARCHAR(50),
    notes               TEXT DEFAULT '',
    parent_goal_id      UUID,
    tags                TEXT[] DEFAULT '{}',
    status              VARCHAR(50) NOT NULL DEFAULT 'Draft',
    version             INT NOT NULL DEFAULT 1,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_parent_goal FOREIGN KEY (parent_goal_id) REFERENCES goals(goal_id),
    CONSTRAINT uq_goal_name_type UNIQUE (user_id, name, type),
    CONSTRAINT chk_priority CHECK (priority >= 1)
);

CREATE UNIQUE INDEX uq_goal_active_priority ON goals (user_id, priority) WHERE status = 'Active';

CREATE INDEX idx_goals_status ON goals (user_id, status);
CREATE INDEX idx_goals_importance ON goals (user_id, importance);
CREATE INDEX idx_goals_type ON goals (user_id, type);
CREATE INDEX idx_goals_target_date ON goals (target_date) WHERE target_date IS NOT NULL;
