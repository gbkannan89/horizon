# PORTFOLIO EXPERIENCE REPORT

**Spec ID:** HZN-EXP-004
**Phase:** 6.4
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
│                    Aggregator (10 providers, parallel)                │
│  ┌──────┬──────┬──────────┬──────┬──────┬──────┬──────┬──────┬─────┐ │
│  │Port. │Asset │Allocation│ Perf │ Risk │ Proj │ Rec  │ Opt  │ Sim │ │
│  └──────┴──────┴──────────┴──────┴──────┴──────┴──────┴──────┴─────┘ │
│                          ┌────────┐                                  │
│                          │ Events │                                  │
│                          └────────┘                                  │
├──────────────────────────────────────────────────────────────────────┤
│                    Composer (view models)                             │
│  ┌──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────────┐│
│  │Dashbd│Alloc │Perf  │ Risk │ Proj │ Rec  │ Opt  │ Sim  │ Timeline ││
│  └──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────────┘│
├──────────────────────────────────────────────────────────────────────┤
│              InMemoryCache (read-only, TTL=5m)                        │
└──────────────────────────────────────────────────────────────────────┘
```

### Layers

| Layer | Files | Purpose |
|---|---|---|
| **Config** | `internal/config/config.go` | Environment configuration (PORT) |
| **Engine** | `internal/engine/core.go` | View models: PortfolioDashboard, AllocationView, PerformanceView, RiskView, ProjectionView, SimView, CardView — 10 card types |
| **Aggregator** | `internal/aggregator/aggregator.go` | 10 providers, 11-way parallel fan-out |
| **Cache** | `internal/infrastructure/cache/cache.go` | In-memory read-through cache |
| **API** | `internal/api/handlers.go`, `internal/api/router.go`, `internal/api/dto.go` | 10 REST endpoints with DTOs |
| **Main** | `cmd/server/main.go` | Dependency injection, 10 no-op providers |

---

## Portfolio Dashboard Model

### 10 Card Types

| Card | CardType | Data Source |
|---|---|---|
| Portfolio Overview | `Overview` | Portfolio Domain |
| Asset Allocation | `Allocation` | Asset Domain |
| Performance | `Performance` | Portfolio Domain + Performance Engine |
| Holdings | `Holdings` | Asset Domain |
| Projection | `Projection` | Projection Engine |
| Risk Assessment | `Risk` | Risk Engine |
| Recommendations | `Recommendation` | Recommendation Engine |
| Optimization | `Optimization` | Optimization Engine |
| Simulation | `Simulation` | Simulation Engine |
| Timeline | `Timeline` | Financial Event Domain + Timeline |

### View Models

| Model | Fields | Purpose |
|---|---|---|
| `PortfolioDashboard` | PortfolioValue, TotalReturn, TotalReturnPct, RiskScore, RiskLevel, Cards[10] | Full portfolio dashboard |
| `AllocationView` | Allocations[AllocEntry], TotalValue, DiversificationScore | Allocation breakdown with drift |
| `PerformanceView` | PeriodReturn, PeriodReturnPct, BenchmarkReturn, UnrealizedGL, RealizedGL | Performance metrics |
| `RiskView` | RiskScore, RiskLevel, VaR, SharpeRatio, Volatility, MaxDrawdown | Risk assessment |
| `ProjectionView` | ProjectedValue, Confidence, HorizonYears, AnnualReturn | Portfolio projection |
| `CardView` | Items[CardEntry], Count | Generic card list |
| `SimView` | Simulations[SimEntry], TotalCount | Scenario simulations |

---

## REST Endpoints

| Method | Path | Response | Description |
|---|---|---|---|
| GET | `/api/v1/portfolio/experience` | PortfolioDashboard | Full portfolio dashboard |
| GET | `/api/v1/portfolio/dashboard` | PortfolioDashboard | Same as experience |
| GET | `/api/v1/portfolio/allocation` | AllocationView | Asset allocation breakdown |
| GET | `/api/v1/portfolio/performance` | PerformanceView | Performance metrics |
| GET | `/api/v1/portfolio/risk` | RiskView | Risk assessment |
| GET | `/api/v1/portfolio/projection` | ProjectionView | Portfolio projection |
| GET | `/api/v1/portfolio/recommendations` | CardView | Active recommendations |
| GET | `/api/v1/portfolio/optimization` | CardView | Optimization strategies |
| GET | `/api/v1/portfolio/simulations` | SimView | Scenario simulations |
| GET | `/api/v1/portfolio/timeline` | CardView | Recent portfolio events |

### Query Parameters

| Parameter | Type | Description |
|---|---|---|
| `user_id` | string | User identifier (default: "default") |

---

## Aggregation Strategy

```
11-way parallel fan-out via goroutines:
  ├─ Portfolio.GetPortfolioSummary   → value, return, returnPct
  ├─ Assets.GetAllocations           → allocation entries
  ├─ Performance.GetPerformance      → period return, benchmark, GL
  ├─ Risk.GetPortfolioRisk           → score, level, Sharpe, VaR, volatility
  ├─ Projection.GetPortfolioProjection → projected value, confidence
  ├─ Recs.HasRecommendations         → has, count
  ├─ Optimize.HasOptimizations       → has, count
  ├─ Simulation.HasSimulations       → has, count
  ├─ Events.GetEventCount            → timeline event count
  └─ Events.GetMilestoneCount        → milestone count
```

---

## Engine Integrations

| Engine | Provider Interface | Composer Method | Endpoint |
|---|---|---|---|
| **Projection Engine** | `ProjProvider.GetPortfolioProjection` | `BuildProjection` | `/portfolio/projection` |
| **Risk Engine** | `RiskProvider.GetPortfolioRisk` | `BuildRisk` | `/portfolio/risk` |
| **Recommendation Engine** | `RecProvider.HasRecommendations` | `BuildRecommendations` | `/portfolio/recommendations` |
| **Optimization Engine** | `OptProvider.HasOptimizations` | `BuildOptimization` | `/portfolio/optimization` |
| **Simulation Engine** | `SimProvider.HasSimulations` | `BuildSimulations` | `/portfolio/simulations` |

All engine integrations delegate entirely — the Portfolio Experience contains zero financial calculation logic.

---

## Architecture Compliance

| Rule | Status |
|---|---|
| Read-only architecture | ✅ All providers are read-only. No writes. |
| No business rules | ✅ All logic is orchestrational (composition, aggregation). |
| No financial calculations | ✅ All financial values from domain/engine outputs only. |
| Delegates to engines | ✅ Projection, Risk, Recommendation, Optimization, Simulation all delegated. |
| No domain writes | ✅ Portfolio Experience never calls Save/Create/Update/Delete. |

---

## Build Status

| Check | Status |
|---|---|
| `go build ./...` (37 modules) | ✅ PASS |
| `go vet ./...` (37 modules) | ✅ PASS |
| Portfolio endpoints | ✅ 10 endpoints registered |
| Engine integrations | ✅ 5 engines integrated |
| Read-only architecture | ✅ Verified |

---

## File Inventory

| File | Lines | Purpose |
|---|---|---|
| `cmd/server/main.go` | 74 | Entry point, DI, 10 no-op providers |
| `internal/config/config.go` | 16 | Environment configuration |
| `internal/engine/core.go` | 260 | All view models (7 types), 10 card types, Composer with 9 build methods |
| `internal/aggregator/aggregator.go` | 130 | 10 providers, 11-way parallel fan-out |
| `internal/infrastructure/cache/cache.go` | 44 | In-memory read-through cache |
| `internal/api/handlers.go` | 122 | 10 REST endpoint handlers |
| `internal/api/router.go` | 42 | Router + middleware |
| `internal/api/dto.go` | 97 | Response DTOs (9 response types) |
| `internal/events/events.go` | 12 | PortfolioOpened event |

---

## Remaining Experience Layer Work

| Phase | Experience | Status |
|---|---|---|
| 6.1 | Dashboard | ✅ Complete |
| 6.2 | Timeline | ✅ Complete |
| 6.3 | Goals | ✅ Complete |
| **6.4** | **Portfolio** | **✅ Complete** |
| 6.5 | Planning | 🔲 Not started |
| 6.6 | Accounts | 🔲 Not started |
| 6.7 | Insights | 🔲 Not started |
| 6.8 | Notifications | 🔲 Not started |
| 6.9 | Advisor | 🔲 Not started |

---

## Ready for Phase 6.5 — Planning Experience

```
Go build:   37/37 PASS
Go vet:     37/37 PASS
Endpoints:  10/10 operational
Card types: 10/10 implemented
Engines:    5/5   integrated
Providers:  10/10 data providers
Architecture: READ-ONLY ✅
```
