# Phase 2 — Milestone 3: Transaction Edit & Delete

**Context Reminder:** Architecture frozen. Specs frozen. Stack frozen. Monorepo. No Redis/K8s/Grafana/testing.

**Milestones completed before this:** M1 (Docker infra), M2 (Auth & Profile endpoints).

**Current state:** Transactions have LIST, DETAIL, CREATE (mock), ARCHIVE. No EDIT or DELETE.

## Objective

Implement `PUT /api/v1/transactions/{id}` and `DELETE /api/v1/transactions/{id}` in the backend. Add Edit and Delete UIs in Flutter.

## Backend Audit

Read:
- `services/domains/financial-event/internal/application/service.go` — has `UpdateEvent`, `CancelCommand` in domain layer?
- `services/domains/financial-event/internal/application/dto/command/commands.go` — check for Update/Cancel/Delete commands
- `services/domains/financial-event/register/register.go` — currently has GET endpoints only
- `services/domains/financial-event/internal/domain/entity.go` — FinancialEvent entity with state machine

Confirm what domain commands already exist:
- `UpdateEvent` or similar command — check commands.go
- `CancelEvent` or `DeleteEvent` — check if cancel maps to delete

## Implementation

### Backend

In `services/domains/financial-event/register/register.go`:
- `PUT /api/v1/transactions/{id}` — accepts `{description, category, notes, amount}`, updates the event
- `DELETE /api/v1/transactions/{id}` — soft-delete (sets state to cancelled/archived) or hard delete

Use existing domain service methods if they exist. If domain layer doesn't have update/delete, create simple handlers that update the database directly (MVP pattern matching existing code).

### Flutter

- `transaction_form_page.dart` — currently only "Add" mode. Add "Edit" mode: accepts optional `TransactionEvent`, pre-fills form, calls PUT on save instead of POST
- `transaction_detail_page.dart` — add "Edit" button in AppBar that navigates to form in edit mode
- `transaction_detail_page.dart` — add "Delete" button in PopupMenu with confirmation dialog, calls DELETE API

## Files to Modify

| File | Change |
|------|--------|
| `services/domains/financial-event/register/register.go` | Add PUT + DELETE handlers |
| `transaction_form_page.dart` | Support edit mode with pre-filled fields |
| `transaction_detail_page.dart` | Add edit + delete actions |
| `transaction_repository.dart` | Add `updateTransaction()` + `deleteTransaction()` methods |

## Verification

- [ ] `go build ./cmd/server/...` — backend builds
- [ ] `PUT /transactions/{id}` updates fields
- [ ] `DELETE /transactions/{id}` removes (or archives) the transaction
- [ ] `flutter analyze` — zero errors
- [ ] Detail page shows Edit button
- [ ] Edit form pre-fills with current values
- [ ] Save calls PUT and refreshes
- [ ] Delete shows confirmation dialog
- [ ] Confirm delete calls DELETE and navigates back

## Deliverables

Endpoints created, Flutter edit/delete flow.
