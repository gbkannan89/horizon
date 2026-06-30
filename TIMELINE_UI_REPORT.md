# TIMELINE UI REPORT

**Phase:** 8.4
**Status:** Complete
**Build:** `flutter analyze`: 0 errors | `flutter build apk --debug`: SUCCESS

---

## Screen Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                        Timeline Screen                                │
├──────────────────────────────────────────────────────────────────────┤
│  AppBar: "Timeline" + Search toggle + Filter button                  │
├──────────────────────────────────────────────────────────────────────┤
│  Search Mode: TextField with submit → search endpoint                │
├──────────────────────────────────────────────────────────────────────┤
│  Timeline Feed (infinite scroll via ListView.builder)                │
│  ┌──────────────────────────────────────────────────────────────────┐│
│  │ TimelineCard                                                     ││
│  │  ┌──────┐ ┌────────────────────────────────────────────────────┐ ││
│  │  │ Icon │ │ Title + Timestamp                                  │ ││
│  │  │ 40x40│ │ Summary (2 lines max)                              │ ││
│  │  │      │ │ [CategoryBadge] [SeverityBadge]         [Amount]   │ ││
│  │  └──────┘ └────────────────────────────────────────────────────┘ ││
│  └──────────────────────────────────────────────────────────────────┘│
│  ... (infinite scroll loads more via cursor pagination)              │
├──────────────────────────────────────────────────────────────────────┤
│  Empty State: Icon + message                                         │
│  Error State: Icon + retry button                                    │
│  Loading State: 8 skeleton cards                                     │
└──────────────────────────────────────────────────────────────────────┘
```

## Detail Screen

```
┌──────────────────────────────────────────────────────────────────────┐
│                      Timeline Detail Screen                           │
├──────────────────────────────────────────────────────────────────────┤
│  AppBar: "Event Details" + back button                               │
├──────────────────────────────────────────────────────────────────────┤
│  TimelineDetailCard                                                  │
│  ┌──────────────────────────────────────────────────────────────────┐│
│  │ [Icon] Title + EventType                                        ││
│  │ ──────────────────────────────────────────────────────────────── ││
│  │ Description (if present)                                        ││
│  │ Category, Severity, EventType, Timestamp                         ││
│  │ Related Entity, Goal, Account, Asset (if present)                ││
│  │ Amount (if non-zero)                                             ││
│  │ Metadata (if present, key-value list)                            ││
│  └──────────────────────────────────────────────────────────────────┘│
│  [Back to Timeline] button                                           │
└──────────────────────────────────────────────────────────────────────┘
```

---

## Widget Hierarchy

| Widget | Purpose |
|---|---|
| `TimelineCard` | Compact event card with icon, title, summary, badges, amount |
| `TimelineDetailCard` | Full event detail with all fields |
| `CategoryBadge` | Colored chip showing event category |
| `SeverityBadge` | Colored chip showing severity level |
| `FilterSheet` | Bottom sheet with category, severity, date range filters |
| `categoryIcon()` | Maps category string to Material icon |
| `categoryColor()` | Maps category string to color |
| `severityColor()` | Maps severity string to color |

---

## Icon/Color Mapping

| Category | Icon | Color |
|---|---|---|
| FinancialEvent | `currency_rupee` | green |
| GoalEvent | `flag` | blue |
| AccountEvent | `account_balance` | indigo |
| AssetEvent | `trending_up` | teal |
| LiabilityEvent | `credit_card` | purple |
| PortfolioEvent | `pie_chart` | blue |
| HealthEvent | `favorite` | pink |
| RiskEvent | `shield` | amber |
| RecommendationEvent | `lightbulb` | amber |
| SimulationEvent | `science` | cyan |
| OptimizationEvent | `auto_graph` | cyan |
| Achievement | `emoji_events` | gold |
| Milestone | `flag` | green |
| UserEvent | `person` | grey |

## API Integration

| Endpoint | State Method | Used By |
|---|---|---|
| `GET /api/v1/timeline` | `load()`, `loadMore()` | Initial load + pagination |
| `GET /api/v1/timeline/recent` | (available) | Recent events |
| `GET /api/v1/timeline/filter` | `applyFilters()` | Category/severity/date filter |
| `GET /api/v1/timeline/search` | `search()` | Search bar |
| `GET /api/v1/timeline/:id` | `getById()` | Detail screen |

---

## Pagination Strategy

| Mechanism | Implementation |
|---|---|
| **Infinite scroll** | `ScrollController` detects when near bottom |
| **Cursor pagination** | Backend returns `cursor` + `has_more` |
| **Lazy rendering** | `ListView.builder` with visible-only items |
| **Loading more** | `TimelineStatus.loadingMore` with bottom indicator |
| **Pull-to-refresh** | `RefreshIndicator` → `refresh()` |

---

## Filtering

| Filter | Type | UI |
|---|---|---|
| Categories | Multi-select chips | FilterSheet |
| Severity | Single-select chips | FilterSheet |
| Date Range | Text fields (YYYY-MM-DD) | FilterSheet |
| Search | Text field | AppBar search mode |

---

## State Management

| State | Visual | User Action |
|---|---|---|
| `initial` | — | Auto-loads |
| `loading` | 8 skeleton cards | — |
| `loaded` | Timeline feed | Scroll, tap cards |
| `loadingMore` | Feed + bottom spinner | Auto on scroll |
| `empty` | Empty state icon + message | Adjust filters |
| `error` | Error + retry button | Tap retry |
| `offline` | Offline message | Connect to internet |

---

## Performance

| Mechanism | Implementation |
|---|---|
| **Infinite scroll** | `ScrollController` + threshold detection |
| **Cursor pagination** | No offset-based scans |
| **Lazy rendering** | `ListView.builder` |
| **Pull-to-refresh** | `RefreshIndicator` |
| **Skeleton loading** | 8 animated placeholder cards |
| **Optimistic UI** | Keeps items while loading more |

---

## Build Status

| Check | Status |
|---|---|
| `flutter analyze` | ✅ **0 errors** |
| `flutter build apk --debug` | ✅ **SUCCESS** |
| Timeline feed | ✅ Infinite scroll with cursor pagination |
| Timeline detail | ✅ Full event detail with all fields |
| Search | ✅ Text search via `/timeline/search` |
| Filtering | ✅ Categories, severity, date range |
| Pull-to-refresh | ✅ `RefreshIndicator` |
| Icon/color mapping | ✅ 14 categories mapped |
| Loading skeleton | ✅ 8 placeholder cards |
| Empty/Error states | ✅ Handled |

---

## File Inventory

| File | Lines | Purpose |
|---|---|---|
| `features/timeline/models/timeline_models.dart` | 130 | TimelineItem, TimelineFeed, TimelineOutput, FilterParams |
| `features/timeline/repository/timeline_repository.dart` | 75 | 5 endpoint methods |
| `features/timeline/state/timeline_state.dart` | 103 | TimelineStateNotifier with 7 states + pagination |
| `features/timeline/widgets/timeline_widgets.dart` | 370 | TimelineCard, TimelineDetailCard, CategoryBadge, SeverityBadge, FilterSheet, icon/color mappers |
| `features/timeline/pages/timeline_page.dart` | 192 | Main timeline screen |
| `features/timeline/pages/timeline_detail_page.dart` | 78 | Event detail screen |

---

## Remaining Flutter Work

| Phase | Feature | Status |
|---|---|---|
| 8.1 | Foundation | ✅ Complete |
| 8.2 | Authentication UI | ✅ Complete |
| 8.3 | Dashboard UI | ✅ Complete |
| **8.4** | **Timeline UI** | **✅ Complete** |
| 8.5 | Goals UI | 🔲 |
| 8.6 | Portfolio UI | 🔲 |
| 8.7 | Accounts UI | 🔲 |
| 8.8 | Settings/Profile | 🔲 |

---

## Ready for Phase 8.5 — Goals UI

```
Flutter analyze: 0 errors
APK build:       SUCCESS
Timeline states: 7 states
Widget cards:    5 widgets (+ icons/badges)
API endpoints:   5 integrated
Filter types:    4 categories
```
