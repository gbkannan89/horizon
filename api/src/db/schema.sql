-- PostgreSQL Schema for FinScore

CREATE TABLE IF NOT EXISTS households (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    invite_code VARCHAR(10) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    user_type VARCHAR(50) DEFAULT 'salaried' NOT NULL, -- 'salaried' or 'self-employed'
    risk_profile VARCHAR(50) DEFAULT 'moderate' NOT NULL, -- 'conservative', 'moderate', 'aggressive'
    household_id INTEGER REFERENCES households(id) ON DELETE SET NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS contributing_members (
    id SERIAL PRIMARY KEY,
    household_id INTEGER REFERENCES households(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    monthly_income NUMERIC(20, 2) DEFAULT 0.00 NOT NULL,
    contribution_to_household NUMERIC(20, 2) DEFAULT 0.00 NOT NULL,
    relationship VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS incomes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    label VARCHAR(255),                 -- user-friendly name e.g. 'My Salary', 'Partner Income'
    type VARCHAR(50) NOT NULL,          -- 'salary', 'business', 'passive'
    amount NUMERIC(20, 2) NOT NULL,
    frequency VARCHAR(50) NOT NULL,     -- 'monthly', 'annual'
    pf_employee NUMERIC(20, 2) DEFAULT 0.00 NOT NULL,
    pf_employer NUMERIC(20, 2) DEFAULT 0.00 NOT NULL,
    shares_deduction NUMERIC(20, 2) DEFAULT 0.00 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS assets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'bank', 'fd', 'rd', 'physical'
    subtype VARCHAR(50), -- 'property', 'gold', 'vehicle', etc.
    name VARCHAR(255) NOT NULL,
    amount NUMERIC(20, 2) NOT NULL,
    interest_rate NUMERIC(5, 2) DEFAULT 0.00 NOT NULL,
    start_date DATE,
    maturity_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS liabilities (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'home_loan', 'car_loan', 'personal_loan', 'credit_card'
    name VARCHAR(255) NOT NULL,
    outstanding NUMERIC(20, 2) NOT NULL,
    emi NUMERIC(20, 2) DEFAULT 0.00 NOT NULL,
    interest_rate NUMERIC(5, 2) DEFAULT 0.00 NOT NULL,
    tenure_months INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS recurring_bills (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount NUMERIC(20, 2) NOT NULL,
    frequency VARCHAR(50) NOT NULL, -- 'monthly', 'quarterly', 'half-yearly', 'yearly'
    category VARCHAR(100) NOT NULL, -- 'rent', 'utilities', 'insurance', 'education', etc.
    bucket VARCHAR(50) NOT NULL, -- 'Needs', 'Wants', 'Savings'
    due_day INTEGER DEFAULT 1 NOT NULL,
    monthly_equivalent NUMERIC(20, 2) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE NOT NULL,
    is_subscription BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS health_scores (
    id SERIAL PRIMARY KEY,
    household_id INTEGER REFERENCES households(id) ON DELETE CASCADE NOT NULL,
    score INTEGER NOT NULL,
    max_possible_score INTEGER NOT NULL,
    pillar_scores JSONB NOT NULL,
    adherence_history JSONB,
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
    token_hash VARCHAR(255) PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL
);

-- Schema updates / migrations
ALTER TABLE households ADD COLUMN IF NOT EXISTS invite_code VARCHAR(10) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS name VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS phone VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS user_type VARCHAR(50) DEFAULT 'salaried';
ALTER TABLE users ADD COLUMN IF NOT EXISTS risk_profile VARCHAR(50) DEFAULT 'moderate';
ALTER TABLE users ADD COLUMN IF NOT EXISTS household_id INTEGER REFERENCES households(id) ON DELETE SET NULL;
ALTER TABLE incomes ADD COLUMN IF NOT EXISTS label VARCHAR(255);
ALTER TABLE incomes ADD COLUMN IF NOT EXISTS pf_employee NUMERIC(20, 2) DEFAULT 0.00 NOT NULL;
ALTER TABLE incomes ADD COLUMN IF NOT EXISTS pf_employer NUMERIC(20, 2) DEFAULT 0.00 NOT NULL;
ALTER TABLE incomes ADD COLUMN IF NOT EXISTS shares_deduction NUMERIC(20, 2) DEFAULT 0.00 NOT NULL;

CREATE TABLE IF NOT EXISTS goals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    target_amount NUMERIC(20, 2) NOT NULL,
    current_amount NUMERIC(20, 2) DEFAULT 0.00 NOT NULL,
    status VARCHAR(50) DEFAULT 'On Track' NOT NULL, -- 'On Track', 'Behind', 'Ahead'
    color VARCHAR(20) DEFAULT '#059669',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS expenses (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount NUMERIC(20, 2) NOT NULL,
    category VARCHAR(100) NOT NULL, -- 'Rent', 'Food', 'SIP', etc.
    bucket VARCHAR(50) NOT NULL, -- 'Needs', 'Wants', 'Savings'
    icon VARCHAR(50),
    date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wishlist_items (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    amount NUMERIC(20, 2) NOT NULL,
    added_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    unlock_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(50) DEFAULT 'locked', -- 'locked', 'unlocked', 'bought_early', 'bought', 'discarded'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP);

-- Milestone 2: EMIs, Vehicles, Enhanced Assets, Wishlist custom locks
ALTER TABLE recurring_bills ADD COLUMN IF NOT EXISTS is_emi BOOLEAN DEFAULT FALSE;
ALTER TABLE recurring_bills ADD COLUMN IF NOT EXISTS emi_total_months INTEGER;
ALTER TABLE recurring_bills ADD COLUMN IF NOT EXISTS emi_months_paid INTEGER DEFAULT 0;
ALTER TABLE recurring_bills ADD COLUMN IF NOT EXISTS start_date DATE;

ALTER TABLE assets ADD COLUMN IF NOT EXISTS is_liability BOOLEAN DEFAULT FALSE;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS generates_income BOOLEAN DEFAULT FALSE;
ALTER TABLE assets ADD COLUMN IF NOT EXISTS income_frequency VARCHAR(50); -- 'monthly', 'quarterly', 'yearly'

ALTER TABLE wishlist_items ADD COLUMN IF NOT EXISTS lock_duration_days INTEGER DEFAULT 30;

ALTER TABLE assets ADD COLUMN IF NOT EXISTS purchase_price NUMERIC(20, 2);
ALTER TABLE assets ADD COLUMN IF NOT EXISTS purchase_date DATE;

CREATE TABLE IF NOT EXISTS vehicles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    make_model VARCHAR(255) NOT NULL,
    purchase_cost NUMERIC(20, 2) NOT NULL,
    insurance_renewal_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS net_worth_snapshots (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    snapshot_date DATE NOT NULL,
    net_worth NUMERIC(20, 2) NOT NULL,
    UNIQUE(user_id, snapshot_date)
);

-- Advisor Engine: Simulations, Nudges, Subscriptions
CREATE TABLE IF NOT EXISTS simulations (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    scenario_name VARCHAR(255) NOT NULL,
    scenario_type VARCHAR(50) NOT NULL,
    input_params JSONB NOT NULL,
    results JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS nudges (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    category VARCHAR(50) NOT NULL,
    severity VARCHAR(20) NOT NULL DEFAULT 'info',
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    action_label VARCHAR(100),
    action_link VARCHAR(255),
    is_dismissed BOOLEAN DEFAULT FALSE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS subscription_insights (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    bill_id INTEGER REFERENCES recurring_bills(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    last_used_date DATE,
    monthly_cost NUMERIC(20, 2) NOT NULL,
    annual_cost NUMERIC(20, 2) NOT NULL,
    savings_opportunity NUMERIC(20, 2) DEFAULT 0.00,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS insurances (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'health', 'term', 'life', 'vehicle', 'other'
    provider VARCHAR(255) NOT NULL,
    policy_name VARCHAR(255),
    premium_amount NUMERIC(20, 2) NOT NULL,
    premium_frequency VARCHAR(50) NOT NULL, -- 'monthly', 'quarterly', 'yearly'
    coverage_amount NUMERIC(20, 2),
    renewal_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Collection Tracking (group payments from multiple people)
CREATE TABLE IF NOT EXISTS collections (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    label VARCHAR(255) NOT NULL,
    description TEXT,
    total_expected NUMERIC(20, 2) NOT NULL DEFAULT 0,
    total_collected NUMERIC(20, 2) NOT NULL DEFAULT 0,
    status VARCHAR(50) DEFAULT 'active', -- 'active', 'settled', 'cancelled'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collection_members (
    id SERIAL PRIMARY KEY,
    collection_id INTEGER REFERENCES collections(id) ON DELETE CASCADE NOT NULL,
    name VARCHAR(255) NOT NULL,
    expected_amount NUMERIC(20, 2) NOT NULL,
    paid_amount NUMERIC(20, 2) DEFAULT 0.00,
    paid_date DATE,
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'partial', 'paid'
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Family Invites: email-based invitations to join a household
CREATE TABLE IF NOT EXISTS family_invites (
    id SERIAL PRIMARY KEY,
    household_id INTEGER REFERENCES households(id) ON DELETE CASCADE NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    relationship VARCHAR(100),
    invite_code VARCHAR(10) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(household_id, email)
);

-- ═══════════════════════════════════════════════════════════════════════════
-- Phase 1A: Gold Assets
-- ═══════════════════════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS gold_assets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    carat INTEGER NOT NULL CHECK (carat IN (18, 22, 24)),
    grams NUMERIC(10, 2) NOT NULL CHECK (grams > 0),
    purchase_price_per_gram NUMERIC(20, 2) NOT NULL,
    purchase_date DATE NOT NULL,
    last_current_value NUMERIC(20, 2),
    last_value_updated_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ═══════════════════════════════════════════════════════════════════════════
-- Phase 1B: Stock/MF Holdings
-- ═══════════════════════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS stock_holdings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    ticker VARCHAR(50) NOT NULL,
    exchange VARCHAR(10) DEFAULT 'NSE',
    name VARCHAR(255),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    avg_purchase_price NUMERIC(20, 2) NOT NULL,
    purchase_date DATE NOT NULL,
    total_invested NUMERIC(20, 2) NOT NULL,
    last_current_price NUMERIC(20, 2),
    last_current_value NUMERIC(20, 2),
    last_value_updated_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ═══════════════════════════════════════════════════════════════════════════
-- Phase 1C: FD Maturity Enhancement
-- ═══════════════════════════════════════════════════════════════════════════
ALTER TABLE assets ADD COLUMN IF NOT EXISTS years_of_deposit INTEGER;

-- ═══════════════════════════════════════════════════════════════════════════
-- Phase 1D: PF Assets
-- ═══════════════════════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS pf_assets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    start_date DATE NOT NULL,
    retirement_age INTEGER NOT NULL DEFAULT 58,
    current_balance NUMERIC(20, 2) DEFAULT 0.00 NOT NULL,
    monthly_contribution NUMERIC(20, 2) DEFAULT 0.00 NOT NULL,
    employer_contribution NUMERIC(20, 2) DEFAULT 0.00,
    interest_rate NUMERIC(5, 2) DEFAULT 8.25 NOT NULL,
    user_date_of_birth DATE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- ═══════════════════════════════════════════════════════════════════════════
-- Phase 1E: Gold Loan (extend liabilities)
-- ═══════════════════════════════════════════════════════════════════════════
ALTER TABLE liabilities ADD COLUMN IF NOT EXISTS gold_grams NUMERIC(10, 2);
ALTER TABLE liabilities ADD COLUMN IF NOT EXISTS gold_carat INTEGER;
ALTER TABLE liabilities ADD COLUMN IF NOT EXISTS gold_items_count INTEGER;
ALTER TABLE liabilities ADD COLUMN IF NOT EXISTS loan_date DATE;
ALTER TABLE liabilities ADD COLUMN IF NOT EXISTS bank_name VARCHAR(255);

-- ═══════════════════════════════════════════════════════════════════════════
-- Phase 1F: Lend/Borrow Records
-- ═══════════════════════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS lending_records (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    direction VARCHAR(10) NOT NULL CHECK (direction IN ('lent', 'borrowed')),
    person_name VARCHAR(255) NOT NULL,
    amount NUMERIC(20, 2) NOT NULL,
    date_given DATE NOT NULL,
    promised_return_date DATE,
    actual_return_date DATE,
    returned_amount NUMERIC(20, 2) DEFAULT 0.00 NOT NULL,
    status VARCHAR(20) DEFAULT 'open' CHECK (status IN ('open', 'partial', 'closed', 'overdue')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
