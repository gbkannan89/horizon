-- V009__asset.sql
-- Purpose: Asset domain tables
-- Domain: HZN-DOM-006

CREATE TABLE assets (
    asset_id            UUID PRIMARY KEY,
    asset_name          VARCHAR(200) NOT NULL,
    classification      VARCHAR(50) NOT NULL,
    asset_type          VARCHAR(100),
    ownership_model     VARCHAR(50) NOT NULL DEFAULT 'Individual',
    owner_id            UUID NOT NULL,
    household_id        UUID,
    institution_id      UUID,
    account_id          UUID,
    currency            VARCHAR(3) NOT NULL,
    acquisition_date    DATE,
    disposition_date    DATE,
    asset_health        VARCHAR(50) NOT NULL DEFAULT 'Stable',
    asset_confidence    VARCHAR(50) NOT NULL DEFAULT 'UserVerified',
    cost_basis          NUMERIC(19,4) NOT NULL DEFAULT 0,
    valuation_profile   VARCHAR(50) NOT NULL DEFAULT 'ManualAssessment',
    valuation_method    VARCHAR(50) NOT NULL DEFAULT 'ManualValue',
    valuation_date      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    quantity            NUMERIC(19,4),
    unit_price          NUMERIC(19,4) NOT NULL DEFAULT 0,
    liquidity_profile   VARCHAR(50) NOT NULL DEFAULT 'MediumTerm',
    ownership_percentage NUMERIC(5,2) NOT NULL DEFAULT 100.00,
    status              VARCHAR(50) NOT NULL DEFAULT 'Planned',
    external_identifiers JSONB DEFAULT '{}',
    source_of_truth     VARCHAR(50) NOT NULL DEFAULT 'Manual',
    tags                TEXT[] DEFAULT '{}',
    metadata            JSONB DEFAULT '{}',
    notes               TEXT DEFAULT '',
    version             INT NOT NULL DEFAULT 1,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_asset_institution FOREIGN KEY (institution_id) REFERENCES institutions(institution_id),
    CONSTRAINT fk_asset_account FOREIGN KEY (account_id) REFERENCES accounts(account_id)
);

CREATE INDEX idx_assets_owner ON assets (owner_id);
CREATE INDEX idx_assets_classification ON assets (owner_id, classification);
CREATE INDEX idx_assets_status ON assets (owner_id, status);
CREATE INDEX idx_assets_institution ON assets (institution_id);
CREATE INDEX idx_assets_account ON assets (account_id);

CREATE TABLE valuation_history (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    asset_id        UUID NOT NULL REFERENCES assets(asset_id) ON DELETE CASCADE,
    value           NUMERIC(19,4) NOT NULL,
    method          VARCHAR(50) NOT NULL,
    valuation_date  DATE NOT NULL,
    source_of_truth VARCHAR(50) DEFAULT 'User',
    recorded_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_val_history_asset ON valuation_history (asset_id, valuation_date DESC);
