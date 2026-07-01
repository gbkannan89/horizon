# Phase 2 — Milestone 7: Notifications Completion

**Context Reminder:** Architecture frozen. Specs frozen. Stack frozen. Monorepo. No Redis/K8s/Grafana/testing.

**Milestones completed:** M1-M6 (all backend gaps filled except notification remaining items).

**Current state:** Notifications working with list, unread, history, mark-read, archive, snooze, search. Missing: push notification registration, delete notification.

## Objective

Add push notification registration endpoint + delete notification endpoint.

## Backend Audit

Read:
- `services/experiences/notifications/register/register.go` — current routes
- `services/experiences/notifications/internal/api/handlers.go` — existing handlers

Confirm:
- `POST /api/v1/notifications/register-token` — ❌ DOES NOT EXIST
- `DELETE /api/v1/notifications/{id}` — ❌ DOES NOT EXIST

## Implementation

### Backend

In `services/experiences/notifications/register/register.go`:

1. `POST /api/v1/notifications/register-token` — accepts `{token, platform}`, stores in DB. For MVP, log and return 200.
2. `DELETE /api/v1/notifications/{id}` — hard deletes the notification from DB.

### Flutter

- `notifications_page.dart` — add "delete" action to long-press menu (alongside archive and snooze)
- Add push token registration on app startup (in main.dart or app.dart) using `device_info_plus` or a simple placeholder token for MVP

## Files to Modify

| File | Change |
|------|--------|
| `services/experiences/notifications/register/register.go` | Add 2 routes |
| `services/experiences/notifications/internal/api/handlers.go` | Add 2 handlers |
| `notifications_repository.dart` | Add `deleteNotification()` method |
| `notification_state.dart` | Add `deleteNotification()` method |
| `notifications_page.dart` | Add delete action |
| `main.dart` or `app.dart` | Add push token registration call |

## Verification

- [ ] `go build ./cmd/server/...` passes
- [ ] `POST /notifications/register-token` returns 200
- [ ] `DELETE /notifications/{id}` removes notification
- [ ] `flutter analyze` — zero errors
- [ ] Long-press menu shows Delete option
- [ ] Delete removes notification from list

## Deliverables

Push registration, delete endpoint, Flutter delete action.
