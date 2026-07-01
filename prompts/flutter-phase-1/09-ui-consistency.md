# Horizon — Milestone 9: UI Consistency

**Phase 8 of 9** | **Estimated: 2 sessions** | **Status: Pending**

---

## CONTEXT REMINDER

**Architecture is frozen. Specifications are frozen. Technology stack is frozen.**

### Prior Work Completed

- **M1-M8**: All features implemented (Auth, Settings, Transactions, Search, AI Advisor, Insights, Planning, Notifications, Settings completion)
- **M3** created the initial shared widget library: `SharedLoadingView`, `SharedErrorView`, `SharedEmptyView`, `SharedOfflineView`, `SharedCard`, `SharedSectionHeader`, `SharedStatusBadge`, `SharedSearchBar`, `SharedFilterChip`
- **M3** added theme constants: `AppSpacing`, `AppRadius`, `AppElevation`, `AppIconSize`
- Initial migration was done in M3 for Dashboard, Accounts, Goals, Portfolio, Timeline, Settings

### Design System References

- `horizon-spec/products/DESIGN_SYSTEM.md` — THE primary reference for this milestone. Read the FULL document (501 lines). Key sections:
  - Design Principles (Calm, Trustworthy, Professional, Approachable, Confident, Predictable, Explainable, Minimal Cognitive Load)
  - Typography Philosophy (primary numbers are large and prominent, every number has a label)
  - Spacing Philosophy (cards have consistent padding, Tier 1 vs Tier 2 visual separation)
  - Visual Rhythm (cards repeat at consistent cadence, information appears in same positions)
  - Information Hierarchy (Primary Metric -> Label -> Change Indicator -> Confidence -> Secondary Info -> Context -> Evidence -> Actions)
  - Financial Language (INR formatting with Cr/L suffixes, percentages, dates, growth/loss terminology)
  - Colour Philosophy (functional not decorative: Success=green, Warning=amber, Critical=red, Neutral=grey, Information=blue)
  - Iconography (icons reinforce meaning but never carry meaning alone)
  - Card Philosophy (every card has purpose, primary metric, evidence, primary action)
  - Charts Philosophy (decision-support tools)
  - Empty States, Error States, Loading States specifications
  - Accessibility (WCAG 2.1 AA minimum)
- `horizon-spec/products/UX_PHILOSOPHY.md` — 10 Core Principles, Decision-First Design, Progressive Disclosure (Essential/Standard/Advanced/Expert)
- `horizon-spec/products/NAVIGATION_SYSTEM.md` — 3-tap max rule, context preservation

### Current App State
- All features functional
- Shared widgets exist but adoption may be incomplete
- Theme constants exist but usage across features may be inconsistent
- Some files may still have hardcoded padding/margin/colors

### Flutter App Location

`C:\Kannan\Horizon\horizon\apps\mobile\`

---

## STANDARD OPERATING PROCEDURE

1. **Full UI Audit** — Read EVERY page and widget file. Catalog every inconsistency.
2. **Implementation Plan** — List every file to modify.
3. **Implementation** — Apply fixes. NO new features. NO behavior changes.
4. **Build Verification** — `flutter analyze` — zero errors.
5. **Visual Verification** — Walk through every screen, confirm consistency.
6. **Final Report** — Adoption percentages, remaining issues.
7. **Message to Solution Architect**.

---

## Step 1: Full UI Audit

Read EVERY file in these directories. For EACH file, answer:
1. Does it use shared widgets from `lib/shared/widgets/`? (Yes/No/Partial)
2. Does it use `AppSpacing` constants? (Yes/No/Partial)
3. Does it use `AppRadius` constants? (Yes/No/Partial)
4. Does it use `AppElevation` constants? (Yes/No/Partial)
5. Does it use `AppIconSize` constants? (Yes/No/Partial)
6. Are there hardcoded color values not using theme? (List them)
7. Are there hardcoded padding/margin values? (List them)
8. Does the file follow the information hierarchy? (Primary Metric -> Label -> etc.)
9. Are money amounts formatted using Indian notation (Cr/L/K)?
10. Does every number have a label?

### Files to Audit

```
lib/features/auth/pages/
  splash_page.dart, login_page.dart, forgot_password_page.dart,
  reset_password_page.dart, session_expired_page.dart

lib/features/dashboard/pages/
  dashboard_page.dart
lib/features/dashboard/widgets/
  dashboard_cards.dart

lib/features/accounts/pages/
  accounts_page.dart, account_detail_page.dart
lib/features/accounts/widgets/
  acct_widgets.dart

lib/features/goals/pages/
  goals_page.dart, goal_detail_page.dart
lib/features/goals/widgets/
  goal_widgets.dart

lib/features/portfolio/pages/
  portfolio_page.dart
lib/features/portfolio/widgets/
  pf_widgets.dart

lib/features/timeline/pages/
  timeline_page.dart, timeline_detail_page.dart
lib/features/timeline/widgets/
  timeline_widgets.dart

lib/features/transactions/pages/
  transactions_page.dart, transaction_detail_page.dart, transaction_form_page.dart
lib/features/transactions/widgets/
  transaction_widgets.dart, transaction_filters.dart

lib/features/search/pages/
  search_page.dart
lib/features/search/widgets/
  search_widgets.dart (if any)

lib/features/advisor/pages/
  advisor_page.dart
lib/features/advisor/widgets/
  advisor_widgets.dart

lib/features/insights/pages/
  insights_page.dart, insight_detail_page.dart
lib/features/insights/widgets/
  insight_widgets.dart

lib/features/planning/pages/
  planning_page.dart, projection_page.dart, scenario_page.dart,
  scenario_compare_page.dart, planning_detail_page.dart
lib/features/planning/widgets/
  planning_widgets.dart

lib/features/notifications/pages/
  notifications_page.dart, notification_settings_page.dart
lib/features/notifications/widgets/
  notification_widgets.dart

lib/features/settings/pages/
  settings_page.dart, profile_page.dart, edit_profile_page.dart,
  preferences_page.dart, privacy_page.dart, security_page.dart,
  about_page.dart, ai_settings_page.dart, data_management_page.dart,
  notification_preferences_page.dart
```

---

## Step 2: Implementation

### Phase A: Shared Widget Adoption (Pass 1)

For each file that doesn't use shared widgets, migrate:
- Replace `Center(child: CircularProgressIndicator())` with `SharedLoadingView`
- Replace manual error views (icon + message + retry button) with `SharedErrorView`
- Replace manual empty states (icon + message) with `SharedEmptyView`
- Replace manual offline views with `SharedOfflineView`
- Wrap cards with `SharedCard` wrapper for consistent elevation/radius/clip
- Replace manual section headers with `SharedSectionHeader`
- Replace manual status badges with `SharedStatusBadge`
- Replace manual search bars with `SharedSearchBar`
- Replace manual filter chips with `SharedFilterChip`

### Phase B: Theme Constant Adoption (Pass 2)

Replace ALL hardcoded spacing values:
- `4` -> `AppSpacing.xs`
- `8` -> `AppSpacing.sm`
- `12` -> `AppSpacing.sm` (closest) or `AppSpacing.md` (if 16 is too much)
- `16` -> `AppSpacing.md`
- `24` -> `AppSpacing.lg`
- `32` -> `AppSpacing.xl`
- `48` -> `AppSpacing.xxl`

Replace ALL hardcoded border radius:
- `8` -> `AppRadius.sm`
- `12` -> `AppRadius.md`
- `16` -> `AppRadius.lg`
- `24` -> `AppRadius.xl`

Replace ALL hardcoded elevations:
- `1` -> `AppElevation.low`
- `2` -> `AppElevation.medium`
- `4` -> `AppElevation.high`

Replace ALL hardcoded icon sizes:
- `16` -> `AppIconSize.sm`
- `20` -> `AppIconSize.md`
- `24` -> `AppIconSize.lg`
- `32` -> `AppIconSize.xl`

### Phase C: Design System Compliance (Pass 3)

For each screen, verify:
1. Information hierarchy is correct (Primary Metric first, then label, then supporting info)
2. Every number has a label (no orphaned numbers)
3. Money uses Indian notation (Cr/L/K with ₹ prefix)
4. Functional colors are used correctly (green=positive, red=negative/critical, amber=warning, blue=info, grey=neutral)
5. Icons reinforce meaning but text also conveys the message (accessibility)
6. Touch targets are adequate (minimum 48x48 for interactive elements)
7. Cards have consistent padding (should use `AppSpacing.md` as the standard card padding)

### Phase D: Auth Pages Consistency (Pass 4)

Ensure auth pages match the overall app design:
- Consistent spacing using `AppSpacing`
- Consistent button styling (use `FilledButton` for primary, `OutlinedButton` for secondary)
- Consistent input decoration theme (already in `theme.dart`, verify it's applied)
- Same brand colors (teal seed color)

---

## Step 3: Build Verification

```powershell
cd C:\Kannan\Horizon\horizon\apps\mobile
flutter pub get
flutter analyze
```

Zero errors required.

---

## VERIFICATION CHECKLIST

- [ ] `flutter analyze` — zero errors
- [ ] Every screen uses shared widgets for loading/error/empty/offline states
- [ ] No hardcoded spacing values remaining (all use `AppSpacing`)
- [ ] No hardcoded border radius values remaining (all use `AppRadius`)
- [ ] No hardcoded elevation values remaining (all use `AppElevation`)
- [ ] No hardcoded icon sizes remaining (all use `AppIconSize`)
- [ ] Auth pages visually consistent with app design
- [ ] Every card has consistent padding and elevation
- [ ] Every section header is consistent
- [ ] Money formatted in Indian notation consistently
- [ ] Information hierarchy followed on all screens
- [ ] Accessibility: icons never carry meaning alone
- [ ] No behavior changes (visual diff only)

---

## DELIVERABLE TEMPLATE

```markdown
## 1. UI Audit Results
...
## 2. Shared Widget Adoption %
...
## 3. Theme Constant Adoption %
...
## 4. Files Modified (by phase)
...
## 5. Remaining Inconsistencies
...
## 6. Verification Checklist
...
## 7. Message to Solution Architect
```
