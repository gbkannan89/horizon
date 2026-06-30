-- V004__institution.sql
-- Purpose: Institution domain tables
-- Domain: HZN-DOM-005

CREATE TABLE institutions (
    institution_id      UUID PRIMARY KEY,
    name                VARCHAR(300) NOT NULL,
    institution_type    VARCHAR(50) NOT NULL,
    category            VARCHAR(50) NOT NULL,
    status              VARCHAR(50) NOT NULL DEFAULT 'Draft',
    country             VARCHAR(2) NOT NULL,
    trust_level         VARCHAR(50) NOT NULL DEFAULT 'Unverified',
    health              VARCHAR(50) NOT NULL DEFAULT 'Stable',
    confidence          VARCHAR(50) NOT NULL DEFAULT 'UserDefined',
    products            JSONB NOT NULL DEFAULT '[]',
    connectivity        JSONB NOT NULL DEFAULT '["Manual"]',
    website             TEXT DEFAULT '',
    phone               VARCHAR(50) DEFAULT '',
    email               VARCHAR(255) DEFAULT '',
    address             TEXT DEFAULT '',
    headquarters        TEXT DEFAULT '',
    regulator           VARCHAR(255) DEFAULT '',
    regulatory_license  VARCHAR(255) DEFAULT '',
    metadata            JSONB DEFAULT '{}',
    tags                TEXT[] DEFAULT '{}',
    notes               TEXT DEFAULT '',
    source_of_truth     VARCHAR(50) NOT NULL DEFAULT 'User',
    version             INT NOT NULL DEFAULT 1,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_institutions_type ON institutions (institution_type);
CREATE INDEX idx_institutions_country ON institutions (country);
CREATE INDEX idx_institutions_status ON institutions (status);
CREATE INDEX idx_institutions_name ON institutions USING gin (name gin_trgm_ops);
