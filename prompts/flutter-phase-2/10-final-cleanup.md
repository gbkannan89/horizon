# Phase 2 — Milestone 10: Final Cleanup

**Context Reminder:** Architecture frozen. Specs frozen. Stack frozen. Monorepo. No Redis/K8s/Grafana/testing.

**Milestones completed:** M1-M9 (all features + polish + branding done).

**Current state:** Full-featured app. Need to clean up unused code, verify everything builds, and ensure production readiness.

## Objective

Remove unused packages, clean dead code, verify final build.

## Implementation

### Unused Packages

Check each package in `pubspec.yaml` — is it actually imported anywhere?

| Package | Used? | Action |
|---------|-------|--------|
| `flutter_hooks` | Run grep — if unused, remove | Remove if unused |
| `intl` | Run grep — if unused, remove | Remove if unused |
| `freezed_annotation` | Check if any freezed files exist | Remove if unused |
| `json_annotation` | Check if `@JsonSerializable` used | Keep if used |
| `flutter_secure_storage` | ✅ Used by TokenManager | Keep |
| `fl_chart` | ✅ Used by insights + planning charts | Keep |
| `flutter_markdown` | ✅ Used by Advisor chat | Keep |
| `connectivity_plus` | ✅ Used by ConnectivityMonitor | Keep |
| `local_auth` | Run grep — if unused, remove | Remove if unused |
| `freezed`, `json_serializable`, `riverpod_generator`, `build_runner` | Dev deps — if no generated files, remove | Remove if unused |
| `mocktail` | Dev dep for testing — remove since no tests | Remove |

### Dead Code

- Remove any unused imports across all Dart files (`flutter analyze` will show these)
- Remove any unused helper functions (like `_fmt` duplicates in accounts)
- Remove unused `print()` calls (the 3 in `api_client.dart` are logging-gated — keep them)
- Remove the empty `lib/shared/models/` and `lib/shared/utils/` directories

### Final Build

```powershell
flutter clean
flutter pub get
flutter analyze
```

Confirm zero errors, zero warnings.

## Files to Modify

| File | Change |
|------|--------|
| `pubspec.yaml` | Remove unused packages |
| Various `.dart` files | Fix unused import warnings |
| Remove empty directories | `lib/shared/models/`, `lib/shared/utils/` |

## Verification

- [ ] `flutter pub get` succeeds
- [ ] `flutter analyze` — zero errors, zero warnings
- [ ] App builds: `flutter build apk --debug`
- [ ] No runtime crashes on navigation
- [ ] All removed packages are confirmed unused

## Deliverables

Clean pubspec.yaml, clean codebase, final build verification.
