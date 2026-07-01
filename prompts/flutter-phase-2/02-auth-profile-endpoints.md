# Phase 2 — Milestone 2: Auth & Profile Endpoints

**Context Reminder:** Architecture frozen. Specs frozen. Stack frozen. Monorepo. No Redis/K8s/Grafana/testing. Phase 1 completed all Flutter features. Phase 2 fills backend gaps.

## Objective

Implement missing backend HTTP endpoints: Change Password, User Profile CRUD, Preferences, Privacy.

## Backend Audit

Read these files:

1. `services/experiences/auth/register/register.go` — current login handler, modify to support change-password
2. `services/domains/user/internal/application/service.go` — has `UpdateProfile`, `UpdatePreferences`, `UpdatePrivacy`, `GrantConsent` methods (domain code exists but no HTTP routes)
3. `services/infra/api-gateway/internal/router/routes.go` — has basic user CRUD (create, list, get, activate, suspend, archive) but these are mock stubs
4. `cmd/server/main.go` — verify the user domain is NOT registered in the monolith

Confirm:
- `POST /api/v1/auth/change-password` — ❌ DOES NOT EXIST
- `GET /api/v1/users/me` — ❌ DOES NOT EXIST
- `PUT /api/v1/users/me/profile` — ❌ DOES NOT EXIST
- `PUT /api/v1/users/me/preferences` — ❌ DOES NOT EXIST
- `PUT /api/v1/users/me/privacy` — ❌ DOES NOT EXIST
- `POST /api/v1/users/me/consent` — ❌ DOES NOT EXIST

## Implementation

### Change Password

Create a simple handler in `services/experiences/auth/register/register.go`:
- `POST /api/v1/auth/change-password`
- Accepts `{current_password, new_password}`
- For MVP: accept any password (same stub pattern as login). Returns 200.
- Flutter UI already exists at `/settings/security` — just needs a working backend

### User Profile Routes

Create a new file `services/experiences/user/register/register.go` (or extend auth register):
- `GET /api/v1/users/me` — returns `{user_id, email, name, phone, country, base_currency, timezone, language, created_at}`. Query from `users` table by user_id from JWT.
- `PUT /api/v1/users/me/profile` — updates `{name, phone, date_of_birth, country, base_currency}`. Writes to `users` table.
- `PUT /api/v1/users/me/preferences` — reads/writes `user_preferences` table
- `PUT /api/v1/users/me/privacy` — reads/writes `user_privacy` table

### Wire into monolith

Add `userReg.RegisterRoutes(mux, pool)` in `cmd/server/main.go`

## Files to Create/Modify

| File | Change |
|------|--------|
| `services/experiences/auth/register/register.go` | Add POST /auth/change-password handler |
| `services/experiences/user/register/register.go` | NEW — user profile/preferences/privacy handlers |
| `services/experiences/user/internal/api/handlers.go` | NEW — HTTP handlers |
| `services/experiences/user/internal/api/dto.go` | NEW — request/response types |
| `cmd/server/main.go` | Add userReg import and RegisterRoutes call |

## Verification

- [ ] Backend builds: `go build ./cmd/server/...`
- [ ] `POST /auth/change-password` returns 200
- [ ] `GET /users/me` returns user profile from DB
- [ ] `PUT /users/me/profile` updates and returns profile
- [ ] `PUT /users/me/preferences` updates preferences
- [ ] `PUT /users/me/privacy` updates privacy settings
- [ ] Flutter settings pages load real data

## Deliverables

Files created/modified, endpoint inventory.
