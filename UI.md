# Horizon UI/UX Design System & Theme Guidelines

This document defines the uniform UI/UX specifications for the Horizon mobile application to ensure a highly polished, consistent, and premium aesthetic. All developer agents must strictly adhere to these standards when creating or modifying screens, widgets, and layouts.

---

## 🎨 1. Design Aesthetic & Visual Identity

The design language of Horizon is built on a **soft, premium glassmorphism** style featuring clean pastel overlays, high-contrast typography, and a modern, friendly interface.

### Color Palette

*   **Background Screen**: Off-white / light slate blue (`#F1F5F9` to `#F8FAFC`). Do not use pure white or solid dark grey backgrounds.
*   **Hero Card Gradients**: Warm lavender-indigo gradients:
    *   Start: `#818CF8` (Indigo 400)
    *   End: `#C7D2FE` (Indigo 200)
*   **Text Colors**:
    *   Primary Text: Deep Slate/Blue-grey (`#0F172A` or `#1E293B`).
    *   Secondary/Muted Text: Cool Grey/Slate (`#64748B`).
*   **Pastel Accents & Indicator Colors**:
    *   Active/Teal: `#0D9488` (Teal 600) with background `#E6F4F1`
    *   Orange/Warning: `#EA580C` (Warning/Orange) with background `#FDF2E9`
    *   Purple/Accent: `#7C3AED` (Purple 600) with background `#F5F3FF`
    *   Blue/Info: `#2563EB` (Blue 600) with background `#EFF6FF`

---

## 📐 2. Card Layouts & Containers

All content sections, lists, and summary cards must be structured using highly rounded, frosted glass-like containers with soft shadows.

### Standard Section Card Structure
*   **Border Radius**: `24` or `28` (highly rounded corners).
*   **Border**: Thin white/semi-transparent border to simulate glass refraction.
*   **Shadow**: Large, very soft blur shadow (low opacity).

```dart
// Standard Frosted Glassmorphism Card
Container(
  decoration: BoxDecoration(
    color: Colors.white.withValues(alpha: 0.85),
    borderRadius: BorderRadius.circular(24),
    border: Border.all(
      color: Colors.white.withValues(alpha: 0.6),
      width: 1.2,
    ),
    boxShadow: [
      BoxShadow(
        color: Colors.black.withValues(alpha: 0.03),
        blurRadius: 20,
        offset: const Offset(0, 8),
      ),
    ],
  ),
  child: ...
)
```

---

## 📊 3. Progress Bars & Badges

### Thick Pill Progress Bars
Progress indicators must be thick, smooth, and rounded, utilizing color coding that matches the card's accent.

```dart
// Custom Pill Progress Bar
Container(
  height: 8,
  width: double.infinity,
  decoration: BoxDecoration(
    color: Colors.grey.withValues(alpha: 0.1),
    borderRadius: BorderRadius.circular(4),
  ),
  child: FractionallySizedBox(
    alignment: Alignment.centerLeft,
    widthFactor: percentComplete, // e.g. 0.68
    child: Container(
      decoration: BoxDecoration(
        color: accentColor, // e.g. Color(0xFF6366F1)
        borderRadius: BorderRadius.circular(4),
      ),
    ),
  ),
)
```

### Pill Tab Badges
Used for category selections (e.g. Active, Upcoming, Completed).
*   **Unselected**: Transparent background, grey text.
*   **Selected**: Solid primary indigo (`#6366F1` / `#4F46E5`), white text.

---

## 🧭 4. Navigation & Shell

*   **Bottom Navigation**: Frosted glass floating design with a high blur coefficient (`sigmaX: 30`, `sigmaY: 30`).
*   **Active States**: Highlighted using a soft pill background tint (`color.withValues(alpha: 0.1)`) instead of hard lines.
*   **Center Action Button**: Large circular elevated action button in the center (primary indigo/teal) with a soft drop shadow.

---

## ✍️ 5. Typography Hierarchy

*   **Page Title**: `fontSize: 28`, `fontWeight: FontWeight.w900` (heavy block title).
*   **Sub-Header / Card Label**: `fontSize: 16`, `fontWeight: FontWeight.w700`.
*   **Muted Description**: `fontSize: 13`, `fontWeight: FontWeight.w500`, `color: Color(0xFF64748B)`.
*   **Metric / Value Highlight**: `fontSize: 22` or `24`, `fontWeight: FontWeight.w900` (very bold metrics).
