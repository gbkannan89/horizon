# AI FOUNDATION REPORT

**Phase:** 7.1
**Status:** Complete
**Build:** 38/38 modules PASS (zero external dependencies)

---

## AI Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                     REST API (5 endpoints)                            │
├──────────────────────────────────────────────────────────────────────┤
│                     Handlers (thin layer)                             │
├──────────────────────────────────────────────────────────────────────┤
│                     AI Service (orchestrator)                         │
├──────────────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐ │
│  │ Provider │  │  Prompt  │  │ Session  │  │ Advisor Context      │ │
│  │ Interface│  │ Templates│  │ Manager  │  │ Builder (deterministic)│ │
│  ├──────────┤  ├──────────┤  ├──────────┤  ├──────────────────────┤ │
│  │ Mock     │  │ 10 types │  │ Conver-  │  │ 14 context blocks    │ │
│  │ Provider │  │ + ver-   │  │ sation   │  │ → flat map for       │ │
│  │ (offline)│  │ sioning  │  │ history  │  │ prompt injection     │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────────────┘ │
├──────────────────────────────────────────────────────────────────────┤
│              InMemoryCache (read-only, TTL=5m)                        │
└──────────────────────────────────────────────────────────────────────┘
```

---

## Provider Interface

```go
type AIProvider interface {
    Explain(ctx context.Context, req ExplainRequest) (*ExplainResponse, error)
    Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error)
    Health(ctx context.Context) (*HealthResponse, error)
}
```

| Method | Input | Output |
|---|---|---|
| `Explain` | PromptType, Data (map), Version | Explanation text, Confidence, Provider, Version |
| `Chat` | SessionID, Message, History, Context | Reply, Confidence, Provider, SessionID |
| `Health` | — | Status, Provider, Message |

**Provider config:** `AI_PROVIDER` env var. Currently only `"mock"` is implemented.

---

## Mock Provider

- **Zero external dependencies** — no network, no LLM, no Ollama
- **Deterministic responses** — same inputs always produce same outputs
- **10 prompt type handlers** — explain_health, explain_risk, explain_goal, explain_portfolio, explain_recommendation, explain_optimization, explain_projection, explain_simulation, financial_summary, advisor_summary
- **Chat handler** — keyword-matched replies with fallback for unknown queries
- **Context-aware** — uses data map to parameterize all responses

---

## Advisor Context (14 blocks)

| Block | Fields | Source Engine/Domain |
|---|---|---|
| `user_summary` | UserID | User Domain |
| `financial_summary` | NetWorth, CashBalance, MonthlyIncome, MonthlyExp | Dashboard Experience |
| `goal_summary` | TotalGoals, OnTrack, AtRisk, FundingGap | Goals Experience |
| `account_summary` | TotalAccounts, TotalBalance | Accounts Experience |
| `portfolio_summary` | PortfolioValue, TotalReturn, ReturnPct, RiskScore | Portfolio Experience |
| `health` | Score, Grade, Change | Health Score Engine |
| `risk` | Score, Level | Risk Engine |
| `projection` | NetWorthProjected, OnTrack, Confidence | Projection Engine |
| `recommendations` | HasRecs, Count | Recommendation Engine |
| `optimizations` | HasOpts, Count | Optimization Engine |
| `simulations` | HasSims, Count | Simulation Engine |
| `timeline_events` | RecentCount | Timeline Experience |
| `notifications` | UnreadCount | Notifications Experience |
| `insights` | TotalCount, CriticalCount | Insights Experience |

`ToFlatMap()` serializes all blocks into a single flat `map[string]interface{}` for prompt template injection.

---

## Prompt Templates

| Prompt Type | Version 1 | Version 2 |
|---|---|---|
| `explain_health` | ✅ Basic health explanation | ✅ Detailed with strengths/weaknesses |
| `explain_risk` | ✅ Risk level explanation | — |
| `explain_goal` | ✅ Goal status explanation | — |
| `explain_portfolio` | ✅ Portfolio summary | — |
| `explain_recommendation` | ✅ Recommendation priority | — |
| `explain_optimization` | ✅ Optimization explanation | — |
| `explain_projection` | ✅ Projection methodology | — |
| `explain_simulation` | ✅ Simulation explanation | — |
| `financial_summary` | ✅ Financial overview | — |
| `advisor_summary` | ✅ Comprehensive summary | — |

Template placeholders use `{{placeholder}}` syntax. `Render()` substitutes from the context data map.

---

## Conversation Model

```
Manager
├── GetOrCreate(sessionID, userID) → Conversation
├── AddMessage(sessionID, Message)
├── GetHistory(sessionID) → []Message
├── SetContext(sessionID, map)
└── Get(sessionID) → Conversation

Conversation
├── SessionID, UserID
├── Messages[] (role: user/assistant/system)
├── CreatedAt, UpdatedAt
└── Context (map for deterministic data attachment)
```

Session history is bounded by `MaxHistory` (default: 50 messages). Oldest messages are pruned.

---

## REST Endpoints

| Method | Path | Request Body | Response | Description |
|---|---|---|---|---|
| GET | `/api/v1/ai/health` | — | `{provider, status, message}` | Provider health check |
| POST | `/api/v1/ai/chat` | `{session_id, message, context}` | `{reply, confidence, provider}` | Conversational chat |
| POST | `/api/v1/ai/explain` | `{prompt_type, data, version}` | `{explanation, confidence, provider}` | Explain an engine output |
| GET | `/api/v1/ai/prompts` | — | `[{type, version, versions}]` | List prompt templates |
| GET | `/api/v1/ai/context` | — | Full AdvisorContext (14 blocks) | Get AI-ready context |

---

## Configuration

| Env Variable | Default | Purpose |
|---|---|---|
| `AI_PROVIDER` | `"mock"` | AI provider to use |
| `PORT` | `"8200"` | HTTP server port |

**Future provider support** (configurable but not implemented): `ollama`, `openai`, `azure`, `anthropic`

---

## Architecture Compliance

| Rule | Status |
|---|---|
| No real AI/LLM | ✅ Mock provider only — offline |
| No external APIs | ✅ Zero external dependencies |
| No Ollama | ✅ Not required |
| No financial calculations | ✅ All values from context (deterministic) |
| No business rules | ✅ Orchestration only |
| AI Boundary | ✅ AI only explains — never computes |
| Deterministic context | ✅ Same inputs → same context |
| Provider-swappable | ✅ AIProvider interface |

---

## Build Status

| Check | Status |
|---|---|
| `go build ./...` (38 modules) | ✅ PASS |
| `go vet ./...` (38 modules) | ✅ PASS |
| Mock provider | ✅ Operational |
| Advisor Context | ✅ 14 blocks, flat map serialization |
| Prompt rendering | ✅ 10 types, versioned |
| Conversation | ✅ Session management, history |
| REST endpoints | ✅ 5 endpoints |
| External dependencies | ✅ **Zero** |
| Architecture boundary | ✅ AI Boundary respected |

---

## File Inventory

| File | Lines | Purpose |
|---|---|---|
| `cmd/server/main.go` | 37 | Entry point, provider selection, DI |
| `internal/config/config.go` | 20 | AI_PROVIDER, PORT configuration |
| `internal/provider/provider.go` | 68 | AIProvider interface (Explain, Chat, Health) |
| `internal/provider/mock.go` | 179 | Mock AI provider — deterministic, offline |
| `internal/prompts/templates.go` | 140 | Versioned prompt templates, Render() |
| `internal/session/session.go` | 97 | Conversation/session manager |
| `internal/context/builder.go` | 132 | AdvisorContext (14 blocks), ToFlatMap() |
| `internal/api/handlers.go` | 186 | 5 REST endpoint handlers |
| `internal/api/router.go` | 42 | Router + middleware |
| `internal/api/dto.go` | 85 | Response DTOs |
| `internal/infrastructure/cache/cache.go` | 44 | In-memory read-through cache |

---

## Ready for Phase 7.2 — Ollama Provider

```
Go build:         38/38 PASS
Endpoints:        5/5   operational
Prompt Types:     10/10 registered
Context Blocks:   14/14 available
External Deps:    0     (zero)
AI Boundary:      CLEAN ✅
```
