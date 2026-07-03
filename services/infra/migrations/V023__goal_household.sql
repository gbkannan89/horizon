-- V023__goal_household.sql
-- Add household_id to goals table

ALTER TABLE goals ADD COLUMN household_id UUID;
CREATE INDEX IF NOT EXISTS idx_goals_household ON goals (household_id);
