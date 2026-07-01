# Horizon — Milestone 8: Settings & Profile Completion

**Phase 7 of 9** | **Estimated: 1 session** | **Status: Pending**

---

## CONTEXT REMINDER

**Architecture is frozen. Specifications are frozen. Technology stack is frozen.**

### Prior Work Completed

- **M1** Foundation Stabilization: Settings baseline (profile, edit profile with API, preferences persistence, theme persistence, security page with disabled change-password, logout-all-devices fix)
- **M2-M7**: All core features completed (Transactions, Search + Shared Widgets, AI Advisor, Insights, Financial Planning, Notifications)

### Settings Already Built (from M1)

- `settings_page.dart` — hub with sections (Profile, Preferences, Privacy & Security, About), Sign Out, theme picker
- `profile_page.dart` — loads from API, displays avatar, name, email, phone, country, currency, timezone, language, member since
- `edit_profile_page.dart` — loads current profile, edits name/phone, calls `PUT /users/me` with loading/error/success states
- `preferences_page.dart` — loads/saves currency, language, notification toggles, compact timeline
- `privacy_page.dart` — loads privacy settings, toggles for data sharing/analytics/personalization/third-party/marketing, export data (no-op), delete account (dialog only)
- `security_page.dart` — change password (disabled), active sessions (no-op), logout all devices (calls repo.logout())
- `about_page.dart` — app version, licenses, privacy policy, terms links

### Design System References

- `horizon-spec/products/DESIGN_SYSTEM.md`: Information hierarchy, functional colors, accessibility (WCAG 2.1 AA minimum).
- `horizon-spec/products/NAVIGATION_SYSTEM.md`: Settings is secondary navigation — accessible from More page.
- `horizon-spec/engineering/FRONTEND_ENGINEERING_GUIDE.md`: Presentation Without Business Logic, State Isolation.

### Flutter App Location

`C:\Kannan\Horizon\horizon\apps\mobile\`

---

## STANDARD OPERATING PROCEDURE

1. **Backend Audit** — Verify what exists vs missing.
2. **Implementation Plan** — List every file.
3. **Implementation** — Build. NO PLACEHOLDERS. For features with missing endpoints, show "coming soon" with clear message. Do NOT leave no-op buttons.
4. **Build Verification** — `flutter analyze` — zero errors.
5. **Functional Verification** — Checklist.
6. **Final Report**.
7. **Message to Solution Architect**.

---

## Step 1: Backend Audit

### Check These Endpoints

| Endpoint | Expected Finding | Action |
|----------|------------------|--------|
| `POST /api/v1/auth/change-password` | ❌ DOES NOT EXIST anywhere | Keep disabled with message |
| `PUT /api/v1/users/{id}/profile` | ✅ Domain code exists (`UpdateProfile`) but ❌ no HTTP route | Document gap |
| `PUT /api/v1/users/{id}/preferences` | ✅ Domain code exists (`UpdatePreferences`) but ❌ no HTTP route | Document gap |
| `PUT /api/v1/users/{id}/privacy` | ✅ Domain code exists (`UpdatePrivacy`) but ❌ no HTTP route | Document gap |
| `POST /api/v1/users/{id}/consent` | ✅ Domain code exists but ❌ no HTTP route | Document gap |
| `GET /api/v1/users/{id}/export` | ❌ DOES NOT EXIST | Show "coming soon" |
| `POST /api/v1/users/{id}/import` | ❌ DOES NOT EXIST | Show "coming soon" |
| `POST /api/v1/users/{id}/archive` | ✅ EXISTS (API gateway route) | Use for delete account |
| `GET /api/v1/ai/providers` | ✅ EXISTS (AI service) | Use for AI settings |
| `GET /api/v1/notifications/preferences` | ✅ EXISTS | Use for notification prefs page |

---

## Step 2: Implementation Plan

### Files to Modify

1. `privacy_page.dart` — wire Export My Data (if endpoint exists or show coming soon), wire Delete Account to `POST /api/v1/users/{id}/archive`
2. `security_page.dart` — keep change-password disabled with clear message, wire Active Sessions (or show coming soon), ensure logout-all works
3. `edit_profile_page.dart` — add date picker for Date of Birth field, add country/currency/language pickers
4. `preferences_page.dart` — wire currency and language pickers to open selection dialogs

### New Files to Create

5. `ai_settings_page.dart` — AI provider selection, show active provider, allow switching
6. `data_management_page.dart` — data export, import, backup, restore (show coming soon if backend endpoints missing)
7. `notification_preferences_page.dart` — connect to `/api/v1/notifications/preferences` for per-category notification settings

### Router Changes

8. `router.dart` — add routes for new settings pages

---

## Step 3: Implementation

### Fix Existing Pages

#### Privacy Page
- "Export My Data" button: If no export endpoint exists, show coming soon dialog: "Data export will be available in a future update."
- "Delete Account" button: Confirm dialog -> call `POST /api/v1/users/{id}/archive` -> on success show message and logout -> on error show error

#### Security Page
- "Change Password": Keep disabled with subtitle "Feature coming soon. Password management will be available in a future update."
- "Active Sessions": If no endpoint, show "Coming soon"
- "Logout All Devices": Already works (calls repo.logout())

#### Edit Profile (enhance)
- Date of Birth field: Replace read-only text field with tappable field that opens `showDatePicker`
- Country: Add dropdown picker with common countries (India, USA, UK, UAE, Singapore, Canada, Australia)
- Currency: Add dropdown picker (INR, USD, EUR, GBP, AED, SGD, CAD, AUD)
- Language: Add dropdown picker (English, Hindi, Tamil, Telugu, Kannada, Malayalam, Marathi, Gujarati, Bengali)

#### Preferences
- Currency tile: Open bottom sheet picker, on selection call `updatePreferences`
- Language tile: Open bottom sheet picker, on selection call `updatePreferences`

### New Pages

#### AI Settings Page (`features/settings/pages/ai_settings_page.dart`)
- Load available providers from `GET /api/v1/ai/providers`
- Show currently active provider
- Allow user to switch provider (calls `POST /api/v1/ai/provider/switch`)
- Show provider capabilities (icon, description)

#### Data Management Page (`features/settings/pages/data_management_page.dart`)
- "Export My Data" — show coming soon or call endpoint if it exists
- "Import Data" — show coming soon
- "Backup" — show coming soon
- "Restore" — show coming soon
- Each card: icon + title + description + status badge (Available / Coming Soon)

#### Notification Preferences Page (`features/settings/pages/notification_preferences_page.dart`)
- If already created in M7, enhance it
- Load from `GET /api/v1/notifications/preferences`
- Per-category: in_app toggle, push toggle, email toggle, digest frequency dropdown
- Save on change via `PUT /api/v1/notifications/preferences`

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
- [ ] Edit profile: Date of Birth opens date picker
- [ ] Edit profile: Country/Currency/Language pickers work
- [ ] Edit profile: Save updates backend
- [ ] Preferences: Currency picker opens and saves
- [ ] Preferences: Language picker opens and saves
- [ ] Privacy: Export My Data shows coming soon or actual export
- [ ] Privacy: Delete Account confirms, calls archive API, logs out
- [ ] Security: Change Password shows disabled with message
- [ ] Security: Active Sessions shows coming soon or actual data
- [ ] Security: Logout All Devices works
- [ ] AI Settings: Providers load, switch works
- [ ] Data Management: Each card shows appropriate status
- [ ] Notification Preferences: Loads and saves
- [ ] All no-op buttons fixed (no more silent failures)

---

## DELIVERABLE TEMPLATE

```markdown
## 1. Backend Endpoints Used
...
## 2. Missing Backend Endpoints Documented
...
## 3. Files Modified
...
## 4. Files Created
...
## 5. Fixed Issues (previous no-op buttons)
...
## 6. Verification Checklist
...
## 7. Message to Solution Architect
```
