# Horizon — Milestone 5: Insights

**Phase 4 of 9** | **Estimated: 1-2 sessions** | **Status: Pending**

---

## CONTEXT REMINDER

**Architecture is frozen. Specifications are frozen. Technology stack is frozen. Do NOT redesign the architecture.**

### Prior Work Completed

- **M1** Foundation Stabilization (Auth, Settings, Networking, Error handling)
- **M2** Transactions Experience (list, detail, create, archive, search, filters, pagination, infinite scroll)
- **M3** Global Search + Shared Widgets (cross-module search, `SharedLoadingView`, `SharedErrorView`, `SharedEmptyView`, `SharedOfflineView`, `SharedCard`, `SharedSectionHeader`, `SharedStatusBadge`, `SharedSearchBar`, `SharedFilterChip`, theme constants `AppSpacing`, `AppRadius`, `AppElevation`, `AppIconSize`)
- **M4** AI Advisor (chat interface with markdown rendering, suggested prompts, context awareness, conversation management, provider switching)

### Design System References

- `horizon-spec/products/DESIGN_SYSTEM.md`: Charts Philosophy — charts are decision-support tools, every chart must lead to an action. Information hierarchy, functional colors.
- `horizon-spec/products/UX_PHILOSOPHY.md`: Decision-First Design framework. Every screen answers "What Should I Do Next?". Financial Storytelling.
- `horizon-spec/engineering/FRONTEND_ENGINEERING_GUIDE.md`: Experience First, State Isolation, Data Loading Philosophy.
- All insight category colors and priority levels from `INSIGHTS_EXPERIENCE_REPORT.md` (16 insight categories, 5 priority levels).

### Current App State
- `/insights` route exists with a PLACEHOLDER (7-line stub)
- `fl_chart: ^0.68.0` already in `pubspec.yaml` but never used
- Shared widget library available in `lib/shared/widgets/`
- Theme constants available in `lib/app/theme.dart`

### Flutter App Location

`C:\Kannan\Horizon\horizon\apps\mobile\`

---

## STANDARD OPERATING PROCEDURE

1. **Backend Audit** — Verify every required endpoint exists.
2. **Implementation Plan** — List every file.
3. **Implementation** — Build. NO PLACEHOLDERS. NO TODOs.
4. **Build Verification** — `flutter analyze` — zero errors.
5. **Functional Verification** — Checklist.
6. **Final Report**.
7. **Message to Solution Architect**.

---

## Step 1: Backend Audit

### Required Endpoints

Check `C:\Kannan\Horizon\horizon\services\experiences\insights\register\register.go` and `handlers.go`:

| Method | Path | Purpose | Expected |
|--------|------|---------|----------|
| `GET` | `/api/v1/insights` | Feed of all insights | ✅ EXISTS |
| `GET` | `/api/v1/insights/summary` | Dashboard summary | ✅ EXISTS |
| `GET` | `/api/v1/insights/opportunities` | Opportunities list | ✅ EXISTS |
| `GET` | `/api/v1/insights/warnings` | Warnings list | ✅ EXISTS |
| `GET` | `/api/v1/insights/achievements` | Achievements list | ✅ EXISTS |
| `GET` | `/api/v1/insights/forecast` | Forecast data | ✅ EXISTS |
| `GET` | `/api/v1/insights/trends` | Trend analysis | ✅ EXISTS |
| `GET` | `/api/v1/insights/timeline` | Insights timeline | ✅ EXISTS |
| `GET` | `/api/v1/insights/search?q=` | Search insights | ✅ EXISTS |
| `GET` | `/api/v1/insights/{id}` | Single insight detail | ✅ EXISTS |

Record exact response shapes for each. If ANY endpoint is missing, STOP and report.

### Expected Data Shapes

From `handlers.go`, each endpoint returns `{success, data, metadata}`. The data shapes include:
- Insights have categories (Spending, Income, Goal, Portfolio, Risk, Opportunity, Warning, Achievement, Forecast, Trend, Behavior, Saving, Investment, Tax, Subscription, Custom)
- Priorities: Critical (P1), High (P2), Medium (P3), Low (P4), Info (P5)
- Summary includes total, critical count, new count, etc.

---

## Step 2: Implementation Plan

Structure:
```
features/insights/
├── models/insight_models.dart
├── repository/insight_repository.dart
├── state/insight_state.dart
├── pages/
│   ├── insights_page.dart        (REPLACE placeholder)
│   └── insight_detail_page.dart
├── widgets/insight_widgets.dart
```

Use `fl_chart` for trend charts (already in `pubspec.yaml`). Use shared widgets from `lib/shared/widgets/`.

---

## Step 3: Implementation

### Models (`features/insights/models/insight_models.dart`)

- `InsightResponse` — `{success, data: InsightData, metadata}`
- `InsightData` — `{insights[], total, critical_count, new_count, summary_cards[]}`
- `InsightItem` — `id, category, subcategory, title, summary, description, priority (P1-P5), impact, confidence, affected_entities[], actions[], created_at, read, dismissed`
- `InsightSummary` — dashboard-level summary data
- `TrendData` — `period, values[], change_percent` for trend charts
- `ForecastItem` — `period, projected_value, lower_bound, upper_bound, confidence`

### Repository (`features/insights/repository/insight_repository.dart`)

- `getInsights()` — `GET /api/v1/insights`
- `getSummary()` — `GET /api/v1/insights/summary`
- `getOpportunities()` — `GET /api/v1/insights/opportunities`
- `getWarnings()` — `GET /api/v1/insights/warnings`
- `getAchievements()` — `GET /api/v1/insights/achievements`
- `getForecast()` — `GET /api/v1/insights/forecast`
- `getTrends()` — `GET /api/v1/insights/trends`
- `search({query})` — `GET /api/v1/insights/search?q=`
- `getById(id)` — `GET /api/v1/insights/{id}`

### State (`features/insights/state/insight_state.dart`)

- `InsightStatus` enum: `initial, loading, loaded, error, empty`
- `InsightState` with `copyWith`: `status`, `insights[]`, `summary`, `opportunities[]`, `warnings[]`, `achievements[]`, `trends[]`, `forecast[]`, `activeTab`, `error`
- `InsightNotifier`:
  - `init()` — load summary + main feed
  - `loadTab(String tab)` — load opportunities, warnings, or achievements based on selected tab
  - `refresh()` — reload all
  - `search(String query)` — search with debounce
  - `filterByCategory(String? category)` — filter insights

### Pages

#### Insights Page (`features/insights/pages/insights_page.dart`)

REPLACE placeholder.

**Layout:**
- `AppBar` with "Insights" title + search icon
- Summary card at top (total insights, critical count, new count)
- Tab bar: All | Opportunities | Warnings | Achievements | Forecast
- Filter chips row: by category (Spending, Income, Goal, Portfolio, Risk, etc.)
- List of insight cards below

**Each insight card:**
- Category icon + color (from category mapping)
- Title + summary
- Priority badge (P1=red, P2=orange, P3=amber, P4=grey, P5=blue)
- Relative date
- Tap navigates to detail page

**Chart section (Trends tab):**
- Line chart using `fl_chart`'s `LineChart`
- Period selector (1M, 3M, 6M, 1Y)
- Show trend line with optional confidence band

#### Insight Detail Page (`features/insights/pages/insight_detail_page.dart`)

- Full insight detail
- Category badge, priority badge
- Description
- Impact section
- Affected entities list (tappable to navigate)
- Actions section (if applicable)
- Confidence indicator

### Widgets (`features/insights/widgets/insight_widgets.dart`)

- `InsightCard` — reusable insight list item card
- `InsightCategoryBadge` — colored badge for insight category
- `InsightPriorityBadge` — priority level indicator
- `TrendChart` — fl_chart line chart wrapper
- `InsightSummaryCard` — dashboard summary

### Color Mapping (from Insights Experience Report)

Map each insight category to an icon + color:
- Spending: `Icons.shopping_cart`, Colors.red
- Income: `Icons.trending_up`, Colors.green
- Goal: `Icons.flag`, Colors.blue
- Portfolio: `Icons.pie_chart`, Colors.indigo
- Risk: `Icons.warning`, Colors.amber
- Opportunity: `Icons.lightbulb`, Colors.teal
- Warning: `Icons.error_outline`, Colors.orange
- Achievement: `Icons.emoji_events`, Colors.amber (gold)
- Forecast: `Icons.query_stats`, Colors.cyan
- Trend: `Icons.show_chart`, Colors.blue
- etc.

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
- [ ] Insight page shows summary + feed on load
- [ ] Tab switching loads opportunities/warnings/achievements/forecast
- [ ] Filter chips filter by category
- [ ] Search works with debounce
- [ ] Each insight card navigates to detail page
- [ ] Detail page shows full insight data
- [ ] Trend chart renders with fl_chart (no chart rendering errors)
- [ ] Period selector changes chart data
- [ ] Loading/empty/error/offline states work
- [ ] Shared widgets used for common states
- [ ] Pull-to-refresh works

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
## 5. Shared Widgets Used
...
## 6. Remaining Backend Gaps
...
## 7. Verification Checklist
...
## 8. Message to Solution Architect
```
