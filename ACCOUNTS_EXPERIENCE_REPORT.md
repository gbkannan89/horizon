# ACCOUNTS EXPERIENCE REPORT

**Spec ID:** HZN-EXP-006
**Phase:** 6.6
**Status:** Complete
**Build:** 37/37 modules PASS

---

## Experience Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                     REST API (10 endpoints)                           │
├──────────────────────────────────────────────────────────────────────┤
│                     Handlers (thin layer)                             │
├──────────────────────────────────────────────────────────────────────┤
│                    Aggregator (9 providers, parallel)                 │
│  ┌──────┬──────┬──────────┬──────┬──────┬──────┬──────┬──────┬──────┐│
│  │Accts │ Trans│Projection│ Risk │Health│ Rec  │ Opt  │  Sim │Events││
│  └──────┴──────┴──────────┴──────┴──────┴──────┴──────┴──────┴──────┘│
├──────────────────────────────────────────────────────────────────────┤
│                    Composer + Engine (view models)                    │
│  ┌──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────────┐│
│  │Dashbd│Summary│Bal.  │Cash  │Health│ Rec  │ Proj │Detail│ Timeline ││
│  │      │(list)│      │Flow  │      │      │      │      │          ││
│  └──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────────┘│
├──────────────────────────────────────────────────────────────────────┤
│              InMemoryCache (read-only, TTL=5m)                        │
└──────────────────────────────────────────────────────────────────────┘
```

### Layers

| Layer | Files | Purpose |
|---|---|---|
| **Config** | `internal/config/config.go` | Environment configuration (PORT=8100) |
| **Engine** | `internal/engine/core.go` | AccountsDashboard, BalanceSummary, CashFlowSummary, AccountHealthSummary, 10 card types |
| **Engine** | `internal/engine/composer.go` | BuildAccountsList, BuildAccountDetail, filtering, pagination |
| **Engine** | `internal/engine/types.go` | AccountCardData, BalanceEntry, AccountDetailData, TransPreview, AccountsList, Inputs |
| **Aggregator** | `internal/aggregator/aggregator.go` | 9 providers, 9-way parallel, account detail enrichment |
| **Cache** | `internal/infrastructure/cache/cache.go` | In-memory read-through cache |
| **API** | `internal/api/handlers.go`, `internal/api/router.go`, `internal/api/dto.go` | 10 REST endpoints with DTOs |
| **Main** | `cmd/server/main.go` | Dependency injection, 9 no-op providers |

---

## Accounts Dashboard Model

### 10 Card Types

| Card | CardType | Data Source |
|---|---|---|
| Accounts Overview | `Overview` | Account Domain |
| Cash Balance | `CashBalance` | Account Domain |
| Account Health | `AccountHealth` | Health Score Engine |
| Cash Flow | `CashFlow` | Financial Event Domain |
| Risk Assessment | `Risk` | Risk Engine |
| Projection | `Projection` | Projection Engine |
| Recommendations | `Recommendation` | Recommendation Engine |
| Timeline | `Timeline` | Financial Event Domain |

### View Models

| Model | Fields | Purpose |
|---|---|---|
| `AccountsDashboard` | TotalBalance, TotalAccounts, AccountsByType, Cards[8] | Dashboard view |
| `AccountsList` | Accounts[AccountCardData], Total, Cursor, HasMore, Summary | Paginated account list |
| `AccountDetailData` | All account fields, Balances[7], Transactions[] | Single account detail |
| `BalanceSummary` | TotalBalance, TotalAvailable, TotalSpendable, CreditUtilized, AccountBalances | Balance breakdown |
| `CashFlowSummary` | PeriodInflow, PeriodOutflow, NetFlow, Projected | Cash flow summary |
| `AccountHealthSummary` | AccountsByHealth map, HealthyPct | Health aggregation |

### Balance Types (7)

`current`, `available`, `cleared`, `uncleared`, `reserved`, `allocated`, `spendable`

---

## REST Endpoints

| Method | Path | Response | Description |
|---|---|---|---|
| GET | `/api/v1/accounts/experience` | AccountsDashboard | Full accounts dashboard |
| GET | `/api/v1/accounts/dashboard` | AccountsDashboard | Same as experience |
| GET | `/api/v1/accounts/summary` | AccountsList | Paginated account list with filters |
| GET | `/api/v1/accounts/balances` | BalanceSummary | Balance breakdown |
| GET | `/api/v1/accounts/cashflow` | CashFlowSummary | Cash flow summary |
| GET | `/api/v1/accounts/health` | AccountHealthSummary | Health aggregation |
| GET | `/api/v1/accounts/recommendations` | Card[] | Active recommendations |
| GET | `/api/v1/accounts/projections` | Card | Cash flow projection |
| GET | `/api/v1/accounts/timeline` | Card | Recent activity |
| GET | `/api/v1/accounts/{id}` | AccountDetailData | Single account with transactions |

### Query Parameters

| Parameter | Type | Description |
|---|---|---|
| `user_id` | string | User identifier |
| `type` | []string | Filter by account type(s) |
| `status` | []string | Filter by status(es) |
| `currency` | []string | Filter by currency(ies) |
| `cursor` | string | Pagination cursor |
| `limit` | int | Page size (default: 25) |

---

## Aggregation Strategy

```
9-way parallel fan-out via goroutines:
  ├─ Accounts.GetAccounts              → account list
  ├─ Projection.GetCashFlowProjection  → inflow, outflow, projected
  ├─ Risk.GetRiskScore                 → score, level
  ├─ Health.GetHealthScore             → score, grade
  ├─ Recs.HasRecommendations           → has, count
  ├─ Optimize.HasOptimizations         → has flag
  ├─ Simulation.HasSimulations         → has flag
  └─ Events.GetEventCount              → event count

Post-aggregation: cash balance = sum of all account current balances
```

### Account Detail Flow
```
1. Full aggregation (cache hit possible)
2. Find specific account in list
3. Fetch transactions from TransProvider
4. Compose AccountDetailData via BuildAccountDetail
```

---

## Engine Integrations

| Engine | Provider Interface | Composer Method | Endpoint |
|---|---|---|---|
| **Projection Engine** | `ProjProvider.GetCashFlowProjection` | `BuildCashFlowSummary` | `/accounts/projections` |
| **Risk Engine** | `RiskProvider.GetRiskScore` | `BuildDashboard` | `/accounts/risk` |
| **Health Score Engine** | `HealthProvider.GetHealthScore` | `BuildHealthSummary` | `/accounts/health` |
| **Recommendation Engine** | `RecProvider.HasRecommendations` | Dashboard card | `/accounts/recommendations` |
| **Optimization Engine** | `OptProvider.HasOptimizations` | Dashboard flag | `/accounts/dashboard` |
| **Simulation Engine** | `SimProvider.HasSimulations` | Dashboard flag | `/accounts/dashboard` |

All engine integrations delegate entirely — the Accounts Experience contains zero financial calculation logic.

---

## Architecture Compliance

| Rule | Status |
|---|---|
| Read-only architecture | ✅ All providers are read-only. No writes. |
| No business rules | ✅ Composition, aggregation, filtering only. |
| No financial calculations | ✅ All financial values from domain/engine outputs. |
| Delegates to engines | ✅ 6 engines integrated. |
| No domain writes | ✅ Accounts Experience never calls Save/Create/Update/Delete. |

---

## Build Status

| Check | Status |
|---|---|
| `go build ./...` (37 modules) | ✅ PASS |
| `go vet ./...` (37 modules) | ✅ PASS |
| Accounts endpoints | ✅ 10 endpoints registered |
| Engine integrations | ✅ 6 engines integrated |
| Account detail with 7 balance types | ✅ |
| Cursor pagination | ✅ |
| Account filtering | ✅ by type, status, currency |
| Read-only architecture | ✅ Verified |

---

## File Inventory

| File | Lines | Purpose |
|---|---|---|
| `cmd/server/main.go` | 72 | Entry point, DI, 9 no-op providers |
| `internal/config/config.go` | 16 | Environment configuration |
| `internal/engine/core.go` | 195 | AccountsDashboard, 10 card types, Composer with 6 build methods |
| `internal/engine/composer.go` | 136 | BuildAccountsList, BuildAccountDetail, filtering |
| `internal/engine/types.go` | 81 | All data types (AccountCardData, AccountDetailData, etc.) |
| `internal/aggregator/aggregator.go` | 147 | 9 providers, parallel aggregation, account detail |
| `internal/infrastructure/cache/cache.go` | 44 | In-memory read-through cache |
| `internal/api/handlers.go` | 212 | 10 REST endpoint handlers |
| `internal/api/router.go` | 42 | Router + middleware |
| `internal/api/dto.go` | 80 | Response DTOs (7 response types) |
| `internal/events/events.go` | 18 | AccountListViewShown, AccountDetailViewShown, AccountCreated, AccountSearchPerformed |

---

## Remaining Experience Layer Work

| Phase | Experience | Status |
|---|---|---|
| 6.1 | Dashboard | ✅ Complete |
| 6.2 | Timeline | ✅ Complete |
| 6.3 | Goals | ✅ Complete |
| 6.4 | Portfolio | ✅ Complete |
| 6.5 | Planning | ✅ Complete |
| **6.6** | **Accounts** | **✅ Complete** |
| 6.7 | Insights | 🔲 Not started |
| 6.8 | Notifications | 🔲 Not started |
| 6.9 | Advisor | 🔲 Not started |

---

## Ready for Phase 6.7 — Insights Experience

```
Go build:   37/37 PASS
Go vet:     37/37 PASS
Endpoints:  10/10 operational
Card types: 10/10 implemented
Engines:    6/6   integrated
Providers:  9/9   data providers
Architecture: READ-ONLY ✅
```
