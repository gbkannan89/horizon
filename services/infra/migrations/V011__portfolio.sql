-- V011__portfolio.sql
-- Purpose: Portfolio domain tables
-- Domain: HZN-DOM-008

CREATE TABLE portfolios (
    portfolio_id        UUID PRIMARY KEY,
    name                VARCHAR(200) NOT NULL,
    portfolio_type      VARCHAR(50) NOT NULL,
    owner_id            UUID NOT NULL,
    household_id        UUID,
    base_currency       VARCHAR(3) NOT NULL,
    membership_model    VARCHAR(50) NOT NULL DEFAULT 'Manual',
    status              VARCHAR(50) NOT NULL DEFAULT 'Draft',
    portfolio_health    VARCHAR(50) NOT NULL DEFAULT 'Stable',
    portfolio_confidence VARCHAR(50) NOT NULL DEFAULT 'UserVerified',
    risk_profile        VARCHAR(50) DEFAULT 'Balanced',
    benchmark_profile   VARCHAR(50) DEFAULT 'None',
    benchmark           VARCHAR(100) DEFAULT '',
    liquidity_profile   VARCHAR(50) DEFAULT 'MediumTerm',
    source_of_truth     VARCHAR(50) NOT NULL DEFAULT 'User',
    tags                TEXT[] DEFAULT '{}',
    metadata            JSONB DEFAULT '{}',
    notes               TEXT DEFAULT '',
    version             INT NOT NULL DEFAULT 1,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE portfolio_members (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    portfolio_id    UUID NOT NULL REFERENCES portfolios(portfolio_id) ON DELETE CASCADE,
    entity_id       VARCHAR(255) NOT NULL,
    entity_type     VARCHAR(50) NOT NULL,
    weight          NUMERIC(5,2),
    UNIQUE (portfolio_id, entity_id, entity_type)
);

CREATE INDEX idx_portfolios_owner ON portfolios (owner_id);
CREATE INDEX idx_portfolios_type ON portfolios (owner_id, portfolio_type);
CREATE INDEX idx_portfolios_status ON portfolios (status);
CREATE INDEX idx_pf_members_portfolio ON portfolio_members (portfolio_id);
