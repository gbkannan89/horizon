-- V007__account.sql
-- Purpose: Account domain tables
-- Domain: HZN-DOM-003

CREATE TABLE accounts (
    account_id          UUID PRIMARY KEY,
    institution_id      UUID,
    owner_id            UUID NOT NULL,
    household_id        UUID,
    classification      VARCHAR(50) NOT NULL,
    account_type        VARCHAR(50) NOT NULL,
    account_sub_type    VARCHAR(100),
    account_name        VARCHAR(200) NOT NULL,
    currency            VARCHAR(3) NOT NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'Draft',
    opened_date         DATE NOT NULL,
    closed_date         DATE,
    liquidity_profile   VARCHAR(50) NOT NULL DEFAULT 'MediumTerm',
    account_health      VARCHAR(50) NOT NULL DEFAULT 'Healthy',
    credit_limit        NUMERIC(19,4),
    interest_rate       NUMERIC(6,4),
    country             VARCHAR(2),
    source_of_truth     VARCHAR(50),
    visibility          VARCHAR(50) NOT NULL DEFAULT 'Private',
    external_reference  VARCHAR(255),
    tags                TEXT[] DEFAULT '{}',
    metadata            JSONB DEFAULT '{}',
    notes               TEXT DEFAULT '',
    version             INT NOT NULL DEFAULT 1,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_account_institution FOREIGN KEY (institution_id) REFERENCES institutions(institution_id)
);

CREATE INDEX idx_accounts_owner ON accounts (owner_id);
CREATE INDEX idx_accounts_type ON accounts (owner_id, account_type);
CREATE INDEX idx_accounts_status ON accounts (owner_id, status);
CREATE INDEX idx_accounts_institution ON accounts (institution_id);
CREATE INDEX idx_accounts_name ON accounts USING gin (account_name gin_trgm_ops);
