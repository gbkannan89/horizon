-- S001__seed_demo_data.sql
-- Purpose: Seed comprehensive demo data for the Horizon application
-- Login: demo@horizon.app / password123
-- Note: The MVP login handler accepts any valid email/password (stub).
--       This seed provides realistic data for when real providers are wired.

-- ============================================================
-- DEMO USER
-- ============================================================
-- Using a consistent UUID for the demo user.
-- The login handler generates user_id as "user-" + email, but we use
-- a proper UUID here for when real user lookup is implemented.
-- Login email: demo@horizon.app  |  Password: password123

INSERT INTO users (user_id, display_name, legal_name, preferred_name, country, base_currency, locale, timezone, status, user_type, profile_completeness)
VALUES (
    'a1b2c3d4-0000-4000-8000-000000000001',
    'Priya Sharma', 'Priya Sharma', 'Priya',
    'IN', 'INR', 'en-IN', 'Asia/Kolkata',
    'Active', 'Individual', 85.00
);

INSERT INTO user_preferences (user_id, preference_key, preference_value)
VALUES
    ('a1b2c3d4-0000-4000-8000-000000000001', 'theme', '"light"'),
    ('a1b2c3d4-0000-4000-8000-000000000001', 'locale', '"en-IN"'),
    ('a1b2c3d4-0000-4000-8000-000000000001', 'currency', '"INR"');

INSERT INTO user_privacy (user_id, data_sharing, household_visibility, third_party_access)
VALUES ('a1b2c3d4-0000-4000-8000-000000000001', 'OptIn', 'RoleBased', FALSE);

-- ============================================================
-- INSTITUTIONS
-- ============================================================
INSERT INTO institutions (institution_id, name, institution_type, category, country, status, health, confidence, source_of_truth)
VALUES
    ('b1c2d3e4-0000-4000-8000-000000000001', 'HDFC Bank', 'Bank', 'Commercial', 'IN', 'Active', 'Healthy', 'UserVerified', 'User'),
    ('b1c2d3e4-0000-4000-8000-000000000002', 'State Bank of India', 'Bank', 'Commercial', 'IN', 'Active', 'Healthy', 'UserVerified', 'User'),
    ('b1c2d3e4-0000-4000-8000-000000000003', 'Zerodha', 'Broker', 'Brokerage', 'IN', 'Active', 'Healthy', 'UserVerified', 'User');

-- ============================================================
-- ACCOUNTS
-- ============================================================
INSERT INTO accounts (account_id, institution_id, owner_id, classification, account_type, account_name, currency, status, opened_date, liquidity_profile, account_health)
VALUES
    ('c1d2e3f4-0000-4000-8000-000000000001', 'b1c2d3e4-0000-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Personal', 'Savings', 'Primary Savings Account', 'INR', 'Active', '2023-01-15', 'ShortTerm', 'Healthy'),
    ('c1d2e3f4-0000-4000-8000-000000000002', 'b1c2d3e4-0000-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Personal', 'Checking', 'Salary Account', 'INR', 'Active', '2023-01-15', 'ShortTerm', 'Healthy'),
    ('c1d2e3f4-0000-4000-8000-000000000003', 'b1c2d3e4-0000-4000-8000-000000000002', 'a1b2c3d4-0000-4000-8000-000000000001', 'Personal', 'Savings', 'SBI Fixed Deposit', 'INR', 'Active', '2024-03-01', 'LongTerm', 'Healthy'),
    ('c1d2e3f4-0000-4000-8000-000000000004', 'b1c2d3e4-0000-4000-8000-000000000003', 'a1b2c3d4-0000-4000-8000-000000000001', 'Personal', 'Investment', 'Zerodha Trading Account', 'INR', 'Active', '2023-06-01', 'LongTerm', 'Healthy'),
    ('c1d2e3f4-0000-4000-8000-000000000005', 'b1c2d3e4-0000-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Personal', 'Credit', 'HDFC Credit Card', 'INR', 'Active', '2023-02-01', 'ShortTerm', 'Healthy');

-- Update current balances via financial events (below) instead of directly

-- ============================================================
-- FINANCIAL EVENTS
-- ============================================================
INSERT INTO financial_events (event_id, user_id, event_type, event_date, effective_date, currency, amount, source, destination, description, origin, confidence, state)
VALUES
    -- Salary credits
    ('e1f2a3b4-0001-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Salary', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', 'INR', 150000, 'Employer', 'Salary Account', 'Monthly salary credit', 'Import', 'Confirmed', 'POSTED'),
    ('e1f2a3b4-0002-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Salary', NOW() - INTERVAL '35 days', NOW() - INTERVAL '35 days', 'INR', 145000, 'Employer', 'Salary Account', 'Monthly salary credit', 'Import', 'Confirmed', 'POSTED'),
    ('e1f2a3b4-0003-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Salary', NOW() - INTERVAL '65 days', NOW() - INTERVAL '65 days', 'INR', 150000, 'Employer', 'Salary Account', 'Monthly salary credit', 'Import', 'Confirmed', 'POSTED'),

    -- Expenses
    ('e1f2a3b4-0004-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Rent', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', 'INR', -25000, 'Salary Account', 'Landlord', 'Monthly rent payment', 'Import', 'Confirmed', 'POSTED'),
    ('e1f2a3b4-0005-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Groceries', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days', 'INR', -8500, 'Salary Account', 'BigBasket', 'Weekly groceries', 'Import', 'Confirmed', 'POSTED'),
    ('e1f2a3b4-0006-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Dining', NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day', 'INR', -3200, 'Salary Account', 'Swiggy', 'Restaurant delivery', 'Import', 'Confirmed', 'POSTED'),
    ('e1f2a3b4-0007-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Utilities', NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days', 'INR', -4500, 'Salary Account', 'Electricity Board', 'Electricity bill', 'Import', 'Confirmed', 'POSTED'),
    ('e1f2a3b4-0008-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Shopping', NOW() - INTERVAL '6 days', NOW() - INTERVAL '6 days', 'INR', -12000, 'Credit Card', 'Amazon Pay', 'Online shopping - electronics', 'Import', 'Confirmed', 'POSTED'),
    ('e1f2a3b4-0009-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Transport', NOW() - INTERVAL '7 days', NOW() - INTERVAL '7 days', 'INR', -2200, 'Salary Account', 'Uber', 'Cab rides', 'Import', 'Confirmed', 'POSTED'),

    -- Investments
    ('e1f2a3b4-0010-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Investment', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', 'INR', -30000, 'Salary Account', 'Zerodha', 'Monthly SIP investment', 'Import', 'Confirmed', 'POSTED'),
    ('e1f2a3b4-0011-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Investment', NOW() - INTERVAL '35 days', NOW() - INTERVAL '35 days', 'INR', -25000, 'Salary Account', 'Zerodha', 'Monthly SIP investment', 'Import', 'Confirmed', 'POSTED'),
    ('e1f2a3b4-0012-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Investment', NOW() - INTERVAL '65 days', NOW() - INTERVAL '65 days', 'INR', -30000, 'Salary Account', 'Zerodha', 'Monthly SIP investment', 'Import', 'Confirmed', 'POSTED'),

    -- Goal contributions
    ('e1f2a3b4-0013-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Transfer', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', 'INR', -10000, 'Salary Account', 'Emergency Fund', 'Monthly goal contribution', 'User', 'Confirmed', 'POSTED'),
    ('e1f2a3b4-0014-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Transfer', NOW() - INTERVAL '5 days', NOW() - INTERVAL '5 days', 'INR', -20000, 'Salary Account', 'Retirement', 'Monthly goal contribution', 'User', 'Confirmed', 'POSTED'),

    -- EMI
    ('e1f2a3b4-0015-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'EMI', NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days', 'INR', -18500, 'Salary Account', 'HDFC Home Loan', 'Home loan EMI', 'Import', 'Confirmed', 'POSTED');

-- ============================================================
-- GOALS
-- ============================================================
INSERT INTO goals (goal_id, user_id, name, importance, type, subtype, success_criteria, target_date, priority, status, created_at)
VALUES
    ('d1e2f3a4-0001-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Retirement Corpus', 'Mandatory', 'TimeBound', 'Retirement',
     '{"model": "TargetAmount", "target_value": 50000000}', NOW() + INTERVAL '25 years', 1, 'Active', NOW() - INTERVAL '365 days'),
    ('d1e2f3a4-0002-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Emergency Fund', 'Essential', 'OpenEnded', 'Savings',
     '{"model": "EmergencyFundMonths", "target_months": 6}', NULL, 2, 'Active', NOW() - INTERVAL '365 days'),
    ('d1e2f3a4-0003-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Buy a House', 'Dream', 'TimeBound', 'Purchase',
     '{"model": "TargetAmount", "target_value": 15000000}', NOW() + INTERVAL '5 years', 3, 'Active', NOW() - INTERVAL '180 days'),
    ('d1e2f3a4-0004-4000-8000-000000000001', 'a1b2c3d4-0000-4000-8000-000000000001', 'Kids Education', 'Essential', 'TimeBound', 'Education',
     '{"model": "TargetAmount", "target_value": 25000000}', NOW() + INTERVAL '15 years', 4, 'Active', NOW() - INTERVAL '90 days');

-- ============================================================
-- ALLOCATIONS
-- ============================================================
INSERT INTO allocations (allocation_id, goal_id, funding_source_id, funding_source_type, allocation_type, allocation_health, funding_commitment, priority, weight, currency, status, effective_date, reserved_amount, allocated_amount)
VALUES
    ('f1a2b3c4-0001-4000-8000-000000000001', 'd1e2f3a4-0001-4000-8000-000000000001', 'c1d2e3f4-0000-4000-8000-000000000004', 'Account', 'Auto', 'Healthy', 0, 1, 50.00, 'INR', 'Active', NOW(), 0, 500000),
    ('f1a2b3c4-0002-4000-8000-000000000001', 'd1e2f3a4-0002-4000-8000-000000000001', 'c1d2e3f4-0000-4000-8000-000000000001', 'Account', 'Fixed', 'Healthy', 300000, 2, 30.00, 'INR', 'Active', NOW(), 0, 200000),
    ('f1a2b3c4-0003-4000-8000-000000000001', 'd1e2f3a4-0003-4000-8000-000000000001', 'c1d2e3f4-0000-4000-8000-000000000003', 'Account', 'Fixed', 'Healthy', 0, 3, 20.00, 'INR', 'Active', NOW(), 0, 500000),
    ('f1a2b3c4-0004-4000-8000-000000000001', 'd1e2f3a4-0004-4000-8000-000000000001', 'c1d2e3f4-0000-4000-8000-000000000002', 'Account', 'Auto', 'Healthy', 0, 4, 10.00, 'INR', 'Draft', NOW(), 0, 100000);

-- ============================================================
-- ASSETS
-- ============================================================
INSERT INTO assets (asset_id, asset_name, classification, owner_id, currency, cost_basis, valuation_profile, valuation_method, quantity, unit_price, status)
VALUES
    ('a1b2c3d4-1001-4000-8000-000000000001', 'HDFC Bank Savings', 'CashAndCashEquivalent', 'a1b2c3d4-0000-4000-8000-000000000001', 'INR', 500000, 'ManualAssessment', 'StatementBalance', 1, 500000, 'Active'),
    ('a1b2c3d4-1002-4000-8000-000000000001', 'SBI Fixed Deposit', 'FixedDeposit', 'a1b2c3d4-0000-4000-8000-000000000001', 'INR', 500000, 'ManualAssessment', 'StatementBalance', 1, 525000, 'Active'),
    ('a1b2c3d4-1003-4000-8000-000000000001', 'Equity Portfolio', 'Equity', 'a1b2c3d4-0000-4000-8000-000000000001', 'INR', 800000, 'ManualAssessment', 'MarketValue', 1000, 960000, 'Active'),
    ('a1b2c3d4-1004-4000-8000-000000000001', 'Mutual Funds - Large Cap', 'MutualFund', 'a1b2c3d4-0000-4000-8000-000000000001', 'INR', 600000, 'ManualAssessment', 'MarketValue', 500, 720000, 'Active'),
    ('a1b2c3d4-1005-4000-8000-000000000001', 'Gold ETF', 'Gold', 'a1b2c3d4-0000-4000-8000-000000000001', 'INR', 200000, 'ManualAssessment', 'MarketValue', 50, 235000, 'Active'),
    ('a1b2c3d4-1006-4000-8000-000000000001', 'PPF Account', 'GovernmentSecurity', 'a1b2c3d4-0000-4000-8000-000000000001', 'INR', 350000, 'ManualAssessment', 'StatementBalance', 1, 380000, 'Active'),
    ('a1b2c3d4-1007-4000-8000-000000000001', 'Employee PF', 'RetirementFund', 'a1b2c3d4-0000-4000-8000-000000000001', 'INR', 1200000, 'ManualAssessment', 'StatementBalance', 1, 1280000, 'Active');

-- ============================================================
-- LIABILITIES
-- ============================================================
INSERT INTO liabilities (liability_id, name, classification, owner_id, currency, original_principal, interest_rate, interest_model, repayment_profile, installment_amount, remaining_installments, maturity_date, liability_health, status)
VALUES
    ('a1b2c3d4-2001-4000-8000-000000000001', 'Home Loan - HDFC', 'Mortgage', 'a1b2c3d4-0000-4000-8000-000000000001', 'INR', 5000000, 8.5000, 'ReducingBalance', 'FixedInstallment', 45000, 180, NOW() + INTERVAL '15 years', 'Stable', 'Active'),
    ('a1b2c3d4-2002-4000-8000-000000000001', 'Credit Card Dues', 'CreditCard', 'a1b2c3d4-0000-4000-8000-000000000001', 'INR', 45000, 42.0000, 'Simple', 'Revolving', 45000, 1, NOW() + INTERVAL '30 days', 'Warning', 'Active');

-- ============================================================
-- PORTFOLIO
-- ============================================================
INSERT INTO portfolios (portfolio_id, name, portfolio_type, owner_id, base_currency, status, portfolio_health)
VALUES
    ('a1b2c3d4-3001-4000-8000-000000000001', 'Primary Investment Portfolio', 'Investment', 'a1b2c3d4-0000-4000-8000-000000000001', 'INR', 'Active', 'Stable');

INSERT INTO portfolio_members (portfolio_id, entity_id, entity_type, weight)
VALUES
    ('a1b2c3d4-3001-4000-8000-000000000001', 'a1b2c3d4-1003-4000-8000-000000000001', 'Asset', 35.00),
    ('a1b2c3d4-3001-4000-8000-000000000001', 'a1b2c3d4-1004-4000-8000-000000000001', 'Asset', 25.00),
    ('a1b2c3d4-3001-4000-8000-000000000001', 'a1b2c3d4-1005-4000-8000-000000000001', 'Asset', 10.00),
    ('a1b2c3d4-3001-4000-8000-000000000001', 'a1b2c3d4-1006-4000-8000-000000000001', 'Asset', 15.00),
    ('a1b2c3d4-3001-4000-8000-000000000001', 'a1b2c3d4-1007-4000-8000-000000000001', 'Asset', 15.00);

-- ============================================================
-- ENGINE OUTPUTS (for demo display)
-- ============================================================

-- Health Score
INSERT INTO health_scores (score_id, overall_score, score_grade, score_data)
VALUES (
    'hscore-demo-001',
    72,
    'Good',
    '{
        "savings_rate": 18.5,
        "debt_to_income": 28.0,
        "emergency_fund_months": 4.5,
        "insurance_coverage": "Adequate",
        "investment_diversification": "Moderate",
        "retirement_progress": 35.0,
        "score_breakdown": {
            "savings": 68,
            "debt": 75,
            "emergency_fund": 60,
            "insurance": 80,
            "investments": 72,
            "retirement": 65,
            "budgeting": 78,
            "goal_progress": 70,
            "credit_health": 82,
            "tax_planning": 60,
            "financial_literacy": 75,
            "income_stability": 85,
            "net_worth_growth": 70
        }
    }'
);

-- Risk Assessment
INSERT INTO risk_assessments (assessment_id, composite_score, risk_level, output_data)
VALUES (
    'risk-demo-001',
    35,
    'Low',
    '{
        "category_scores": {
            "market_risk": 30,
            "concentration_risk": 45,
            "liquidity_risk": 25,
            "credit_risk": 20,
            "inflation_risk": 50,
            "interest_rate_risk": 40
        },
        "indicators": [
            {"name": "Debt to Income Ratio", "score": 28, "threshold": 40, "status": "healthy"},
            {"name": "Emergency Fund Coverage", "score": 4.5, "threshold": 3, "status": "healthy"},
            {"name": "Portfolio Concentration", "score": 45, "threshold": 60, "status": "moderate"}
        ]
    }'
);

-- Projection
INSERT INTO projection_outputs (output_id, projection_type, output_data)
VALUES (
    'proj-demo-001',
    'GoalCompletion',
    '{
        "goal_id": "d1e2f3a4-0001-4000-8000-000000000001",
        "goal_name": "Retirement Corpus",
        "current_progress": 12.5,
        "projected_completion_date": "2045-08-15",
        "target_date": "2050-01-01",
        "on_track": true,
        "monthly_contribution_required": 45000,
        "monthly_contribution_current": 30000,
        "funding_gap": 1500000,
        "confidence": "Medium",
        "assumptions": {
            "equity_return": 12.0,
            "debt_return": 7.0,
            "inflation_rate": 5.0,
            "salary_growth": 8.0
        }
    }'
);

-- Recommendation
INSERT INTO recommendations (recommendation_id, category, priority, score, rec_data)
VALUES
    ('rec-demo-001', 'GoalFunding', 1, 92.00, '{
        "title": "Increase Retirement Contribution",
        "summary": "Increasing monthly contribution by ₹5,000 will close the ₹15L funding gap",
        "expected_improvement": "15% higher corpus at retirement",
        "action": "Increase SIP",
        "affected_goals": ["Retirement Corpus"],
        "confidence": "High"
    }'),
    ('rec-demo-002', 'EmergencyFund', 2, 85.00, '{
        "title": "Build Emergency Fund to 6 Months",
        "summary": "Current emergency fund covers 4.5 months. Target: 6 months (₹9,00,000)",
        "expected_improvement": "Improved financial health score by 5 points",
        "action": "Increase allocation",
        "affected_goals": ["Emergency Fund"],
        "confidence": "High"
    }'),
    ('rec-demo-003', 'DebtManagement', 3, 78.00, '{
        "title": "Pay Down Credit Card Balance",
        "summary": "Current outstanding: ₹45,000 at 42% interest. Pay in full to avoid interest",
        "expected_improvement": "Save ₹15,750 in annual interest",
        "action": "Make payment",
        "affected_goals": ["Buy a House"],
        "confidence": "High"
    }');

-- Optimization
INSERT INTO optimizations (optimization_id, status, candidates_evaluated, candidates_valid, opt_data)
VALUES (
    'opt-demo-001',
    'Completed', 50, 12,
    '{
        "strategies": [
            {"rank": 1, "strategy": "Increase SIP by 10% annually", "score": 88.5, "summary": "Auto-escalate SIP contributions each year"},
            {"rank": 2, "strategy": "Rebalance to 70:30 equity:debt", "score": 82.0, "summary": "Shift allocation for better growth"},
            {"rank": 3, "strategy": "Consolidate FDs to achieve goals", "score": 75.0, "summary": "Use maturing FDs for goal funding"}
        ]
    }'
);

-- Simulation
INSERT INTO simulations (scenario_id, sim_type, sim_data)
VALUES (
    'sim-demo-001',
    'RetirementPlanning',
    '{
        "baseline": {
            "corpus_at_retirement": 38500000,
            "monthly_income": 125000,
            "goals_on_track": 2
        },
        "scenario_1": {
            "name": "Increase Contribution by ₹10,000/month",
            "corpus_at_retirement": 45200000,
            "monthly_income": 148000,
            "goals_on_track": 3
        },
        "scenario_2": {
            "name": "Retire 3 Years Earlier",
            "corpus_at_retirement": 35000000,
            "monthly_income": 110000,
            "goals_on_track": 2
        }
    }'
);
