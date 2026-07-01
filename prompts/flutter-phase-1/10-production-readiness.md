# Horizon — Milestone 10: Production Readiness

**Phase 9 of 9** | **Estimated: 1-2 sessions** | **Status: Pending**

---

## CONTEXT REMINDER

**Architecture is frozen. Specifications are frozen. Technology stack is frozen.**

### Prior Work Completed

- **M1** Foundation Stabilization: Auth, Settings, Networking (connectivity_plus, ErrorMapper, retry interceptor), Error handling
- **M2** Transactions Experience: List, detail, create, archive, search, filters, pagination, infinite scroll
- **M3** Global Search + Shared Widgets: Cross-module search, shared widget library, theme constants
- **M4** AI Advisor: Chat interface with markdown, suggested prompts, context awareness, provider switching
- **M5** Insights: Spending/income analysis, trends with fl_chart, opportunities, warnings, achievements, forecast
- **M6** Financial Planning: Dashboard hub, projections chart, scenarios, comparison, "coming soon" for missing endpoints
- **M7** Notifications: Center with tabs, swipe actions, preferences, search, unread badge on dashboard
- **M8** Settings & Profile Completion: All settings pages complete, no-op buttons fixed, new pages for AI settings, data management, notification preferences
- **M9** UI Consistency: Shared widget adoption complete, theme constants used everywhere, design system compliance verified

### Design System References

- `horizon-spec/products/DESIGN_SYSTEM.md` — All design principles should be fully realized at this point.
- `horizon-spec/engineering/FRONTEND_ENGINEERING_GUIDE.md` — Offline behaviour, performance principles, testing philosophy.

### Current App State
- All features implemented and visually consistent
- Shared widget library fully adopted
- Theme constants used everywhere

### Flutter App Location

`C:\Kannan\Horizon\horizon\apps\mobile\`

---

## STANDARD OPERATING PROCEDURE

1. **Full Production Audit** — Check every item below.
2. **Implementation Plan** — List every fix.
3. **Implementation** — Apply fixes.
4. **Build Verification** — `flutter analyze` — zero errors, zero warnings. `flutter build` — success.
5. **Final Verification** — Full smoke test of every screen.
6. **Final Report** — Production readiness score.
7. **Message to Solution Architect**.

---

## Step 1: Production Audit

### 1A: Placeholder & Debug Removal

Check every file for:
- [ ] Any remaining placeholder pages (Advisor, Insights, Planning, Notifications should all be real now)
- [ ] Any "Coming soon" screens that should have been implemented
- [ ] `debugShowCheckedModeBanner: false` in `app.dart` (verify it's there)
- [ ] `debugLogDiagnostics: false` in `router.dart` (verify)
- [ ] Any `print()` calls in production code (replace with proper logging or remove)
- [ ] Any mock data or mock providers
- [ ] Any TODOs or FIXMEs in the codebase
- [ ] Any commented-out code

### 1B: Performance

- [ ] Run `flutter pub deps` — check for unused packages
- [ ] Run `flutter clean` then `flutter pub get`
- [ ] Check list views for proper `itemBuilder` (no unnecessary rebuilds)
- [ ] Check Riverpod providers for proper `autoDispose` where applicable
- [ ] Check `const` constructors (add where missing)
- [ ] Verify pagination (timeline, transactions) doesn't degrade with 1000+ items
- [ ] Check image/assets loading
- [ ] Verify `build_runner` is not needed (or run it if freezed/json_serializable are actually used)

### 1C: Network & Error Handling

- [ ] Every screen handles: loading, error, empty, offline
- [ ] Token refresh works end-to-end (test by waiting for 401 or simulating it)
- [ ] 401 responses trigger refresh -> retry -> success
- [ ] If refresh fails, redirect to `/session-expired`
- [ ] Network errors show user-friendly messages (not raw `DioException` text)
- [ ] Connectivity changes are reflected in UI (offline -> online auto-retry)
- [ ] Request timeouts are configured (verify `AppConfig.connectTimeout` and `receiveTimeout`)

### 1D: Security

- [ ] All tokens stored in `FlutterSecureStorage` (no `SharedPreferences` for tokens)
- [ ] Logout clears all local data (tokens, cache)
- [ ] No sensitive data in logs (check `print()` calls)
- [ ] No hardcoded API keys or secrets
- [ ] No tokens in URL query params (only in `Authorization` header)

### 1E: Error Boundary

- [ ] App has a global error boundary to catch unhandled exceptions
- [ ] Flutter error widget is customized (not the red error screen)
- [ ] `FlutterError.onError` is configured
- [ ] `PlatformDispatcher.instance.onError` is configured

### 1F: Build Configuration

- [ ] `pubspec.yaml` version is correct (`version: 0.1.0-dev` or updated for release)
- [ ] `environment.dart` production config is correct (`apiBaseUrl`, `enableLogging: false`)
- [ ] App icon and splash screen are correct (check `flutter pub run flutter_launcher_icons` if configured)
- [ ] Android: check `AndroidManifest.xml` for internet permission, app label
- [ ] iOS: check `Info.plist` for app name, permissions

### 1G: Accessibility

- [ ] All images have semantic labels
- [ ] Touch targets are minimum 48x48
- [ ] Color contrast meets WCAG 2.1 AA (not just color to convey info)
- [ ] Screen reader support (Semantics widget where needed)
- [ ] System font size changes don't break layouts

---

## Step 2: Implementation

Apply all fixes found in the audit. Work through each category:

1. Remove placeholders, debug code, print statements
2. Fix performance issues (const constructors, autoDispose)
3. Fix any missing error/offline handling
4. Fix any security issues
5. Add error boundary if missing
6. Update build configuration
7. Fix accessibility issues

---

## Step 3: Build Verification

```powershell
cd C:\Kannan\Horizon\horizon\apps\mobile

# Clean build
flutter clean
flutter pub get
flutter analyze

# Android build (if SDK available)
flutter build apk --release --target-platform android-arm64

# Run tests
flutter test
```

`flutter analyze` — zero errors, zero warnings.
`flutter build` — builds successfully.
`flutter test` — all tests pass.

---

## Step 4: Final Smoke Test

Run through EVERY screen in the app:

1. Splash -> Login -> Forgot Password -> Reset Password -> Session Expired
2. Dashboard: all widget cards load
3. Goals: list -> detail -> progress
4. Portfolio: allocation -> performance -> risk
5. Timeline: list -> filter -> search -> detail
6. Transactions: list -> search -> filter -> add -> detail -> archive
7. More -> Planning: dashboard -> projections -> scenarios
8. More -> Accounts: list -> detail
9. More -> Insights: feed -> tabs -> trends -> search -> detail
10. More -> Advisor: chat -> suggested prompts -> send message -> receive response
11. More -> Notifications: tabs -> swipe actions -> preferences
12. Settings: profile -> edit -> preferences -> privacy -> security -> AI settings -> about
13. Global Search: open -> search -> navigate to result
14. Logout -> verify redirect to login
15. Offline: disable network -> verify all screens show offline state
16. Session expiry: simulate 401 -> verify refresh -> verify session-expired page

---

## DELIVERABLE TEMPLATE

```markdown
## 1. Production Audit Results
### Placeholder/Debug Issues Found
...
### Performance Issues Found
...
### Network/Error Issues Found
...
### Security Issues Found
...
### Build Configuration Issues Found
...
### Accessibility Issues Found
...

## 2. Issues Fixed
...

## 3. Remaining Issues (known limitations)
...

## 4. Production Readiness Score: X/10
...

## 5. Build Verification
- flutter analyze: ✅ (0 errors, 0 warnings)
- flutter build apk: ✅
- flutter test: ✅ / ⚠️ (N tests pass)

## 6. Files Modified
...

## 7. Message to Solution Architect
```
