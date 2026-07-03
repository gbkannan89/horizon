-- V021__household_shared_entities.sql
-- Household domain tables expansion for shared entities

ALTER TABLE households 
ADD COLUMN IF NOT EXISTS linked_accounts JSONB NOT NULL DEFAULT '[]'::jsonb,
ADD COLUMN IF NOT EXISTS linked_goals JSONB NOT NULL DEFAULT '[]'::jsonb,
ADD COLUMN IF NOT EXISTS linked_budgets JSONB NOT NULL DEFAULT '[]'::jsonb,
ADD COLUMN IF NOT EXISTS goal_contributions JSONB NOT NULL DEFAULT '[]'::jsonb;
