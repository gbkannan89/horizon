# Phase 2 — Milestone 4: Financial Planning Endpoints

**Context Reminder:** Architecture frozen. Specs frozen. Stack frozen. Monorepo. No Redis/K8s/Grafana/testing.

**Milestones completed:** M1 (Docker), M2 (Auth/Profile), M3 (Transaction Edit/Delete).

**Current state:** Planning hub shows "Coming soon" cards for Budget, Retirement, Emergency Fund, Debt Payoff, Investment Planning. Only Projections and Scenarios are active.

## Objective

Implement the 5 missing financial planning backend endpoints. Then wire Flutter UI to them.

## Backend Audit

Read:
- `services/experiences/planning/register/register.go` — current endpoints (10 exist but none for the 5 missing features)
- `services/domains/` — check if any domain exists for these planning sub-types
- The planning experience currently has generic projections, scenarios, compare — we need specific endpoints

Since dedicated endpoints don't exist and we need a pragmatic MVP solution, create simple CSV-based or formula-based stub endpoints that return realistic demo data for each planning area.

## Implementation

### Backend — Add 5 endpoints to planning register

In `services/experiences/planning/register/register.go`:

1. `GET /api/v1/planning/budget` — returns budget categories with planned vs actual, surplus/deficit
2. `GET /api/v1/planning/retirement` — returns retirement readiness score, projected corpus, gap analysis
3. `GET /api/v1/planning/emergency-fund` — returns months covered, target amount, current savings
4. `GET /api/v1/planning/debt-payoff` — returns debt list, payoff projections, snowball/avalanche strategies
5. `GET /api/v1/planning/investment` — returns investment allocation suggestions, risk-adjusted returns

Each returns `{success, data: {...}, metadata}` format. Start with simple stub data that reads from financial_events and goals tables to derive realistic values (e.g., budget reads actual spending from events, retirement reads from retirement goal).

### Flutter — Replace "coming soon" cards with real screens

Replace the 5 "coming soon" cards in `planning_page.dart` with working navigation to:
- `/planning/budget` → BudgetPage
- `/planning/retirement` → RetirementPage
- `/planning/emergency-fund` → EmergencyFundPage
- `/planning/debt-payoff` → DebtPayoffPage
- `/planning/investment` → InvestmentPage

Each page shows:
- Summary card with key metric
- List of items or chart
- Loading/empty/error states using shared widgets

## Files to Modify

| File | Change |
|------|--------|
| `services/experiences/planning/register/register.go` | Add 5 routes |
| `services/experiences/planning/internal/api/handlers.go` | Add 5 handlers |
| `lib/features/planning/pages/planning_page.dart` | Replace coming-soon with working navigation |
| `lib/features/planning/pages/budget_page.dart` | NEW |
| `lib/features/planning/pages/retirement_page.dart` | NEW |
| `lib/features/planning/pages/emergency_fund_page.dart` | NEW |
| `lib/features/planning/pages/debt_payoff_page.dart` | NEW |
| `lib/features/planning/pages/investment_page.dart` | NEW |
| `lib/app/router.dart` | Add 5 routes |

## Verification

- [ ] `go build ./cmd/server/...` passes
- [ ] 5 new endpoints return data
- [ ] Planning hub shows active cards
- [ ] Each planning page loads and displays data
- [ ] `flutter analyze` — zero errors

## Deliverables

Backend endpoints, Flutter screens, routes.
