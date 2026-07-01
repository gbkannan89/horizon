# Horizon Backend API Integration Guide

**Target Audience:** Flutter Development Team
**Version:** 1.0.0
**Status:** Production Ready
**Last Updated:** 2026-07-01

---

## Table of Contents

1. [Backend Architecture](#1-backend-architecture)
2. [Quick Start](#2-quick-start)
3. [Authentication & JWT Flow](#3-authentication--jwt-flow)
4. [API Conventions](#4-api-conventions)
5. [Endpoint Reference](#5-endpoint-reference)
   - [5.1 Auth](#51-auth)
   - [5.2 Dashboard](#52-dashboard)
   - [5.3 Accounts](#53-accounts)
   - [5.4 Portfolio](#54-portfolio)
   - [5.5 Goals](#55-goals)
   - [5.6 Planning](#56-planning)
   - [5.7 Insights](#57-insights)
   - [5.8 Timeline](#58-timeline)
   - [5.9 Advisor](#59-advisor)
   - [5.10 Notifications](#510-notifications)
   - [5.11 Transactions](#511-transactions)
   - [5.12 AI](#512-ai)
6. [Error Handling](#6-error-handling)
7. [Pagination](#7-pagination)
8. [Filtering & Search](#8-filtering--search)
9. [Environment Variables](#9-environment-variables)
10. [Docker & Deployment](#10-docker--deployment)
11. [Database Migrations & Seed Data](#11-database-migrations--seed-data)
12. [Recommended Flutter Architecture](#12-recommended-flutter-architecture)
13. [Known Limitations](#13-known-limitations)
14. [Appendix: Status Codes](#14-appendix-status-codes)

---

## 1. Backend Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   Flutter Mobile App                         │
│                  (Dio HTTP Client)                           │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTP (localhost:8080 / staging.horizon.app)
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                 Horizon Modular Monolith                     │
│                                                              │
│  cmd/server/main.go                                          │
│                                                              │
│  ┌─────────────┐  ┌──────────────┐  ┌────────────────────┐  │
│  │ Auth (JWT)  │  │  CORS        │  │  Recovery           │  │
│  │ Middleware   │  │  Middleware  │  │  Middleware         │  │
│  └─────────────┘  └──────────────┘  └────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────────┐│
│  │                  http.ServeMux                            ││
│  │                                                          ││
│  │  /api/v1/auth/*        → Auth Module                     ││
│  │  /api/v1/dashboard/*   → Dashboard Experience            ││
│  │  /api/v1/accounts/*    → Accounts Experience             ││
│  │  /api/v1/portfolio/*   → Portfolio Experience            ││
│  │  /api/v1/goals/*       → Goals Experience                ││
│  │  /api/v1/planning/*    → Planning Experience             ││
│  │  /api/v1/insights/*    → Insights Experience             ││
│  │  /api/v1/timeline/*    → Timeline Experience             ││
│  │  /api/v1/advisor/*     → Advisor Experience              ││
│  │  /api/v1/notifications/* → Notifications Experience      ││
│  │  /api/v1/transactions/* → Financial Event Domain         ││
│  │  /api/v1/ai/*           → AI Service                     ││
│  └──────────────────────────────────────────────────────────┘│
│                                                              │
│  ┌──────────────────────────────────────────────────────────┐│
│  │              PostgreSQL Connection Pool                   ││
│  └──────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

**Key architectural decisions:**
- **Modular Monolith** — single Go binary, single HTTP server, single PostgreSQL pool
- **No microservices** — all experiences run in the same process
- **No service discovery** — everything resolves locally
- **No gRPC** — REST/JSON only
- **JWT-based auth** — HS256 tokens, no session store needed

---

## 2. Quick Start

### Prerequisites

- Go 1.24+
- Docker Desktop
- Flutter SDK 3.5+

### Start the backend

```bash
cd horizon
export DATABASE_URL="postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable"
go run ./cmd/server
# Server starts on :8080
```

### Start with Docker (full stack)

```bash
cd horizon/deploy/docker
docker compose up
# Postgres :5432, Redis :6379, API :8080
```

### Run migrations

```bash
cd horizon/deploy/docker
docker compose exec postgres psql -U horizon -d horizon \
  -f /migrations/V001__extensions.sql
# ... repeat for V002 through V017
```

### Load seed data

```bash
cd horizon/deploy/docker
docker compose exec -T postgres psql -U horizon -d horizon \
  < ../../services/infra/migrations/S001__seed_demo_data.sql
```

### Verify

```bash
curl http://localhost:8080/health/live
# {"status":"live"}
```

### Login (MVP accepts any valid email/password)

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@horizon.app","password":"password123"}'
```

---

## 3. Authentication & JWT Flow

### 3.1 Login

```
POST /api/v1/auth/login

Request:  { "email": "demo@horizon.app", "password": "password123" }
Response: {
  "access_token": "eyJ...",
  "refresh_token": "eyJ...",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

**MVP note:** Currently accepts any valid `email` + `password` (min 8 chars). Future versions will verify against stored credentials.

### 3.2 Using the access token

Include the token in every request:

```
Authorization: Bearer eyJ...
```

### 3.3 Token refresh

```
POST /api/v1/auth/refresh

Request:  { "refresh_token": "eyJ..." }
Response: { "access_token": "eyJ...", "refresh_token": "eyJ...", ... }
```

### 3.4 Token lifecycle

| Token | TTL | Usage |
|---|---|---|
| `access_token` | 15 minutes | Authenticate API requests |
| `refresh_token` | 7 days | Obtain new access tokens |

### 3.5 Token format

HS256 JWT with payload:

```json
{
  "uid": "user-demo@horizon.app",
  "email": "demo@horizon.app",
  "roles": ["member"],
  "ttype": "access",
  "exp": 1700000000,
  "iat": 1699999100
}
```

### 3.6 Logout

```
POST /api/v1/auth/logout
Authorization: Bearer eyJ...

Response: 204 No Content
```

Currently client-side only (discard tokens). Future: token blacklisting via Redis.

### 3.7 Get current user

```
GET /api/v1/auth/me
Authorization: Bearer eyJ...

Response: { "user_id": "user-demo@horizon.app", "email": "demo@horizon.app", "roles": ["member"] }
```

### 3.8 Flutter token refresh flow

```dart
// Recommended implementation using Dio interceptor:
class AuthInterceptor extends Interceptor {
  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401) {
      // Attempt token refresh
      try {
        final response = await dio.post('/auth/refresh', data: {
          'refresh_token': await secureStorage.read('refresh_token'),
        });
        await secureStorage.write('access_token', response.data['access_token']);
        await secureStorage.write('refresh_token', response.data['refresh_token']);
        // Retry the original request with new token
        final retryResponse = await dio.fetch(err.requestOptions);
        handler.resolve(retryResponse);
        return;
      } catch (_) {
        // Refresh failed — redirect to login
      }
    }
    handler.next(err);
  }
}
```

---

## 4. API Conventions

### 4.1 Base URL

- Development: `http://localhost:8080`
- Staging: `https://staging.horizon.app`
- Production: `https://api.horizon.app`

### 4.2 Content Type

All requests and responses use `application/json`.

### 4.3 Success response format

```json
{
  "success": true,
  "data": { ... },
  "metadata": {
    "timestamp": "2026-07-01T14:30:00Z"
  }
}
```

### 4.4 Error response format

```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "Account not found"
  },
  "metadata": {
    "timestamp": "2026-07-01T14:30:00Z"
  }
}
```

### 4.5 Paginated response format

```json
{
  "success": true,
  "data": {
    "items": [ ... ],
    "cursor": "item-42",
    "has_more": true,
    "total": 100
  },
  "metadata": {
    "timestamp": "2026-07-01T14:30:00Z"
  }
}
```

### 4.6 Common query parameters

| Parameter | Type | Description | Used by |
|---|---|---|---|
| `user_id` | string | User identifier | All endpoints |
| `limit` | int | Page size (default 25, max 100) | All list endpoints |
| `cursor` | string | Pagination cursor | All list endpoints |

### 4.7 CORS

The API allows all origins with:
- `Access-Control-Allow-Origin: *`
- Methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
- Headers: Authorization, Content-Type, Idempotency-Key, X-Correlation-ID

---

## 5. Endpoint Reference

---

### 5.1 Auth

#### POST /api/v1/auth/login

**Auth:** None

**Request:**
```json
{
  "email": "demo@horizon.app",
  "password": "password123"
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 900,
  "token_type": "Bearer"
}
```

**Errors:** `400` VALIDATION_ERROR, `500` INTERNAL_ERROR

---

#### POST /api/v1/auth/refresh

**Auth:** None

**Request:**
```json
{
  "refresh_token": "eyJ..."
}
```

**Response 200:** Same as login.

**Errors:** `400` VALIDATION_ERROR, `401` AUTHENTICATION_ERROR

---

#### POST /api/v1/auth/logout

**Auth:** Bearer token

**Response:** `204 No Content`

---

#### GET /api/v1/auth/me

**Auth:** Bearer token

**Response 200:**
```json
{
  "user_id": "user-demo@horizon.app",
  "email": "demo@horizon.app",
  "roles": ["member"]
}
```

**Errors:** `401` AUTHENTICATION_ERROR

---

### 5.2 Dashboard

Base path: `/api/v1/dashboard`

#### GET /api/v1/dashboard

**Auth:** None

**Query:** `?user_id=demo@horizon.app&mode=Overview`

**Response 200:**
```json
{
  "success": true,
  "data": {
    "dashboard_id": "dash-user-demo@horizon.app",
    "mode": "Overview",
    "state": "Healthy",
    "tier1": [
      { "widget_id": "health", "widget_type": "health_score", "title": "Health: 72 (Good)", "tier": 1, "priority": 1, "confidence": "High", "visible": true },
      { "widget_id": "goals", "widget_type": "goal_progress", "title": "3/5 Goals On Track", "tier": 1, "priority": 2, "visible": true }
    ],
    "tier2": [
      { "widget_id": "nw", "widget_type": "net_worth", "title": "₹46,00,000", "tier": 2, "priority": 1, "visible": true }
    ],
    "tier3": [],
    "last_refresh_time": "2026-07-01T14:30:00Z"
  }
}
```

#### GET /api/v1/dashboard/summary

**Query:** `?user_id=demo@horizon.app`

**Response:**
```json
{
  "success": true,
  "data": {
    "net_worth": 4600000,
    "health_score": 72,
    "health_grade": "Good",
    "risk_score": 35,
    "risk_level": "Low",
    "goals_on_track": 3,
    "total_goals": 5,
    "cash_balance": 500000,
    "portfolio_value": 3500000,
    "total_debt": 5045000
  }
}
```

#### GET /api/v1/dashboard/health

Returns the Health Score widget.

#### GET /api/v1/dashboard/recommendations

Returns active recommendations.

#### GET /api/v1/dashboard/widgets

Returns all widgets flattened into a single list.

#### GET /api/v1/dashboard/milestones

Returns upcoming milestones.

---

### 5.3 Accounts

Base path: `/api/v1/accounts`

#### GET /api/v1/accounts/experience

Full accounts dashboard — balance card, cash flow, health, filter chips, account list.

#### GET /api/v1/accounts/summary

**Query:** `?user_id=demo@horizon.app&type=Savings&status=Active&limit=25&cursor=`

Returns paginated account list with filter support.

**Response:**
```json
{
  "success": true,
  "data": {
    "accounts": [
      {
        "account_id": "c1d2e3f4-...",
        "account_name": "Primary Savings Account",
        "account_type": "Savings",
        "classification": "Personal",
        "status": "Active",
        "currency": "INR",
        "current_balance": 500000,
        "available_balance": 480000,
        "institution_name": "HDFC Bank",
        "has_import": false
      }
    ],
    "total": 5,
    "has_more": false,
    "summary": {
      "total_balance": 1500000,
      "count_by_type": { "Savings": 2, "Checking": 1, "Investment": 1, "Credit": 1 },
      "count_by_status": { "Active": 5 }
    }
  }
}
```

#### GET /api/v1/accounts/balances

**Query:** `?user_id=demo@horizon.app`

Returns balance summary with Total/Available/Spendable breakdown.

#### GET /api/v1/accounts/cashflow

**Query:** `?user_id=demo@horizon.app`

Returns cash flow summary (inflow, outflow, net, projected).

#### GET /api/v1/accounts/health

Returns account health breakdown by status.

#### GET /api/v1/accounts/recommendations

Returns account-related recommendations.

#### GET /api/v1/accounts/projections

Returns cash flow projection.

#### GET /api/v1/accounts/timeline

Returns recent account activity count.

#### GET /api/v1/accounts/{id}

Returns full account detail with balances (7 types) and recent transactions.

**Response:**
```json
{
  "success": true,
  "data": {
    "account_id": "c1d2e3f4-...",
    "account_name": "Primary Savings Account",
    "account_type": "Savings",
    "status": "Active",
    "currency": "INR",
    "balances": [
      { "type": "current", "value": 500000, "label": "Current Balance", "description": "Net sum of all posted transactions" },
      { "type": "available", "value": 480000, "label": "Available Balance", "description": "Current + uncleared - reserved" }
    ],
    "transactions": [
      {
        "event_id": "e1f2a3b4-...",
        "event_type": "Salary",
        "amount": 150000,
        "currency": "INR",
        "description": "Monthly salary credit",
        "category": "",
        "event_date": "2026-06-25T14:30:00Z",
        "running_balance": 0
      }
    ],
    "institution_name": "HDFC Bank",
    "opened_date": "2023-01-15",
    "account_health": "Healthy"
  }
}
```

---

### 5.4 Portfolio

Base path: `/api/v1/portfolio`

#### GET /api/v1/portfolio/experience

Full portfolio dashboard.

#### GET /api/v1/portfolio/allocation

**Query:** `?user_id=demo@horizon.app`

**Response:**
```json
{
  "success": true,
  "data": {
    "allocations": [
      { "label": "CashAndCashEquivalent", "value": 500000, "percent": 10.6 },
      { "label": "Equity", "value": 960000, "percent": 20.4 },
      { "label": "MutualFund", "value": 720000, "percent": 15.3 },
      { "label": "RetirementFund", "value": 1280000, "percent": 27.2 },
      { "label": "FixedDeposit", "value": 525000, "percent": 11.2 },
      { "label": "Gold", "value": 235000, "percent": 5.0 },
      { "label": "GovernmentSecurity", "value": 380000, "percent": 8.1 }
    ],
    "total_value": 4600000,
    "diversification_score": 5.2
  }
}
```

#### GET /api/v1/portfolio/performance

#### GET /api/v1/portfolio/risk

#### GET /api/v1/portfolio/projection

#### GET /api/v1/portfolio/recommendations

#### GET /api/v1/portfolio/optimization

#### GET /api/v1/portfolio/simulations

#### GET /api/v1/portfolio/timeline

---

### 5.5 Goals

Base path: `/api/v1/goals`

#### GET /api/v1/goals/experience

**Query:** `?user_id=demo@horizon.app`

**Response:**
```json
{
  "success": true,
  "data": {
    "goals": [
      { "goal_id": "d1e2f3a4-...", "name": "Retirement Corpus", "progress": 12, "status": "Active", "importance": "Mandatory", "priority": 1, "has_recommendation": true },
      { "goal_id": "d1e2f3a4-...", "name": "Emergency Fund", "progress": 40, "status": "Active", "importance": "Essential", "priority": 2, "has_recommendation": false }
    ],
    "total": 4
  }
}
```

#### GET /api/v1/goals/{id}/dashboard

Full goal dashboard with header card + 8 dashboard cards.

**Response:**
```json
{
  "success": true,
  "data": {
    "goal_id": "d1e2f3a4-...",
    "goal_name": "Retirement Corpus",
    "progress": 12,
    "status": "Active",
    "importance": "Mandatory",
    "target_amount": 50000000,
    "current_value": 6000000,
    "funding_gap": 44000000,
    "target_date": "2050-01-01",
    "cards": [
      { "card_id": "ov", "card_type": "Overview", "title": "Retirement Corpus", "summary": "12% funded — Active" },
      { "card_id": "prog", "card_type": "Progress", "title": "Progress", "summary": "₹60,00,000 / ₹5,00,00,000" },
      { "card_id": "fund", "card_type": "Funding", "title": "Funding Status", "summary": "Gap: ₹4,40,00,000" }
    ]
  }
}
```

#### GET /api/v1/goals/{id}/progress

Detailed progress with remaining amount and months to target.

#### GET /api/v1/goals/{id}/projection

Returns projected completion date and on-track status.

#### GET /api/v1/goals/{id}/recommendations

Returns goal-specific recommendations.

#### GET /api/v1/goals/{id}/optimization

Returns goal optimization strategies.

#### GET /api/v1/goals/{id}/timeline

Returns timeline events for this goal.

#### GET /api/v1/goals/{id}/milestones

Returns milestones for this goal.

---

### 5.6 Planning

Base path: `/api/v1/planning`

#### GET /api/v1/planning

#### GET /api/v1/planning/dashboard

**Query:** `?user_id=demo@horizon.app&mode=Overview`

Returns planning workspace with cards and current plan summary.

#### GET /api/v1/planning/projections

Returns net worth, income, expenses, portfolio projections.

#### GET /api/v1/planning/scenarios

Returns saved planning scenarios (currently empty).

#### GET /api/v1/planning/compare

**Query:** `?user_id=demo@horizon.app`

Compares current plan against first scenario with deltas.

#### GET /api/v1/planning/recommendations

#### GET /api/v1/planning/optimizations

#### GET /api/v1/planning/risk

#### GET /api/v1/planning/health

#### GET /api/v1/planning/timeline

---

### 5.7 Insights

Base path: `/api/v1/insights`

#### GET /api/v1/insights

**Query:** `?user_id=demo@horizon.app&limit=10`

Returns paginated insights feed.

#### GET /api/v1/insights/summary

Returns dashboard with counts by priority and category.

#### GET /api/v1/insights/opportunities

Returns opportunity-category insights.

#### GET /api/v1/insights/warnings

Returns warning + risk-category insights.

#### GET /api/v1/insights/achievements

Returns achievement + milestone insights.

#### GET /api/v1/insights/forecast

Returns forecast + cash flow insights.

#### GET /api/v1/insights/trends

Returns savings, net worth, and portfolio trend insights.

#### GET /api/v1/insights/timeline

Recent insight feed.

#### GET /api/v1/insights/search

**Query:** `?user_id=demo@horizon.app&q=retirement`

Full-text search across insight titles and summaries.

#### GET /api/v1/insights/{id}

Returns a single insight by ID.

---

### 5.8 Timeline

Base path: `/api/v1/timeline`

#### GET /api/v1/timeline

**Query:** `?user_id=demo@horizon.app&view=ThisMonth&limit=25&cursor=`

Returns paginated, cursor-based timeline feed.

**Response:**
```json
{
  "success": true,
  "data": {
    "view": "ThisMonth",
    "feed": {
      "items": [
        {
          "timeline_id": "fin-e1f2a3b4-...",
          "timestamp": "2026-06-25T14:30:00Z",
          "event_type": "Salary",
          "category": "FinancialEvent",
          "title": "Salary: Monthly salary credit",
          "summary": "+₹1,50,000",
          "severity": "info",
          "amount": 150000,
          "icon": "currency_rupee",
          "color": "green"
        },
        {
          "timeline_id": "goal-d1e2f3a4-...-Active",
          "timestamp": "2026-06-01T00:00:00Z",
          "event_type": "GoalActive",
          "category": "GoalEvent",
          "title": "Active: Retirement Corpus",
          "summary": "Goal Active",
          "severity": "info",
          "icon": "track_changes",
          "color": "blue"
        }
      ],
      "cursor": "fin-e1f2a3b4-...",
      "has_more": false,
      "total": 2
    },
    "filters": {},
    "total": 2
  }
}
```

#### GET /api/v1/timeline/recent

Recent items (default 10).

#### GET /api/v1/timeline/filter

**Query:** `?user_id=demo@horizon.app&categories=FinancialEvent,GoalEvent&severity=info&start_date=2026-01-01&end_date=2026-12-31`

Full filtering with date range, categories, severity.

#### GET /api/v1/timeline/search

**Query:** `?user_id=demo@horizon.app&q=salary`

Full-text search across title, summary, description, event type.

#### GET /api/v1/timeline/{id}

Single timeline item detail with all fields.

---

### 5.9 Advisor

Base path: `/api/v1/advisor`

#### GET /api/v1/advisor

Returns full advisor workspace with cards.

#### GET /api/v1/advisor/dashboard

Same as advisor.

#### GET /api/v1/advisor/summary

Returns financial summary (net worth, cash balance, monthly income/expenses).

#### GET /api/v1/advisor/context

Returns the structured AdvisorContext (14 blocks) for AI consumption.

```json
{
  "success": true,
  "data": {
    "user_summary": { "user_id": "demo@horizon.app" },
    "financial_summary": { "net_worth": 4600000, "cash_balance": 500000, ... },
    "goal_summary": { "total_goals": 5, "on_track": 3, "at_risk": 1, "funding_gap": 500000 },
    "health": { "score": 72, "grade": "Good", "change": 3 },
    "risk": { "score": 35, "level": "Low" },
    "recommendations": { "has_recommendations": true, "count": 3 },
    ...
  }
}
```

#### GET /api/v1/advisor/health

#### GET /api/v1/advisor/risk

#### GET /api/v1/advisor/recommendations

#### GET /api/v1/advisor/insights

#### GET /api/v1/advisor/timeline

#### GET /api/v1/advisor/notifications

---

### 5.10 Notifications

Base path: `/api/v1/notifications`

#### GET /api/v1/notifications

**Query:** `?user_id=demo@horizon.app&limit=25&cursor=`

Returns paginated notification feed.

**Response:**
```json
{
  "success": true,
  "data": {
    "unread_count": 3,
    "by_priority": { "P1": 1, "P2": 1, "P3": 1 },
    "by_category": { "Health": 1, "Goal": 1, "Recommendation": 1 },
    "items": [
      {
        "notif_id": "notif-health-critical",
        "category": "Health",
        "priority": "P1",
        "state": "Unread",
        "title": "Critical Health Score",
        "summary": "Your health score is 35 — immediate attention",
        "timestamp": "2026-07-01T14:30:00Z",
        "action_url": "/dashboard"
      }
    ],
    "has_more": false
  }
}
```

#### GET /api/v1/notifications/unread

#### GET /api/v1/notifications/history

#### GET /api/v1/notifications/preferences

#### PUT /api/v1/notifications/preferences

**Request:** Array of preference objects with category, enabled, min_priority.

#### GET /api/v1/notifications/search

**Query:** `?q=health`

#### GET /api/v1/notifications/{id}

#### POST /api/v1/notifications/{id}/read

#### POST /api/v1/notifications/{id}/archive

#### POST /api/v1/notifications/{id}/snooze

**Request:** `{ "until": "2026-07-02T14:30:00Z" }`

---

### 5.11 Transactions

Base path: `/api/v1/transactions`

#### GET /api/v1/transactions

**Query:** `?user_id=demo@horizon.app&limit=25&cursor=`

**Response:**
```json
{
  "success": true,
  "data": {
    "transactions": [
      {
        "event_id": "e1f2a3b4-...",
        "event_type": "Salary",
        "amount": 150000,
        "currency": "INR",
        "description": "Monthly salary credit",
        "source": "Employer",
        "destination": "Salary Account",
        "effective_date": "2026-06-25T14:30:00Z",
        "state": "POSTED"
      }
    ],
    "cursor": "e1f2a3b4-...",
    "has_more": false,
    "total": 15
  }
}
```

#### GET /api/v1/transactions/{id}

Full transaction detail including origin, confidence, correlation_id.

#### GET /api/v1/transactions/search

**Query:** `?user_id=demo@horizon.app&q=salary&limit=25`

Full-text search across description, event_type, source, destination.

#### GET /api/v1/transactions/summary

**Query:** `?user_id=demo@horizon.app`

Monthly income/expense summary.

```json
{
  "success": true,
  "data": {
    "period_income": 150000,
    "period_expenses": 85000,
    "net_flow": 65000,
    "income_count": 1,
    "expense_count": 5,
    "total_count": 6
  }
}
```

---

### 5.12 AI

Base path: `/api/v1/ai`

All AI endpoints work with the **mock provider** by default — no Ollama required.

#### GET /api/v1/ai/health

```json
{
  "success": true,
  "data": { "status": "operational", "provider": "mock", "message": "Mock AI provider — ready." }
}
```

#### POST /api/v1/ai/chat

```json
{
  "session_id": "session-20260701-143000",
  "message": "How is my health score?",
  "context": { "health_score": 72 }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "reply": "[Mock] Your financial health is calculated from 13 dimensions...",
    "confidence": "Medium",
    "provider": "mock",
    "session_id": "session-20260701-143000"
  }
}
```

#### POST /api/v1/ai/explain

```json
{
  "prompt_type": "explain_health",
  "data": { "score": 72, "grade": "Good", "change": 3 },
  "version": 1
}
```

#### POST /api/v1/ai/summarize

```json
{ "topic": "financial_summary", "data": { "net_worth": 4600000, "income": 150000, "expenses": 100000 } }
```

#### GET /api/v1/ai/prompts

Lists all registered prompt templates with their versions.

#### GET /api/v1/ai/context

Returns the full AdvisorContext (same as `/api/v1/advisor/context`).

#### GET /api/v1/ai/providers

Lists all registered AI providers (mock, ollama, openai, azure, anthropic).

#### GET /api/v1/ai/providers/active

Returns the currently active provider.

#### GET /api/v1/ai/providers/capabilities

Returns the active provider's capabilities (chat, explain, summarize, streaming, etc.).

#### POST /api/v1/ai/provider/switch

```json
{ "provider": "mock" }
```

#### GET /api/v1/ai/runtime/health

#### GET /api/v1/ai/runtime/config

#### GET /api/v1/ai/runtime/status

---

## 6. Error Handling

### 6.1 Error response structure

Every error follows the same format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable description"
  },
  "metadata": {
    "timestamp": "2026-07-01T14:30:00Z"
  }
}
```

### 6.2 Error codes

| HTTP Status | Code | Meaning |
|---|---|---|
| 400 | `VALIDATION_ERROR` | Invalid request body or parameters |
| 401 | `AUTHENTICATION_ERROR` | Missing or invalid JWT |
| 403 | `AUTHORIZATION_ERROR` | Insufficient permissions |
| 404 | `NOT_FOUND` | Resource not found |
| 409 | `CONFLICT` | State conflict (duplicate, invalid transition) |
| 429 | `RATE_LIMITED` | Too many requests |
| 500 | `INTERNAL_ERROR` | Unexpected server error |
| 503 | `SERVICE_UNAVAILABLE` | Downstream dependency unavailable |

### 6.3 Flutter error handling pattern

```dart
class ApiException implements Exception {
  final int statusCode;
  final String code;
  final String message;

  bool get isAuthError => statusCode == 401;
  bool get isNotFound => statusCode == 404;
  bool get isValidation => statusCode == 400;
  bool get isServerError => statusCode >= 500;
}
```

---

## 7. Pagination

All list endpoints use **cursor-based pagination**.

### Request

| Parameter | Type | Default | Description |
|---|---|---|---|
| `limit` | int | 25 | Items per page (max 100) |
| `cursor` | string | — | Opaque cursor from previous response |

### Response

| Field | Type | Description |
|---|---|---|
| `cursor` | string | Pass as `cursor` in next request. Empty = no more pages. |
| `has_more` | bool | `true` if more items exist |
| `total` | int | Items in this page |

### Flutter pattern

```dart
class PaginatedState<T> {
  final List<T> items;
  final String? cursor;
  final bool isLoading;

  Future<void> loadNext() async {
    if (isLoading || cursor == null) return;
    final response = await api.get('/path', queryParameters: {
      'limit': 25, 'cursor': cursor,
    });
    items.addAll(response.data['items']);
    cursor = response.data['cursor'];
  }
}
```

---

## 8. Filtering & Search

### Filters (Timeline and Account list)

| Parameter | Type | Example |
|---|---|---|
| `categories` | string (comma-sep) | `FinancialEvent,GoalEvent` |
| `severity` | string | `info`, `warning`, `critical`, `success` |
| `start_date` | string (RFC3339) | `2026-01-01T00:00:00Z` |
| `end_date` | string (RFC3339) | `2026-12-31T23:59:59Z` |
| `type` | string | `Savings`, `Checking` (account type) |
| `status` | string | `Active`, `Inactive` |
| `q` | string | Search query |

### Search endpoints

| Endpoint | Fields searched |
|---|---|
| `/timeline/search` | title, summary, description, event_type |
| `/insights/search` | title, summary, description |
| `/transactions/search` | description, event_type, source, destination |
| `/notifications/search` | title, summary, description |

---

## 9. Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP server port |
| `DATABASE_URL` | `postgres://horizon:horizon@localhost:5432/horizon?sslmode=disable` | PostgreSQL connection string |
| `APP_ENV` | `development` | Application environment |
| `LOG_LEVEL` | `info` | Logging level |
| `AI_PROVIDER` | `mock` | AI provider (mock, ollama, openai, azure, anthropic) |
| `AI_ENABLED` | `true` | Enable AI features |
| `OLLAMA_URL` | `http://localhost:11434` | Ollama server URL |
| `OLLAMA_MODEL` | `llama3` | Default Ollama model |
| `AI_TIMEOUT_SECONDS` | `30` | AI request timeout |
| `AI_RETRIES` | `2` | AI request retry count |

---

## 10. Docker & Deployment

### Docker build

```bash
cd horizon
docker build -t horizon-api .
docker run -p 8080:8080 \
  -e DATABASE_URL="postgres://horizon:horizon@host.docker.internal:5432/horizon?sslmode=disable" \
  horizon-api
```

### Docker Compose (full stack)

```bash
cd horizon/deploy/docker
docker compose up --build
```

This starts:
- `postgres` — PostgreSQL 17 on `:5432`
- `redis` — Redis 7 on `:6379` (optional, for future rate limiting)
- `horizon-api` — The monolithic backend on `:8080`

**Note:** The monitoring stack (prometheus, grafana, tempo, loki, otel-collector) has known file-bind-mount issues on Windows Docker Desktop. If you encounter mount errors, start only essential services:

```bash
docker compose up postgres redis horizon-api
```

---

## 11. Database Migrations & Seed Data

### Migration files (apply in order)

Located at `services/infra/migrations/`:

```
V001__extensions.sql
V002__foundation.sql
V003__user.sql
V004__institution.sql
V005__financial_event.sql
V006__goal.sql
V007__account.sql
V008__allocation.sql
V009__asset.sql
V010__liability.sql
V011__portfolio.sql
V012__projection_engine.sql
V013__risk_engine.sql
V014__health_score_engine.sql
V015__recommendation_engine.sql
V016__optimization_engine.sql
V017__simulation_engine.sql
```

### Apply migrations

```bash
cd horizon/deploy/docker
for f in ../../services/infra/migrations/V*.sql; do
  docker compose exec -T postgres psql -U horizon -d horizon < "$f"
done
```

### Seed data

```bash
docker compose exec -T postgres psql -U horizon -d horizon \
  < ../../services/infra/migrations/S001__seed_demo_data.sql
```

### Seed data includes

| Table | Records | Demo user |
|---|---|---|
| users | 1 | Priya Sharma |
| accounts | 5 | Savings, Checking, FD, Investment, Credit |
| goals | 4 | Retirement, Emergency Fund, House, Education |
| assets | 7 | Equity, Mutual Funds, Gold, PPF, EPF, etc. |
| liabilities | 2 | Home Loan, Credit Card |
| financial_events | 15 | Salary, expenses, investments |
| health_scores | 1 | Score 72 (Good) |
| risk_assessments | 1 | Score 35 (Low) |
| recommendations | 3 | Actionable recommendations |

---

## 12. Recommended Flutter Architecture

### 12.1 Networking layer

```
lib/
└── core/
    └── network/
        ├── api_client.dart       # Dio instance with interceptors
        ├── auth_interceptor.dart  # JWT injection + 401 refresh
        ├── error_mapper.dart      # Error → AppError mapping
        └── api_response.dart     # Generic response wrapper
```

### 12.2 API client setup

```dart
final apiClient = Dio(BaseOptions(
  baseUrl: 'http://localhost:8080/api/v1',
  connectTimeout: Duration(seconds: 15),
  receiveTimeout: Duration(seconds: 30),
  headers: {'Content-Type': 'application/json'},
));

apiClient.interceptors.addAll([
  AuthInterceptor(),     // Add Bearer token, handle 401 refresh
  LogInterceptor(),      // Debug logging
]);
```

### 12.3 Response wrapper

```dart
class ApiResponse<T> {
  final bool success;
  final T? data;
  final String? errorCode;
  final String? errorMessage;
  final String? cursor;
  final bool hasMore;
}
```

### 12.4 Repository pattern

```dart
class DashboardRepository {
  final Dio client;
  
  Future<DashboardData> getDashboard({String? userId}) async {
    final response = await client.get('/dashboard', queryParameters: {
      if (userId != null) 'user_id': userId,
    });
    return DashboardData.fromJson(response.data['data']);
  }
}
```

### 12.5 State management (Riverpod)

```dart
final dashboardProvider = FutureProvider.family<DashboardData, String?>((ref, userId) async {
  final repo = ref.read(dashboardRepositoryProvider);
  return repo.getDashboard(userId: userId);
});
```

---

## 13. Known Limitations

| Limitation | Impact | Planned Fix |
|---|---|---|
| **MVP Login** — accepts any email/password without verification | No real auth | Wire User domain service to verify credentials |
| **No auth middleware** — JWT validation is optional, not enforced on experience endpoints | Endpoints publicly accessible | Add JWT middleware to all protected routes |
| **No CORS for credentials** — `Access-Control-Allow-Origin: *` doesn't support cookies | Cookies not usable | Use token-based auth (recommended) |
| **In-memory notification state** — read/archive/snooze resets on restart | State lost after restart | Persist notification states to database |
| **No pagination on some endpoints** — Insights, Notifications use simple limit | Large datasets unsupported | Add cursor-based pagination |
| **Goals N+1 queries** — each goal triggers 7 additional DB queries | ~71 queries for 10 goals | Batch query projections, events per goal |
| **No rate limiting** | No DoS protection | Add Redis-based rate limiting |
| **No structured logging** — uses stdlib log | Hard to integrate with log aggregators | Switch to zerolog/zap |

---

## 14. Appendix: Status Codes

| Code | Name | Usage |
|---|---|---|
| 200 | OK | Successful GET, PUT, PATCH |
| 201 | Created | Successful POST (resource created) |
| 204 | No Content | Successful DELETE, logout |
| 400 | Bad Request | Validation errors, malformed JSON |
| 401 | Unauthorized | Missing or invalid JWT |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Duplicate, invalid state transition |
| 422 | Unprocessable | Semantic validation failure |
| 429 | Too Many Requests | Rate limit exceeded |
| 500 | Internal Server Error | Unexpected server failure |
| 503 | Service Unavailable | Database down, dependency unavailable |
