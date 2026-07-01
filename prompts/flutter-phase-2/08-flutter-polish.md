# Phase 2 — Milestone 8: Flutter Polish

**Context Reminder:** Architecture frozen. Specs frozen. Stack frozen. Monorepo. No Redis/K8s/Grafana/testing.

**Milestones completed:** M1-M7 (all backend gaps filled).

**Current state:** All features functional. UI is clean but lacks animations, haptic feedback, and responsive/tablet layouts.

## Objective

Add page transitions, card entrance animations, haptic feedback, and basic tablet responsiveness.

## Implementation

### Page Transitions

In `lib/app/router.dart`, customize GoRouter transitions:
- Add `CustomTransitionPage` for smooth slide/fade transitions between screens
- Dashboard tabs: no transition (instant)
- Push navigation: slide from right
- Modal routes (settings, forms): slide up from bottom

### Card Entrance Animations

- Create `AnimatedCard` shared widget that fades in + slides up when entering viewport
- Use `AnimationController` with `SlideTransition` + `FadeTransition`
- Wrap cards in list views with staggered delay (50ms per card)
- Apply to: transaction list, goals list, timeline list, insight cards

### Haptic Feedback

- `FilledButton` taps: `HapticFeedback.lightImpact()`
- Swipe-to-archive (notifications): `HapticFeedback.mediumImpact()`
- Pull-to-refresh: `HapticFeedback.heavyImpact()` on completion

Add in a shared widget wrapper or extension method.

### Responsive Layouts

- Check all pages for hardcoded widths (e.g., `Container(width: 96)`)
- Replace with `MediaQuery` or `LayoutBuilder` where appropriate
- Ensure minimum touch targets (48x48) on all interactive elements
- Test basic tablet layout (screens should not look stretched)

## Files to Modify

| File | Change |
|------|--------|
| `lib/app/router.dart` | Custom page transitions |
| `lib/shared/widgets/shared_widgets.dart` | AnimatedCard wrapper |
| `lib/features/*/pages/*.dart` | Apply AnimatedCard, haptic feedback |
| Select pages | Fix hardcoded widths for responsive |

## Verification

- [ ] `flutter analyze` — zero errors
- [ ] Page transitions animate smoothly
- [ ] List items animate in on scroll
- [ ] Buttons trigger haptic feedback
- [ ] App layout works on tablet (doesn't break)
- [ ] Touch targets are minimum 48x48

## Deliverables

Smooth animations, haptic feedback, responsive foundations.
