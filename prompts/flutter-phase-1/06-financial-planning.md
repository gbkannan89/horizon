# Horizon — Milestone 6: Financial Planning

**Phase 5 of 9** | **Estimated: 1-2 sessions** | **Status: Pending**

---

## CONTEXT REMINDER

**Architecture is frozen. Specifications are frozen. Technology stack is frozen. Do NOT redesign the architecture.**

### Prior Work Completed

- **M1** Foundation Stabilization (Auth, Settings, Networking, Error handling)
- **M2** Transactions Experience (list, detail, create, archive, search, filters, pagination)
- **M3** Global Search + Shared Widgets (cross-module search, shared UI library, theme constants)
- **M4** AI Advisor (chat interface with markdown, suggested prompts, context awareness, provider switching)
- **M5** Insights (spending/income analysis, trends with fl_chart, opportunities, warnings, achievements, forecast, search, filter by category)

### Design System References

- `horizon-spec/products/UX_PHILOSOPHY.md`: Forward-Looking principle, Decision-First Design, Progressive Disclosure levels (Essential, Standard, Advanced, Expert).
- `horizon-spec/products/DESIGN_SYSTEM.md`: Confidence indicators for projections, information hierarchy for planning data, card philosophy.
- `horizon-spec/products/NAVIGATION_SYSTEM.md` Section 7 (Planning Journey Flow): Goals -> Planning -> Adjust Assumptions -> Run Simulation -> Compare -> Accept as plan.
- `horizon-spec/engineering/FRONTEND_ENGINEERING_GUIDE.md`: Deterministic Rendering, State Isolation.

### Current App State
- `/planning` route exists with a PLACEHOLDER (7-line stub)
- Shared widgets available
- fl_chart available
- Theme constants available

### Flutter App Location

`C:\Kannan\Horizon\horizon\apps\mobile\`

---

## STANDARD OPERATING PROCEDURE

1. **Backend Audit** — Verify every required endpoint. DO NOT PROCEED if critical endpoints are missing.
2. **Implementation Plan** — List every file.
3. **Implementation** — Build. NO PLACEHOLDERS. NO TODOs. For features with missing backend endpoints, show clear "coming soon" messages.
4. **Build Verification** — `flutter analyze` — zero errors.
5. **Functional Verification** — Checklist.
6. **Final Report**.
7. **Message to Solution Architect**.

---

## Step 1: Backend Audit

### Required Endpoints

Check `C:\Kannan\Horizon\horizon\services\experiences\planning\register\register.go`:

| Method | Path | Purpose | Expected Finding |
|--------|------|---------|------------------|
| `GET` | `/api/v1/planning` | Planning dashboard overview | ✅ EXISTS |
| `GET` | `/api/v1/planning/dashboard` | Dashboard data | ✅ EXISTS |
| `GET` | `/api/v1/planning/projections` | Financial projections | ✅ EXISTS |
| `GET` | `/api/v1/planning/scenarios` | Scenario list | ✅ EXISTS |
| `GET` | `/api/v1/planning/compare` | Scenario comparison | ✅ EXISTS |
| `GET` | `/api/v1/planning/recommendations` | Planning recommendations | ✅ EXISTS |
| `GET` | `/api/v1/planning/optimizations` | Planning optimizations | ✅ EXISTS |
| `GET` | `/api/v1/planning/risk` | Risk assessment | ✅ EXISTS |
| `GET` | `/api/v1/planning/health` | Planning health | ✅ EXISTS |
| `GET` | `/api/v1/planning/timeline` | Planning timeline | ✅ EXISTS |

**Check for dedicated endpoints:**
| Endpoint | Expected Finding |
|----------|------------------|
| `/api/v1/planning/budget` | ❌ DOES NOT EXIST |
| `/api/v1/planning/retirement` | ❌ DOES NOT EXIST |
| `/api/v1/planning/emergency-fund` | ❌ DOES NOT EXIST |
| `/api/v1/planning/debt-payoff` | ❌ DOES NOT EXIST |
| `/api/v1/planning/investment` | ❌ DOES NOT EXIST |

Document all missing endpoints. For the missing ones, show "Coming soon" in the UI with a description of the planned feature.

Record exact response shapes for all existing endpoints.

---

## Step 2: Implementation Plan

Structure:
```
features/planning/
├── models/planning_models.dart
├── repository/planning_repository.dart
├── state/planning_state.dart
├── pages/
│   ├── planning_page.dart         (REPLACE placeholder — dashboard hub)
│   ├── projection_page.dart
│   ├── scenario_page.dart
│   ├── scenario_compare_page.dart
│   └── planning_detail_page.dart  (generic for recommendations/risk/health)
├── widgets/planning_widgets.dart
```

---

## Step 3: Implementation

### Models (`features/planning/models/planning_models.dart`)

Follow the backend response shapes from the audit.

- `PlanningDashboardResponse` — `{success, data: PlanningDashboard, metadata}`
- `PlanningDashboard` — summary cards, scenario count, health score, risk level
- `ProjectionData` — period projections with income/expenses/net worth over time
- `Scenario` — id, name, description, type, assumptions[], projected_outcomes[], created_at
- `PlanComparison` — scenarios compared side by side with metrics
- `PlanningCard` — generic card for recommendations/optimizations/risk/health (like the CardView from backend)

### Repository (`features/planning/repository/planning_repository.dart`)

- `getDashboard()` — `GET /api/v1/planning/dashboard`
- `getProjections()` — `GET /api/v1/planning/projections`
- `getScenarios()` — `GET /api/v1/planning/scenarios`
- `getComparison(scenarioIds)` — `GET /api/v1/planning/compare`
- `getRecommendations()` — `GET /api/v1/planning/recommendations`
- `getOptimizations()` — `GET /api/v1/planning/optimizations`
- `getRisk()` — `GET /api/v1/planning/risk`
- `getHealth()` — `GET /api/v1/planning/health`

### Pages

#### Planning Dashboard (`features/planning/pages/planning_page.dart`)

REPLACE placeholder with hub page.

**Hub layout showing cards for each planning area:**
1. **Budget Planning** — "Coming soon" card with description
2. **Retirement Planning** — "Coming soon" card with description
3. **Emergency Fund** — "Coming soon" card with description
4. **Debt Payoff** — "Coming soon" card with description
5. **Investment Planning** — "Coming soon" card with description
6. **Projections** — Active card linking to projection page
7. **Scenarios** — Active card linking to scenario page
8. **Health/Risk** — Active cards with current status

Each card: icon + title + short description/subtitle + arrow for navigation.

#### Projection Page (`features/planning/pages/projection_page.dart`)

- Shows financial projections over time
- Line chart using fl_chart (income/expenses/net worth lines)
- Period selector (1Y, 3Y, 5Y, 10Y)
- Key metrics: projected net worth, annual income, annual expenses

#### Scenario Page (`features/planning/pages/scenario_page.dart`)

- List of saved scenarios
- Each scenario card: name, description, type, key metrics
- Tap to view scenario detail
- Compare button to select scenarios for comparison

#### Comparison Page (`features/planning/pages/scenario_compare_page.dart`)

- Side-by-side comparison of selected scenarios
- Metrics table: current vs scenario A vs scenario B
- Visual comparison bars

### Color Mapping

Planning card colors (from PLANNING_EXPERIENCE_REPORT.md):
- Overview: teal
- Projections: cyan
- Scenarios: indigo
- Comparison: purple
- Recommendations: amber
- Optimizations: cyan
- Risk: red (score-based)
- Health: green (score-based)
- Timeline: brown

---

## Step 4: Build Verification

```powershell
cd C:\Kannan\Horizon\horizon\apps\mobile
flutter pub get
flutter analyze
```

Zero errors required.

---

## VERIFICATION CHECKLIST

- [ ] `flutter analyze` — zero errors
- [ ] Planning dashboard shows all planning area cards
- [ ] Budget/Retirement/Emergency/Debt/Investment cards show "Coming soon" with descriptions
- [ ] Projections page loads and chart renders
- [ ] Period selector changes chart data
- [ ] Scenarios page loads list
- [ ] Comparison page shows selected scenarios
- [ ] Recommendations/optimizations/risk/health load
- [ ] Loading/empty/error/offline states work
- [ ] Shared widgets used for common states

---

## DELIVERABLE TEMPLATE

```markdown
## 1. API Inventory
...
## 2. Screens Implemented
...
## 3. Files Created
...
## 4. Files Modified
...
## 5. Missing Backend Endpoints (documented with "coming soon")
...
## 6. Verification Checklist
...
## 7. Message to Solution Architect
```
