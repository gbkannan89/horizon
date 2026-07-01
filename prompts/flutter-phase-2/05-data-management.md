# Phase 2 — Milestone 5: Data Management

**Context Reminder:** Architecture frozen. Specs frozen. Stack frozen. Monorepo. No Redis/K8s/Grafana/testing.

**Milestones completed:** M1 (Docker), M2 (Auth/Profile), M3 (Transaction Edit/Delete), M4 (Planning endpoints).

**Current state:** Data Management page shows 4 "coming soon" cards for Export, Import, Backup, Restore.

## Objective

Implement simple data export/import/backup/restore endpoints and wire the Flutter UI.

## Backend Audit

No existing endpoints for data management. Need to create new handlers.

## Implementation

### Backend

Create `services/experiences/data/register/register.go`:

1. `GET /api/v1/data/export` — returns JSON dump of user's data: accounts, transactions, goals, portfolio. Accepts `?format=json|csv`. For MVP, JSON only.
2. `POST /api/v1/data/import` — accepts JSON file, creates entities. For MVP, just log received data and return success.
3. `POST /api/v1/data/backup` — triggers export and saves to server temp storage. Returns backup_id.
4. `POST /api/v1/data/restore` — accepts backup_id, restores from saved backup. For MVP, return success message.

Wire into monolith: `dataReg.RegisterRoutes(mux, pool)` in `cmd/server/main.go`.

### Flutter

Replace the "coming soon" cards in `data_management_page.dart`:
- "Export My Data" — calls GET /data/export, saves response to device storage, shows share sheet
- "Import Data" — opens file picker, sends file to POST /data/import
- "Backup" — calls POST /data/backup, shows backup ID
- "Restore" — shows text field for backup ID, calls POST /data/restore

## Files to Create/Modify

| File | Change |
|------|--------|
| `services/experiences/data/register/register.go` | NEW |
| `services/experiences/data/internal/api/handlers.go` | NEW |
| `cmd/server/main.go` | Add dataReg import |
| `data_management_page.dart` | Wire to real APIs |

## Verification

- [ ] `go build ./cmd/server/...` passes
- [ ] `GET /data/export` returns JSON
- [ ] `POST /data/import` accepts data
- [ ] `POST /data/backup` returns ID
- [ ] `flutter analyze` — zero errors
- [ ] Data management page shows working buttons

## Deliverables

Backend endpoints, working data management page.
