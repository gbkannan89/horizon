# Horizon — Milestone 3: Global Search & Shared Component Adoption

**Phase 2 of 9** | **Estimated: 1-2 sessions** | **Status: Pending**

---

## CONTEXT REMINDER

You are building Horizon — a Goal-Centric Personal Financial Operating System. Architecture is frozen. Specifications are frozen. Technology stack is frozen. Do NOT redesign the architecture.

### Prior Work Completed

- **Milestone 1 — Foundation Stabilization**: Auth (login, logout, splash, JWT persistence, token refresh), Settings (profile, preferences, theme persistence, security), Networking (connectivity_plus, ErrorMapper, retry interceptor).
- **Milestone 2 — Transactions Experience**: Complete transactions module with list, detail, create, archive, search, filters, pagination, infinite scroll, loading/empty/error/offline states. Routes at `/transactions`, `/transactions/:id`, `/transactions/add`.

### Design System References

- `horizon-spec/products/NAVIGATION_SYSTEM.md` Section 9 (Search Philosophy): Search is a global destination, not a filter. It searches across all financial entities, goals, recommendations, and events. Searchable entities: Goals (name, type, milestones), Accounts (name, type, institution), Assets, Liabilities, Recommendations, Timeline Events (description, category, amount), Insights (summary, category), Institutions.
- `horizon-spec/products/DESIGN_SYSTEM.md`: Information hierarchy, empty states, error states, loading states specifications. Functional color usage.
- `horizon-spec/engineering/FRONTEND_ENGINEERING_GUIDE.md`: Experience First, State Isolation, Offline First, Data Loading Philosophy (cache-first, graceful degradation).

### Current App Theme
- Seed color: `Color(0xFF00897B)` (Teal)
- Material 3 enabled
- Shared constants: `pageSize = 25`, `cacheDefaultTTLMinutes = 5`
- Money format: Indian numbering (Cr, L, K) with ₹ prefix
- Existing shared: `lib/shared/constants/app_constants.dart`, `lib/shared/extensions/extensions.dart`

### Flutter App Location

`C:\Kannan\Horizon\horizon\apps\mobile\`

---

## STANDARD OPERATING PROCEDURE

1. **Backend Audit** — Verify every required endpoint exists.
2. **Implementation Plan** — List every file to create/modify.
3. **Implementation** — Build. No placeholders. No TODOs.
4. **Build Verification** — `flutter analyze` — zero errors.
5. **Functional Verification** — Checklist below.
6. **Final Report** — Files modified, API inventory, remaining gaps.
7. **Message to Solution Architect**

---

## PART A: GLOBAL SEARCH

### Step 1: Backend Audit

#### 1A: Unified Search Endpoint

Search for a single unified endpoint: `GET /api/v1/search?q=...` or `GET /api/v1/global-search?q=...`

Check in:
- `services/infra/api-gateway/internal/router/routes.go`
- `cmd/server/main.go` (monolith routes)
- ALL `services/*/register/*.go` files

**Expected finding**: DOES NOT EXIST. If it does exist, document its response shape and use it as the primary endpoint.

#### 1B: Per-Module Search Endpoints

For each module, verify these EXIST:

| Module | Endpoint | Expected Response Shape |
|--------|----------|------------------------|
| Transactions | `GET /api/v1/transactions/search?q=` | `{success, data: {transactions[], total}, metadata}` |
| Timeline | `GET /api/v1/timeline/search?q=` | `{success, data: {items[], count}, metadata}` |
| Insights | `GET /api/v1/insights/search?q=` | `{success, data: {insights[], count}, metadata}` |
| Accounts | `GET /api/v1/accounts?q=` or `GET /api/v1/accounts/search?q=` | Check if any search param exists |
| Goals | `GET /api/v1/goals?q=` or `GET /api/v1/goals/experience?q=` | Check if any search param exists |
| Portfolio | `GET /api/v1/portfolio?q=` or similar | Check if any search param exists |

**Expected finding**: Transactions and Timeline have search endpoints. Insights has search. Accounts, Goals, Portfolio likely do NOT have dedicated search endpoints — will need client-side filtering or broader query.

#### 1C: Search History / Recent Searches

Check for: `GET /api/v1/search-history` or `GET /api/v1/recent-searches`

**Expected finding**: DOES NOT EXIST. Implement recent searches locally (in-memory list in the StateNotifier, persisted to a simple file or ignored).

### Implementation Plan

Since no unified search endpoint exists, implement **client-side parallel search**:
- `Future.wait` across all available per-module search endpoints
- For modules without search endpoints, skip them (document as limitation)
- Aggregate results grouped by module
- Map each result to a common `SearchResultItem` model that includes navigation route info

### Models (`features/search/models/search_models.dart`)

- `SearchResultItem`:
  - `String id` — entity ID
  - `String title` — primary display text
  - `String subtitle` — secondary text (category, date, etc.)
  - `String module` — one of: "account", "transaction", "goal", "portfolio", "timeline"
  - `String? amount` — formatted amount if applicable
  - `String? date` — relative date if applicable
  - `IconData icon` — module-specific icon
  - `String route` — GoRouter path to navigate to (e.g., `/transactions/{id}`)
  - `Map<String, dynamic>? routeParams` — optional extra route data
- `SearchGroup` — module name + `List<SearchResultItem>` + icon
- `SearchState` — `status`, `query`, `groups[]`, `recentSearches[]`, `error`

### Repository (`features/search/repository/search_repository.dart`)

- `GlobalSearchRepository` with `ApiClient` dependency
- `Future<List<SearchResultItem>> searchAll(String query)`:
  - Uses `Future.wait` to call all module search endpoints in parallel
  - Maps each module response to `SearchResultItem` list
  - Wraps each call in try/catch so one module failure doesn't block others
  - Returns flattened list of results

### State (`features/search/state/search_state.dart`)

- `SearchStatus` enum: `idle, searching, loaded, error`
- `SearchNotifier`:
  - `search(String query)` — 300ms debounce, calls repository, groups results by module
  - `clearSearch()` — reset state
  - `loadRecent()` — load from in-memory list
  - `saveRecent(String query)` — add to recent list (max 10)
  - `clearHistory()` — clear recent list

### Page (`features/search/pages/search_page.dart`)

- Full-screen search page at route `/search`
- `TextField` with autofocus, search icon prefix, clear button suffix
- 300ms debounce via `Timer`
- **States**:
  - `idle` (no query): Show recent searches section with "Clear history" button
  - `searching`: Show skeleton per module (module header + 2-3 shimmer rows)
  - `loaded`: Results grouped by module with `SectionHeader` chips, each result is a tappable card with icon + title + subtitle + amount
  - `error`: Per-module error message (not full-page failure)
  - `empty`: "No results found" message
- Each result navigates to the appropriate detail page via GoRouter

### Navigation Integration

- Add search icon to the Dashboard AppBar (left side or right side)
- Route: `/search` -> `SearchPage`
- Result tap navigates:
  - Account: `/accounts/{id}`
  - Transaction: `/transactions/{id}`
  - Goal: `/goals/{id}`
  - Portfolio: `/portfolio` or `/portfolio/{id}`
  - Timeline: `/timeline/{id}`

---

## PART B: SHARED COMPONENT MIGRATION

### Create Shared Widget Library (`lib/shared/widgets/`)

Create these reusable widgets (one file per widget, or group related ones):

1. **`shared_loading.dart`**
   - `SharedLoadingView` — centered `CircularProgressIndicator` with optional message
   - `SharedSkeletonList` — configurable number of shimmer card placeholders (like the 8-card skeleton in timeline)
   - `SharedSkeletonCard` — single shimmer placeholder card with configurable height

2. **`shared_error.dart`**
   - `SharedErrorView` — error icon + message + "Try Again" `FilledButton`
   - Accepts `String? message`, `VoidCallback? onRetry`, `VoidCallback? onDismiss`

3. **`shared_empty.dart`**
   - `SharedEmptyView` — icon + title + subtitle + optional action button + optional clear-filters button
   - Accepts `IconData icon`, `String title`, `String? subtitle`, `String? actionLabel`, `VoidCallback? onAction`

4. **`shared_offline.dart`**
   - `SharedOfflineView` — `wifi_off` icon + "You're offline" text

5. **`shared_section_header.dart`**
   - `SharedSectionHeader` — consistent section title with optional "See all" action
   - Accepts `String title`, `VoidCallback? onSeeAll`, `Widget? trailing`

6. **`shared_card.dart`**
   - `SharedCard` — wrapper with consistent Material 3 card styling (Clip.antiAlias, elevation, border radius)
   - Accepts `Widget child`, `VoidCallback? onTap` for InkWell

7. **`shared_status_badge.dart`**
   - `SharedStatusBadge` — consistent colored chip
   - Accepts `String label`, `Color color`, `IconData? icon`

8. **`shared_search_bar.dart`**
   - `SharedSearchBar` — debounced text field with search icon, clear button, autofocus option
   - Accepts `ValueChanged<String> onChanged`, `VoidCallback? onCancel`, `TextEditingController? controller`
   - Internally manages 300ms debounce

9. **`shared_filter_chip.dart`**
   - `SharedFilterChip` — consistent filter chip wrapping `FilterChip` with visual density

### Theme Standardization (`lib/app/theme.dart`)

Add these constants:
```dart
class AppSpacing {
  static const double xs = 4;
  static const double sm = 8;
  static const double md = 16;
  static const double lg = 24;
  static const double xl = 32;
  static const double xxl = 48;
}

class AppRadius {
  static const double sm = 8;
  static const double md = 12;
  static const double lg = 16;
  static const double xl = 24;
}

class AppElevation {
  static const double low = 1;
  static const double medium = 2;
  static const double high = 4;
}

class AppIconSize {
  static const double sm = 16;
  static const double md = 20;
  static const double lg = 24;
  static const double xl = 32;
}
```

### Migrate Existing Features

For each feature below, replace duplicated UI patterns with shared widgets. Do NOT change feature behavior — only reduce duplication.

#### Dashboard (`features/dashboard/`)
- `pages/dashboard_page.dart`: Replace loading skeleton with `SharedSkeletonList`
- `pages/dashboard_page.dart`: Replace error state with `SharedErrorView`
- `pages/dashboard_page.dart`: Replace empty state with `SharedEmptyView`

#### Accounts (`features/accounts/`)
- `pages/accounts_page.dart`: Replace loading/error/empty with shared widgets
- `widgets/acct_widgets.dart`: Use `SharedCard` wrapper, `SharedStatusBadge` for status indicators

#### Goals (`features/goals/`)
- `pages/goals_page.dart`: Replace loading/error/empty with shared widgets
- `pages/goal_detail_page.dart`: Replace loading/error/empty with shared widgets
- `widgets/goal_widgets.dart`: Use `SharedCard`, `SharedStatusBadge`

#### Portfolio (`features/portfolio/`)
- `pages/portfolio_page.dart`: Replace loading/error/empty with shared widgets
- `widgets/pf_widgets.dart`: Use `SharedCard`, `SharedSectionHeader`

#### Timeline (`features/timeline/`)
- `pages/timeline_page.dart`: Replace skeleton with `SharedSkeletonList`, error with `SharedErrorView`, empty with `SharedEmptyView`, offline with `SharedOfflineView`

#### Settings (`features/settings/`)
- `pages/edit_profile_page.dart`: Replace loading state with `SharedLoadingView`
- `pages/preferences_page.dart`: Replace loading/error with shared widgets
- `pages/privacy_page.dart`: Replace error with `SharedErrorView`
- `pages/profile_page.dart`: Replace loading/error with shared widgets

---

## VERIFICATION CHECKLIST

### Search
- [ ] Search icon visible on Dashboard
- [ ] Typing query triggers 300ms debounced search
- [ ] Results appear grouped by module with section headers
- [ ] Each result navigates to correct detail page
- [ ] Recent searches shown when query is empty
- [ ] Clear history removes recent searches
- [ ] Empty state shown when no results
- [ ] Per-module error handling (one failing module doesn't block others)

### Shared Widgets
- [ ] All 9 shared widgets created in `lib/shared/widgets/`
- [ ] Theme constants (`AppSpacing`, `AppRadius`, etc.) added to `theme.dart`
- [ ] Dashboard uses shared widgets for loading/error/empty
- [ ] Accounts uses shared widgets for loading/error/empty + card/status badge
- [ ] Goals uses shared widgets for loading/error/empty + card/status badge
- [ ] Portfolio uses shared widgets for loading/error/empty + card/section header
- [ ] Timeline uses shared widgets for skeleton/error/empty/offline
- [ ] Settings uses shared widgets for loading/error

### Build
- [ ] `flutter analyze` — zero errors
- [ ] App builds without issues
- [ ] Feature behavior unchanged (visual diff is only cosmetic)

---

## DELIVERABLE TEMPLATE

```markdown
## 1. Search Architecture
...
## 2. Search API Inventory
...
## 3. Shared Widgets Created
...
## 4. Screens Migrated
...
## 5. Remaining Duplicated Widgets
...
## 6. Files Modified
...
## 7. Verification Checklist
...
## 8. Message to Solution Architect
```
