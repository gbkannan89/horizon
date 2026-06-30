# GOALS UI REPORT

**Phase:** 8.5
**Status:** Complete
**Build:** `flutter analyze`: 0 errors | `flutter build apk --debug`: SUCCESS

---

## Screen Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                      Goals List Screen                                │
├──────────────────────────────────────────────────────────────────────┤
│  AppBar: "Goals" + Search toggle                                     │
├──────────────────────────────────────────────────────────────────────┤
│  Goal Cards (pull-to-refresh list)                                   │
│  ┌──────────────────────────────────────────────────────────────────┐│
│  │ GoalListCard                                                     ││
│  │  [Goal Name]                              [Status Badge]        ││
│  │  [Importance] [Priority: N]                     [Rec Icon]      ││
│  │  ──────────────────────────────────────────────────────────────── ││
│  │  [████████████████░░░░░░░░░░░]                          75%      ││
│  └──────────────────────────────────────────────────────────────────┘│
├──────────────────────────────────────────────────────────────────────┤
│  Goal Detail Screen (push on tap)                                    │
└──────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────┐
│                      Goal Detail Screen                               │
├──────────────────────────────────────────────────────────────────────┤
│  AppBar: Goal Name                                                   │
├──────────────────────────────────────────────────────────────────────┤
│  GoalHeaderCard                                                      │
│  ┌──────────────────────────────────────────────────────────────────┐│
│  │ [Flag Icon] Goal Name                     [Status Badge]        ││
│  │ ──────────────────────────────────────────────────────────────── ││
│  │ ₹5,00,000                            75%                        ││
│  │ of ₹10,00,000                        Dec 2030                   ││
│  │ [███████████████████████                                           ││
│  │ ⚠ Gap: ₹2,00,000                                                  ││
│  └──────────────────────────────────────────────────────────────────┘│
├──────────────────────────────────────────────────────────────────────┤
│  Dashboard Cards (8 types from API)                                  │
│  ┌──────────────────────┐ ┌──────────────────────┐                   │
│  │ Overview             │ │ Progress             │                   │
│  ├──────────────────────┤ ├──────────────────────┤                   │
│  │ Funding              │ │ Projection           │                   │
│  ├──────────────────────┤ ├──────────────────────┤                   │
│  │ Recommendation       │ │ Optimization         │                   │
│  ├──────────────────────┤ ├──────────────────────┤                   │
│  │ Timeline             │ │ Milestone            │                   │
│  └──────────────────────┘ └──────────────────────┘                   │
├──────────────────────────────────────────────────────────────────────┤
│  Projection Section                                                  │
│  Recommendation Section                                              │
│  Optimization Section                                                │
│  Timeline Section                                                    │
│  Milestone Section                                                   │
└──────────────────────────────────────────────────────────────────────┘
```

---

## Widget Hierarchy

| Widget | Purpose |
|---|---|
| `GoalListCard` | Compact card with name, status badge, importance, progress bar |
| `GoalHeaderCard` | Full goal header with icon, progress, amount, target date, funding gap |
| `GoalDashboardCard` | Reusable card with icon, title, summary, value (8 card types) |
| `GoalDetailPage` | Full detail with dashboard, projection, recs, opts, timeline, milestones |

---

## API Integration (8 endpoints)

| Endpoint | State Method | Section |
|---|---|---|
| `GET /goals/experience` | `load()` | List screen |
| `GET /goals/:id/dashboard` | `load()` (parallel) | Header + 8 cards |
| `GET /goals/:id/progress` | `load()` (parallel) | Progress section |
| `GET /goals/:id/projection` | `load()` (parallel) | Projection section |
| `GET /goals/:id/recommendations` | `load()` (parallel) | Recs section |
| `GET /goals/:id/optimization` | `load()` (parallel) | Opts section |
| `GET /goals/:id/timeline` | `load()` (parallel) | Timeline section |
| `GET /goals/:id/milestones` | `load()` (parallel) | Milestones section |

Detail screen loads all 7 goal-specific endpoints in parallel via `Future.wait`.

---

## State Management

### Goals List

| State | Visual | Action |
|---|---|---|
| `loading: true` | CircularProgressIndicator | Auto-loads |
| `goals: []` | Empty state (icon + message) | — |
| `error` | Error + retry button | Tap retry |
| `loaded` | Goal cards list | Scroll, tap card |

### Goal Detail

Uses `StateNotifierProvider.family` for per-goal state. 7 endpoints loaded in parallel.

---

## Goal Importance Mapping

| Importance | Color |
|---|---|
| Mandatory | red |
| Essential | orange |
| Lifestyle | blue |
| Dream | purple |

## Goal Status Mapping

| Status | Color |
|---|---|
| Active | green |
| Paused | orange |
| Completed | blue |
| AtRisk | red |

## Goal Card Types (8)

| Card | Icon | Color |
|---|---|---|
| Overview | `info_outline` | blue |
| Progress | `pie_chart` | green |
| Funding | `account_balance_wallet` | teal |
| Projection | `trending_up` | cyan |
| Recommendation | `lightbulb` | amber |
| Optimization | `auto_graph` | purple |
| Timeline | `history` | indigo |
| Milestone | `flag` | orange |

---

## Performance

| Mechanism | Implementation |
|---|---|
| **Parallel loading** | `Future.wait` for all 7 detail endpoints |
| **Pull-to-refresh** | `RefreshIndicator` on both screens |
| **Search** | Client-side filter on goals list |
| **Lazy rendering** | `ListView.builder` |
| **Skeleton loading** | Circular progress for simplicity |

---

## Build Status

| Check | Status |
|---|---|
| `flutter analyze` | ✅ **0 errors** |
| `flutter build apk --debug` | ✅ **SUCCESS** |
| Goals list | ✅ Cards, search, filter |
| Goal detail | ✅ Header, 8 cards, projection, recs, opts, timeline, milestones |
| API integration | ✅ 8 endpoints (7 parallel on detail) |
| Pull-to-refresh | ✅ Both screens |
| Status/importance colors | ✅ Mapped |

---

## File Inventory

| File | Lines | Purpose |
|---|---|---|
| `features/goals/models/goal_models.dart` | 310 | All response models (10 classes) |
| `features/goals/repository/goal_repository.dart` | 85 | 8 endpoint methods |
| `features/goals/widgets/goal_widgets.dart` | 200 | GoalListCard, GoalHeaderCard, GoalDashboardCard |
| `features/goals/pages/goals_page.dart` | 105 | Goals list with search |
| `features/goals/pages/goal_detail_page.dart` | 235 | Goal detail with 7 API sections |

---

## Remaining Flutter Work

| Phase | Feature | Status |
|---|---|---|
| 8.1 | Foundation | ✅ Complete |
| 8.2 | Authentication UI | ✅ Complete |
| 8.3 | Dashboard UI | ✅ Complete |
| 8.4 | Timeline UI | ✅ Complete |
| **8.5** | **Goals UI** | **✅ Complete** |
| 8.6 | Portfolio UI | 🔲 |
| 8.7 | Accounts UI | 🔲 |
| 8.8 | Settings/Profile | 🔲 |

---

## Ready for Phase 8.6 — Portfolio UI

```
Flutter analyze: 0 errors
APK build:       SUCCESS
Endpoints:       8 integrated
Screens:         2 (list + detail)
Card types:      8 mapped
Parallel loads:  7 endpoints via Future.wait
```
