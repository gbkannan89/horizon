-- V022__account_household_index.sql
-- Add index for querying accounts by household_id

CREATE INDEX IF NOT EXISTS idx_accounts_household ON accounts (household_id);
