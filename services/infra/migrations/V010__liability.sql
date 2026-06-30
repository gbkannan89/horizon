-- V010__liability.sql
-- Purpose: Liability domain tables
-- Domain: HZN-DOM-007

CREATE TABLE liabilities (
    liability_id            UUID PRIMARY KEY,
    name                    VARCHAR(200) NOT NULL,
    classification          VARCHAR(50) NOT NULL,
    liability_type          VARCHAR(100),
    owner_id                UUID NOT NULL,
    household_id            UUID,
    institution_id          UUID,
    servicing_account_id    UUID,
    currency                VARCHAR(3) NOT NULL,
    original_principal      NUMERIC(19,4) NOT NULL,
    interest_rate           NUMERIC(6,4) NOT NULL DEFAULT 0,
    interest_model          VARCHAR(50) NOT NULL DEFAULT 'Simple',
    repayment_profile       VARCHAR(50) NOT NULL DEFAULT 'FixedInstallment',
    repayment_model         VARCHAR(50) NOT NULL DEFAULT 'FixedEMI',
    installment_amount      NUMERIC(19,4) NOT NULL DEFAULT 0,
    installment_frequency   VARCHAR(50) DEFAULT 'Monthly',
    remaining_installments  INT NOT NULL DEFAULT 0,
    maturity_date           DATE NOT NULL,
    collateral              TEXT DEFAULT '',
    liability_health        VARCHAR(50) NOT NULL DEFAULT 'Stable',
    liability_confidence    VARCHAR(50) NOT NULL DEFAULT 'UserVerified',
    status                  VARCHAR(50) NOT NULL DEFAULT 'Planned',
    source_of_truth         VARCHAR(50) NOT NULL DEFAULT 'User',
    external_reference      VARCHAR(255),
    tags                    TEXT[] DEFAULT '{}',
    metadata                JSONB DEFAULT '{}',
    notes                   TEXT DEFAULT '',
    version                 INT NOT NULL DEFAULT 1,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_liability_institution FOREIGN KEY (institution_id) REFERENCES institutions(institution_id),
    CONSTRAINT fk_liability_account FOREIGN KEY (servicing_account_id) REFERENCES accounts(account_id)
);

CREATE INDEX idx_liabilities_owner ON liabilities (owner_id);
CREATE INDEX idx_liabilities_classification ON liabilities (owner_id, classification);
CREATE INDEX idx_liabilities_status ON liabilities (owner_id, status);
CREATE INDEX idx_liabilities_institution ON liabilities (institution_id);
CREATE INDEX idx_liabilities_maturity ON liabilities (maturity_date) WHERE status IN ('Planned', 'Approved', 'Active');
