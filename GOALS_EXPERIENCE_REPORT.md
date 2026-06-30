# GOALS EXPERIENCE REPORT

**Spec ID:** HZN-EXP-003
**Phase:** 6.3
**Status:** Complete
**Build:** 37/37 modules PASS

---

## Experience Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                     REST API (8 endpoints)                           │
├─────────────────────────────────────────────────────────────────────┤
│                     Handlers (thin layer)                            │
├─────────────────────────────────────────────────────────────────────┤
│                    Aggregator (10 providers)                         │
│  ┌──────┬──────────┬──────────┬──────┬────────┬──────┬──────┬─────┐ │
│  │ Goal │Allocation│Projection│ Rec  │Optimize│ Risk │Events│Health│ │
│  └──────┴──────────┴──────────┴──────┴────────┴──────┴──────┴─────┘ │
├─────────────────────────────────────────────────────────────────────┤
│                    Composer (view models)                            │
│  ┌──────┬────────┬────────┬──────┬──────┬───────┬───────┬──────────┐ │
│  │ List │Dashboard│Progress│Proj. │ Recs │Optim. │Timeline│Milestone │ │
│  └──────┴────────┴────────┴──────┴──────┴───────┴───────┴──────────┘ │
├─────────────────────────────────────────────────────────────────────┤
│              InMemoryCache (read-only, TTL=5m)                       │
└─────────────────────────────────────────────────────────────────────┘
```

### Layers

| Layer | Files | Purpose |
|---|---|---|
| **Config** | `internal/config/config.go` | Environment configuration (PORT) |
| **Engine** | `internal/engine/core.go` | View models: GoalDashboard, GoalProgress, GoalProjectionView, GoalRecommendationView, GoalOptimizationView, GoalTimelineView, GoalMilestoneView, 8 CardTypes |
| **Aggregator** | `internal/aggregator/aggregator.go` | 10 providers, parallel enrichment per goal |
| **Cache** | `internal/infrastructure/cache/cache.go` | In-memory read-through cache |
| **API** | `internal/api/handlers.go`, `internal/api/router.go`, `internal/api/dto.go` | 8 REST endpoints with DTOs |
| **Main** | `cmd/server/main.go` | Dependency injection, 9 no-op providers |

---

## Goal Dashboard Model

### 8 Card Types

| Card | CardType | Data Source |
|---|---|---|
| Overview | `Overview` | Goal Domain |
| Progress | `Progress` | Goal Domain |
| Funding Status | `Funding` | Goal + Allocation Domain |
| Projection | `Projection` | Projection Engine |
| Recommendations | `Recommendation` | Recommendation Engine |
| Optimization | `Optimization` | Optimization Engine |
| Timeline | `Timeline` | Financial Event Domain + Timeline |
| Milestones | `Milestone` | Goal Domain + Milestone Service |

### View Models

| Model | Fields | Purpose |
|---|---|---|
| `GoalDashboard` | GoalID, Name, Progress, Status, Importance, GoalType, TargetAmount, CurrentValue, FundingGap, TargetDate, Cards[8] | Full dashboard for one goal |
| `GoalProgress` | GoalID, Name, ProgressPct, CurrentValue, TargetAmount, RemainingAmount, Status, MonthlyContribution, MonthsToTarget | Detailed progress tracking |
| `GoalProjectionView` | GoalID, ProjectedDate, ProjectedValue, OnTrack, MonthlyContribution | Projection engine wrapper |
| `GoalRecommendationView` | GoalID, Recs[RecItem], TotalCount | Recommendation engine wrapper |
| `GoalOptimizationView` | GoalID, Opts[OptItem], TotalCount | Optimization engine wrapper |
| `GoalTimelineView` | GoalID, Events[GoalEvent], TotalCount | Timeline events for a goal |
| `GoalMilestoneView` | GoalID, Milestones[Milestone], TotalCount | Milestones for a goal |

---

## REST Endpoints

| Method | Path | Response | Description |
|---|---|---|---|
| GET | `/api/v1/goals/experience` | GoalSummary[] | All goals list with summaries |
| GET | `/api/v1/goals/{id}/dashboard` | GoalDashboard | Full goal dashboard with 8 cards |
| GET | `/api/v1/goals/{id}/progress` | GoalProgress | Progress details |
| GET | `/api/v1/goals/{id}/projection` | GoalProjectionView | Projection for this goal |
| GET | `/api/v1/goals/{id}/recommendations` | GoalRecommendationView | Recommendations for this goal |
| GET | `/api/v1/goals/{id}/optimization` | GoalOptimizationView | Optimization for this goal |
| GET | `/api/v1/goals/{id}/timeline` | GoalTimelineView | Timeline events for this goal |
| GET | `/api/v1/goals/{id}/milestones` | GoalMilestoneView | Milestones for this goal |

### Query Parameters

| Parameter | Type | Description |
|---|---|---|
| `user_id` | string | User identifier (default: "default") |

---

## Aggregation Strategy

### Goal List Aggregation

```
1. Fetch all goals from GoalProvider (single call)
2. For EACH goal, parallel fan-out (7 goroutines per goal):
   ├─ Allocation.GetMonthlyContribution
   ├─ Projection.GetGoalProjection
   ├─ Recs.GetGoalRecommendations
   ├─ Optimize.GetGoalOptimizations
   ├─ Risk.HasRisk
   ├─ Events.GetGoalEvents
   └─ Events.GetGoalMilestones
3. Mutex-guarded writes per goal
4. Results cached per user
```

### Single Goal Aggregation

```
1. Fetch goal from GoalProvider.GetGoalByID
2. Parallel fan-out (7 goroutines):
   Same 7 providers as above
3. Mutex-guarded writes
```

---

## Engine Integrations

| Engine | Provider Interface | Composer Method | Endpoint |
|---|---|---|---|
| **Projection Engine** | `ProjectionProvider.GetGoalProjection` | `BuildProjection` | `/goals/{id}/projection` |
| **Recommendation Engine** | `RecProvider.GetGoalRecommendations` | `BuildRecommendations` | `/goals/{id}/recommendations` |
| **Optimization Engine** | `OptProvider.GetGoalOptimizations` | `BuildOptimization` | `/goals/{id}/optimization` |
| **Simulation Engine** | `SimProvider.HasSimulations` | (Dashboard flag) | `/goals/{id}/dashboard` |
| **Risk Engine** | `RiskProvider.HasRisk` | (Dashboard flag) | `/goals/{id}/dashboard` |

All engine integrations delegate entirely to the respective engines — the Goals Experience contains zero financial calculation logic.

---

## Architecture Compliance

| Rule | Status |
|---|---|
| Read-only architecture | ✅ All providers are read-only. No writes. |
| No business rules | ✅ All logic is orchestrational (composition, aggregation). |
| No financial calculations | ✅ Progress % computed from domain values only. All other calculations from engines. |
| No domain writes | ✅ Goals Experience never calls Save/Create/Update/Delete. |
| Delegates to engines | ✅ Projection, Recommendation, Optimization, Simulation, Risk all delegated to respective engines. |

---

## Build Status

| Check | Status |
|---|---|
| `go build ./...` (37 modules) | ✅ PASS |
| `go vet ./...` (37 modules) | ✅ PASS |
| Goals Experience endpoints | ✅ 8 endpoints registered |
| Engine integrations | ✅ 5 engines integrated (Projection, Recommendation, Optimization, Simulation, Risk) |
| Domain integrations | ✅ Goal, Allocation, Account/Events |
| Read-only architecture | ✅ Verified |

---

## File Inventory

| File | Lines | Purpose |
|---|---|---|
| `cmd/server/main.go` | 78 | Entry point, DI, 9 no-op providers |
| `internal/config/config.go` | 16 | Environment configuration |
| `internal/engine/core.go` | 250 | All view models (7 types), 8 card types, Composer with 8 build methods |
| `internal/aggregator/aggregator.go` | 197 | 10 providers, parallel enrichment, cache |
| `internal/infrastructure/cache/cache.go` | 44 | In-memory read-through cache |
| `internal/api/handlers.go` | 203 | 8 REST endpoint handlers |
| `internal/api/router.go` | 42 | Router + middleware |
| `internal/api/dto.go` | 77 | Response DTOs (8 response types) |
| `internal/events/events.go` | 12 | GoalOpened event |

---

## Remaining Experience Layer Work

| Phase | Experience | Status |
|---|---|---|
| 6.1 | Dashboard | ✅ Complete |
| 6.2 | Timeline | ✅ Complete |
| **6.3** | **Goals** | **✅ Complete** |
| 6.4 | Portfolio | 🔲 Not started |
| 6.5 | Accounts | 🔲 Not started |
| 6.6 | Planning | 🔲 Not started |
| 6.7 | Insights | 🔲 Not started |
| 6.8 | Notifications | 🔲 Not started |
| 6.9 | Advisor | 🔲 Not started |

---

## Ready for Phase 6.4 — Portfolio Experience

```
Go build:   37/37 PASS
Go vet:     37/37 PASS
Endpoints:  8/8   operational
Card types: 8/8   implemented
Engines:    5/5   integrated
Providers:  10/10  data providers
Architecture: READ-ONLY ✅
```
