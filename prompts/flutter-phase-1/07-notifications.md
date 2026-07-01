# Horizon — Milestone 7: Notifications

**Phase 6 of 9** | **Estimated: 1 session** | **Status: Pending**

---

## CONTEXT REMINDER

**Architecture is frozen. Specifications are frozen. Technology stack is frozen.**

### Prior Work Completed

- **M1** Foundation Stabilization (Auth, Settings, Networking, Error handling)
- **M2** Transactions Experience (list, detail, create, archive, search, filters, pagination)
- **M3** Global Search + Shared Widgets (cross-module search, shared UI library, theme constants)
- **M4** AI Advisor (chat interface with markdown, suggested prompts, context awareness, provider switching)
- **M5** Insights (spending/income analysis, trends with fl_chart, opportunities, warnings, achievements, forecast)
- **M6** Financial Planning (dashboard hub, projections chart, scenarios, comparison, "coming soon" for missing budget/retirement/emergency/debt/investment endpoints)

### Design System References

- `horizon-spec/products/DESIGN_SYSTEM.md`: Calm Technology principle — the system is quiet by default, silence is a feature. Notification philosophy, information hierarchy.
- `horizon-spec/products/UX_PHILOSOPHY.md`: Non-judgmental language, confidence-building patterns, celebration design principles.
- `horizon-spec/products/NAVIGATION_SYSTEM.md`: Deep linking from notifications to specific goals, recommendations, risk alerts.

### Current App State
- `/dashboard/notifications` route exists with a PLACEHOLDER (7-line stub)
- Dashboard AppBar may show notification bell icon (needs wiring)
- Shared widgets available
- Theme constants available

### Flutter App Location

`C:\Kannan\Horizon\horizon\apps\mobile\`

---

## STANDARD OPERATING PROCEDURE

1. **Backend Audit** — Verify every required endpoint.
2. **Implementation Plan** — List every file.
3. **Implementation** — Build. NO PLACEHOLDERS. NO TODOs.
4. **Build Verification** — `flutter analyze` — zero errors.
5. **Functional Verification** — Checklist.
6. **Final Report**.
7. **Message to Solution Architect**.

---

## Step 1: Backend Audit

### Required Endpoints

Check `C:\Kannan\Horizon\horizon\services\experiences\notifications\register\register.go`:

| Method | Path | Purpose | Expected |
|--------|------|---------|----------|
| `GET` | `/api/v1/notifications` | Notification center (all) | ✅ EXISTS |
| `GET` | `/api/v1/notifications/unread` | Unread notifications | ✅ EXISTS |
| `GET` | `/api/v1/notifications/history` | Notification history | ✅ EXISTS |
| `GET` | `/api/v1/notifications/preferences` | Get notification preferences | ✅ EXISTS |
| `PUT` | `/api/v1/notifications/preferences` | Update preferences | ✅ EXISTS |
| `GET` | `/api/v1/notifications/search?q=` | Search notifications | ✅ EXISTS |
| `GET` | `/api/v1/notifications/{id}` | Single notification | ✅ EXISTS |
| `POST` | `/api/v1/notifications/{id}/read` | Mark as read | ✅ EXISTS |
| `POST` | `/api/v1/notifications/{id}/archive` | Archive | ✅ EXISTS |
| `POST` | `/api/v1/notifications/{id}/snooze` | Snooze | ✅ EXISTS |

**Missing endpoints (document but proceed):**
| Endpoint | Status |
|----------|--------|
| `DELETE /api/v1/notifications/{id}` | ❌ DOES NOT EXIST |
| `POST /api/v1/notifications/register-push-token` | ❌ DOES NOT EXIST |

Record exact request/response shapes.

---

## Step 2: Implementation Plan

Structure:
```
features/notifications/
├── models/notification_models.dart
├── repository/notification_repository.dart
├── state/notification_state.dart
├── pages/
│   ├── notifications_page.dart     (REPLACE placeholder)
│   └── notification_settings_page.dart
└── widgets/notification_widgets.dart
```

Also modify:
- `router.dart` — update route if needed
- `dashboard_page.dart` or `shell.dart` — add unread badge to notification bell icon

---

## Step 3: Implementation

### Models

- `NotificationCenterResponse`, `NotificationItem`, `NotificationPreference`, `NotificationCategory`
- Follow backend response shapes

### Repository

- `getNotifications()`, `getUnread()`, `getHistory()`
- `markRead(id)`, `archive(id)`, `snooze(id, until)`
- `getPreferences()`, `updatePreferences(prefs)`
- `search(query)`

### State

- `NotificationState` with `copyWith`: `notifications[]`, `unreadCount`, `activeTab` (all/unread/history), `preferences[]`, `status`, `error`
- `NotificationNotifier`: `load()`, `loadUnread()`, `loadHistory()`, `markRead()`, `archive()`, `snooze()`, `search()`, `loadPreferences()`, `updatePreference()`

### Pages

#### Notifications Page (`features/notifications/pages/notifications_page.dart`)

REPLACE placeholder.

- Tab bar: All | Unread | History
- Each notification: icon + title + body + time + read/unread indicator
- Swipe actions: swipe left to archive, swipe right to mark read
- Long press: menu with archive, snooze options
- Pull-to-refresh
- Empty state: "No notifications"
- Loading/error/offline states

**Unread badge on Dashboard:**
- Add notification bell icon to Dashboard AppBar
- Show unread count badge
- Tap navigates to `/dashboard/notifications`

#### Notification Settings Page (`features/notifications/pages/notification_settings_page.dart`)

- Per-category preference toggles (in_app, push, email per category)
- Snooze/digest frequency settings
- Load/save via preferences API

### Widgets

- `NotificationCard` — list item with swipe actions
- `NotificationCategoryBadge` — colored category badge
- `NotificationSettingsTile` — preference toggle per category

---

## Step 4: Build Verification

```powershell
cd C:\Kannan\Horizon\horizon\apps\mobile
flutter pub get
flutter analyze
```

Zero errors required.

---

## VERIFICATION CHECKLIST

- [ ] `flutter analyze` — zero errors
- [ ] Notification page loads with tabs (All/Unread/History)
- [ ] Swipe left archives notification
- [ ] Swipe right marks as read
- [ ] Long press shows action menu
- [ ] Pull-to-refresh works
- [ ] Search with debounce works
- [ ] Notification bell icon on Dashboard shows unread count
- [ ] Tap bell navigates to notifications
- [ ] Notification settings loads and saves
- [ ] Empty/loading/error/offline states work
- [ ] Shared widgets used for common states

---

## DELIVERABLE TEMPLATE

```markdown
## 1. API Inventory
...
## 2. Screens Implemented
...
## 3. Files Created/Modified
...
## 4. Missing Backend Endpoints
...
## 5. Verification Checklist
...
## 6. Message to Solution Architect
```
