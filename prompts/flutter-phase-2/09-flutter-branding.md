# Phase 2 — Milestone 9: Flutter Branding

**Context Reminder:** Architecture frozen. Specs frozen. Stack frozen. Monorepo. No Redis/K8s/Grafana/testing.

**Milestones completed:** M1-M8 (all backend + Flutter polish done).

**Current state:** App uses default Flutter icon. Splash screen shows logo + loading. No custom illustrations, empty-state graphics, or branded assets.

## Objective

Set up proper app icon, branded splash screen, and custom illustrations for empty states.

## Implementation

### App Icon

Use `flutter_launcher_icons` package (add to pubspec.yaml):
- Create a simple app icon based on the teal `#00897B` brand color
- Use the trending-up icon as the base
- Run `flutter pub run flutter_launcher_icons` to generate all sizes

### Splash Screen

- Update `android/app/src/main/res/drawable/launch_background.xml` with brand color
- Update `ios/Runner/LaunchScreen.storyboard` with brand background
- Keep the animated splash in Flutter (splash_page.dart) — it's fine

### Empty-State Illustrations

Create simple custom illustration widgets in `lib/shared/widgets/illustrations.dart`:
- `NoDataIllustration` — for empty lists (receipt icon for transactions, flag for goals, etc.)
- `ErrorIllustration` — for error states (cloud-off icon)
- `SearchIllustration` — for search empty states (search icon)
- `NetworkOffIllustration` — for offline states (wifi-off icon)

These should use the brand color palette and be consistent across the app. Currently we use plain Material icons — wrap them in a consistent container with brand-colored background and standard sizing.

### Color & Typography

- Verify all screens use `theme.colorScheme` consistently (no hardcoded colors)
- Verify the teal seed color `#00897B` is the only brand color source
- Ensure financial numbers use the `MoneyFormat` extension consistently

## Files to Modify

| File | Change |
|------|--------|
| `pubspec.yaml` | Add `flutter_launcher_icons` |
| `flutter_launcher_icons.yaml` | NEW — icon config |
| `lib/shared/widgets/illustrations.dart` | NEW — custom illustration widgets |
| `lib/shared/extensions/extensions.dart` | Enhance MoneyFormat |
| Various pages | Replace inline icon+text with illustration widgets |

## Verification

- [ ] App icon shows on home screen
- [ ] Splash screen shows brand color
- [ ] Empty states use branded illustrations
- [ ] All screens use theme colors consistently
- [ ] `flutter analyze` — zero errors

## Deliverables

App icon, splash, illustration widgets.
