# Horizon — Milestone 2: Transactions Experience

**Phase 1 of 9** | **Estimated: 1 session** | **Status: Pending**

---

## CONTEXT REMINDER

You are building Horizon — a Goal-Centric Personal Financial Operating System. Architecture is frozen. Specifications are frozen. Technology stack is frozen. Do NOT redesign the architecture.

### Prior Work Completed (before this milestone)

- **Milestone 1 — Foundation Stabilization**: Authentication (login, logout, splash, JWT persistence, token refresh via `POST /auth/refresh`), Settings (profile, preferences, theme with backend persistence, security page), Networking (connectivity monitor via `connectivity_plus`, ErrorMapper with full `DioException` handling, retry interceptor with retry limit), Error handling on audited screens.

### Design System References

Read these before implementing (key excerpts inline, but refer to full files):
- `horizon-spec/products/DESIGN_SYSTEM.md` — Design principles: Calm, Trustworthy, Professional, Approachable, Confident, Predictable, Explainable, Minimal Cognitive Load. Information hierarchy: Primary Metric -> Label -> Change Indicator -> Confidence -> Secondary Info. Money format: Indian numbering (Cr, L, K) with ₹ prefix.
- `horizon-spec/products/UX_PHILOSOPHY.md` — Financial Events are Source of Truth. Every screen answers "What Should I Do Next?".
- `horizon-spec/products/NAVIGATION_SYSTEM.md` — Transactions are secondary navigation (accessible via More page).
- `horizon-spec/engineering/FRONTEND_ENGINEERING_GUIDE.md` — Experience First, Presentation Without Business Logic, State Isolation, Offline First.

### Flutter App Location

`C:\Kannan\Horizon\horizon\apps\mobile\`

---

## STANDARD OPERATING PROCEDURE

Follow these steps in order. Do NOT skip any step.

### Step 1: Backend Audit

Before writing any Flutter code, verify that EVERY endpoint listed below exists by searching the Go backend code at `C:\Kannan\Horizon\horizon\services\` and `C:\Kannan\Horizon\horizon\services\infra\api-gateway\internal\router\routes.go`.

**Required endpoints:**

| Method | Path | Purpose |
|--------|------|---------|
| `GET` | `/api/v1/transactions` | Cursor-paginated list with filters |
| `GET` | `/api/v1/transactions/{id}` | Single transaction detail |
| `GET` | `/api/v1/transactions/search?q=` | Search by merchant/category/notes |
| `GET` | `/api/v1/transactions/summary` | Monthly income/expense summary |
| `POST` | `/api/v1/events` | Create transaction |
| `POST` | `/api/v1/events/{id}/archive` | Archive/soft-delete |

For each endpoint, record:
- Full path + method
- Request body/query params shape
- Response JSON shape
- File path where it's defined
- Whether it's a real implementation or mock

If ANY endpoint is missing, STOP. Report which ones are missing in the deliverable. Do NOT fabricate client code for endpoints that don't exist.

**Note from prior audit**: `PUT` (update) and `DELETE` (hard delete) endpoints DO NOT EXIST. `POST /api/v1/events` is a mock that accepts only `{type, amount, currency}`. Document these limitations.

### Step 2: Implementation Plan

Describe every file you will create or modify before writing code. Get implicit approval from this prompt by listing the plan.

### Step 3: Implementation

Build the complete feature. No placeholders. No TODOs. Every screen must handle: loading, loaded, empty, error, offline states.

### Step 4: Build Verification

```powershell
cd C:\Kannan\Horizon\horizon\apps\mobile
flutter pub get
flutter analyze
```

Zero errors required. Pre-existing warnings/info are acceptable.

### Step 5: Functional Verification

Run through the verification checklist below. Mark each item.

### Step 6: Final Report

Provide: files modified, API inventory, screens implemented, remaining backend gaps.

### Step 7: Message to Solution Architect

---

## IMPLEMENTATION REQUIREMENTS

### Directory Structure to Create

```
lib/features/transactions/
├── models/transaction_models.dart
├── repository/transaction_repository.dart
├── state/transaction_state.dart
├── pages/
│   ├── transactions_page.dart
│   ├── transaction_detail_page.dart
│   └── transaction_form_page.dart
└── widgets/
    ├── transaction_widgets.dart
    └── transaction_filters.dart
```

### Models (`transaction_models.dart`)

Follow the pattern in `features/timeline/models/timeline_models.dart`.

- `TransactionListResponse` — `{success, data: TransactionListData, metadata}`
- `TransactionListData` — `{transactions: [], cursor, has_more, total}`
- `TransactionDetailResponse` — `{success, data: TransactionEvent, metadata}`
- `TransactionSummaryResponse` — `{success, data: {period_income, period_expenses, net_flow, income_count, expense_count, total_count}}`
- `TransactionEvent` — Map the full backend `EventResult` shape: `event_id`, `user_id`, `event_type`, `amount` (double), `currency`, `event_date`, `effective_date`, `description`, `state`, `origin`, `confidence`, `created_by`, `source`, `destination`, `reference`, `category`, `notes`, `tags`, `imported_from`, `reversal_of_event_id`, `correlation_id`, `order_index`, `created_at`, `updated_at`
  - Helper getters: `isIncome` (amount >= 0), `isExpense` (amount < 0), `formattedDate` (relative like "2d ago"), `formattedAmount` (Indian: Cr/L/K)
- `TransactionFilter` — `startDate`, `endDate`, `accountId`, `category`, `type`, `minAmount`, `maxAmount`, `merchant`
  - `toQuery()` method returning `Map<String, dynamic>`
  - `copyWith()` for immutable updates
  - `hasActiveFilters`, `activeFilterCount` getters
- `TransactionSortField`, `TransactionSortOrder` enums

### Repository (`transaction_repository.dart`)

Follow `features/timeline/repository/timeline_repository.dart`.

- Provider: `Provider<TransactionRepository>` with `ApiClient`
- `getTransactions({userId, cursor, limit, filters})` — `GET /api/v1/transactions`
- `getTransaction({userId, id})` — `GET /api/v1/transactions/{id}`
- `searchTransactions({userId, query, limit})` — `GET /api/v1/transactions/search?q=`
- `getSummary({userId})` — `GET /api/v1/transactions/summary`
- `createTransaction({userId, data})` — `POST /api/v1/events`
- `archiveTransaction({userId, id})` — `POST /api/v1/events/{id}/archive`

### State (`transaction_state.dart`)

Follow `features/timeline/state/timeline_state.dart` EXACTLY.

- `TransactionListStatus` enum: `initial, loading, loaded, loadingMore, error, empty, offline`
- `TransactionListState` with `copyWith`: `status`, `transactions[]`, `cursor`, `hasMore`, `error`, `activeFilters`, `searchQuery`
- `TransactionListNotifier`:
  - `load({userId})` — initial load, sets status to `loading`, handles empty
  - `loadMore({userId})` — appends items, guarded by `_isLoading` and `hasMore`
  - `refresh({userId})` — reset cursor, reload from page 1
  - `search({userId, query})` — search with dedicated search endpoint
  - `applyFilters({userId, filters})` — filter with backend params
  - `clearFilters({userId})` — reset to unfiltered

### Pages

#### Transaction List (`transactions_page.dart`)

Follow `features/timeline/pages/timeline_page.dart` EXACTLY.

**Structure:**
- `ConsumerStatefulWidget` with `ScrollController` for infinite scroll
- `Future.microtask(() => ref.read(transactionListProvider.notifier).load())` in `initState`
- `ScrollController` listener triggers `loadMore()` when within 200px of bottom

**States:**
- `loading`/`initial`: Skeleton list of 8 shimmer cards
- `loaded`: `RefreshIndicator` wrapping `ListView.builder` of `TransactionCard` widgets
- `loadingMore`: Overlay loading spinner at bottom
- `empty`: Icon + message. Contextual: "No transactions found" vs "No transactions yet" vs "Try different search/filters" with clear button
- `error`: Retry button calling `refresh()`
- `offline`: `wifi_off` icon + message

**AppBar:**
- Title: "Transactions"
- Search toggle button — shows `TextField` with autofocus, 300ms debounce via `Timer`, cancel button
- Filter button — shows `Badge` with `activeFilterCount` when filters active

**FAB:**
- `FloatingActionButton` with `Icons.add`
- Navigates to `/transactions/add`
- On pop with `true` result, calls `refresh()`

#### Transaction Detail (`transaction_detail_page.dart`)

Follow `features/goals/pages/goal_detail_page.dart` pattern with `FutureProvider.autoDispose.family`.

**Sections (each in a Card):**
1. **Header Card**: Amount (large, colored red/green by expense/income), currency, description, status badge
2. **Details Card**: Category, date, type, source, destination, reference, origin, confidence
3. **Notes Card**: Notes text or "No notes added" in italic
4. **Attachments Card**: "Not yet supported" placeholder
5. **Tags Card**: Chip list or "No tags" in italic

**AppBar actions:**
- `PopupMenuButton` with "Archive" option
- Archive calls `archiveTransaction()`, shows success snackbar, pops back

#### Transaction Form (`transaction_form_page.dart`)

Follow `features/settings/pages/edit_profile_page.dart`.

**Fields:**
- Type: `SegmentedButton<String>` with "Expense" (arrow_upward) and "Income" (arrow_downward)
- Amount: `TextFormField` with ₹ prefix, keyboard type decimal, validate > 0
- Description: `TextFormField`, maxLines: 2
- Category: `DropdownButtonFormField` — changes options based on type (expense: purchase/bill/food/transport/entertainment/utilities/rent/healthcare/other, income: salary/freelance/investment/refund/gift/other)

**Behavior:**
- Save button in AppBar shows loading spinner while saving
- On success: SnackBar "Transaction created", pop with `true`
- On error: SnackBar with error message

**Footer note:** "Backend currently supports limited fields for new transactions."

### Widgets

#### `TransactionCard` (in `widgets/transaction_widgets.dart`)
- Colored icon (red container for expense, green for income)
- Description + category + relative date in subtitle
- Amount with sign and currency
- `InkWell` with `onTap` for navigation

#### `TransactionSummaryCard`
- Monthly income (green, arrow_down) / expenses (red, arrow_up)
- Net flow row with divider

#### `TransactionStatusBadge`
- Colored chip: posted=green, pending=orange, draft=grey, cancelled=red, reversed=red, archived=grey, confirmed=green

#### `TransactionFilterSheet` (in `widgets/transaction_filters.dart`)
- `showModalBottomSheet` with date range pickers, type chips, amount range, Clear All button, Apply button

### Navigation (`lib/app/router.dart`)

- Add import for transaction pages
- Add routes: `/transactions` -> `TransactionsPage`, `/transactions/:id` -> `TransactionDetailPage`, `/transactions/add` -> `TransactionFormPage`
- Add `ListTile` in `_MorePage` with `Icons.receipt_long_outlined` and title "Transactions"

---

## VERIFICATION CHECKLIST

- [ ] `flutter analyze` — zero errors
- [ ] Navigation: More page -> Transactions list
- [ ] Pull-to-refresh loads data (test with backend running or error state gracefully)
- [ ] Infinite scroll triggers `loadMore` (verify `ScrollController` listener fires)
- [ ] Search bar toggles, debounce fires after 300ms
- [ ] Filter sheet opens with date/type/amount controls
- [ ] FAB opens add form
- [ ] Add form validates amount > 0
- [ ] Add form saves and pops back
- [ ] Detail page shows amount, category, status, notes, tags
- [ ] Archive action from detail works
- [ ] Empty state shows when no transactions
- [ ] Error state shows on network failure
- [ ] Offline state shows when connectivity lost
- [ ] All states render without crashing: loading, loaded, empty, error, offline

---

## DELIVERABLE TEMPLATE

At the end, provide:

```markdown
## 1. Folder Structure
...
## 2. API Inventory Used
...
## 3. Screens Implemented
...
## 4. Components Reused
...
## 5. Files Modified
...
## 6. Remaining Backend Gaps
...
## 7. Verification Checklist (all checked)
...
## 8. Message to Solution Architect
...
```
