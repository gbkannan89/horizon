-- V003__user.sql
-- Purpose: User domain tables
-- Domain: HZN-DOM-010

CREATE TABLE users (
    user_id                     UUID PRIMARY KEY,
    display_name                VARCHAR(100) NOT NULL,
    legal_name                  VARCHAR(255) NOT NULL,
    preferred_name              VARCHAR(100),
    date_of_birth               DATE,
    country                     VARCHAR(2) NOT NULL,
    base_currency               VARCHAR(3) NOT NULL,
    locale                      VARCHAR(10) NOT NULL,
    timezone                    VARCHAR(50) NOT NULL,
    user_health                 VARCHAR(50) NOT NULL DEFAULT 'Healthy',
    user_confidence             VARCHAR(50) NOT NULL DEFAULT 'UserVerified',
    financial_identity_profile  VARCHAR(50) NOT NULL DEFAULT 'Personal',
    status                      VARCHAR(50) NOT NULL DEFAULT 'Registered',
    user_type                   VARCHAR(50) NOT NULL DEFAULT 'Individual',
    financial_behaviour_profile VARCHAR(50) NOT NULL DEFAULT 'Balanced',
    profile_completeness        NUMERIC(5,2) NOT NULL DEFAULT 0,
    source_of_truth             VARCHAR(50) NOT NULL DEFAULT 'System',
    tags                        TEXT[] DEFAULT '{}',
    metadata                    JSONB DEFAULT '{}',
    notes                       TEXT DEFAULT '',
    version                     INT NOT NULL DEFAULT 1,
    created_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE user_consents (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    consent_type    VARCHAR(50) NOT NULL,
    granted         BOOLEAN NOT NULL DEFAULT TRUE,
    granted_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at      TIMESTAMPTZ,
    scope           TEXT NOT NULL DEFAULT ''
);

CREATE TABLE user_preferences (
    user_id         UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    preference_key  VARCHAR(100) NOT NULL,
    preference_value JSONB NOT NULL,
    PRIMARY KEY (user_id, preference_key)
);

CREATE TABLE user_privacy (
    user_id             UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    data_sharing        VARCHAR(50) NOT NULL DEFAULT 'OptIn',
    household_visibility VARCHAR(50) NOT NULL DEFAULT 'RoleBased',
    third_party_access  BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE user_households (
    user_id         UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    household_id    VARCHAR(255) NOT NULL,
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, household_id)
);

CREATE INDEX idx_users_status ON users (status);
CREATE INDEX idx_users_type ON users (user_type);
CREATE INDEX idx_users_health ON users (user_health);
CREATE INDEX idx_users_display_name ON users USING gin (display_name gin_trgm_ops);
CREATE INDEX idx_user_consents_type ON user_consents (user_id, consent_type);
