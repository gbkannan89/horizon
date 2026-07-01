# Horizon Flutter — Phase 1 Implementation Roadmap

**9 milestones to a production-ready MVP Flutter application.**

---

## How to Use

Each `.md` file is a self-contained implementation prompt. Process them in order, one at a time:

1. Open `XX-milestone-name.md`
2. Follow the Standard Operating Procedure:
   - **Backend Audit** — verify APIs before writing code
   - **Implementation Plan** — list files before coding
   - **Implementation** — build with no placeholders, no TODOs
   - **Build Verification** — `flutter analyze` must pass with zero errors
   - **Functional Verification** — check the checklist
   - **Final Report** — document everything
   - **Message to Solution Architect** — highlight concerns
3. When the milestone is complete and verified, move to the next one.

---

## Milestone Overview

| # | File | Feature | Est. Sessions | Depends On |
|---|---|---|---|---|
| 2 | `02-transactions-experience.md` | Transactions Module (list, detail, CRUD, search, filters, pagination) | 1 | M1 Foundation |
| 3 | `03-global-search-and-shared-widgets.md` | Global Search + Shared Widget Library + Theme Constants | 1-2 | M2 |
| 4 | `04-ai-advisor.md` | AI Advisor Chat (markdown, prompts, context, providers) | 2 | M3 |
| 5 | `05-insights.md` | Insights (spending, income, trends, fl_chart, opportunities, warnings) | 1-2 | M4 |
| 6 | `06-financial-planning.md` | Financial Planning (projections, scenarios, "coming soon" for missing) | 1-2 | M5 |
| 7 | `07-notifications.md` | Notifications Center (tabs, swipe, preferences, unread badge) | 1 | M6 |
| 8 | `08-settings-and-profile-completion.md` | Settings Completion (AI settings, data management, fix no-ops) | 1 | M7 |
| 9 | `09-ui-consistency.md` | UI Consistency (full shared widget adoption, design system compliance) | 2 | M8 |
| 10 | `10-production-readiness.md` | Production Readiness (audit, performance, security, final smoke test) | 1-2 | M9 |

---

## Design System Reference

All prompts reference these key design documents. Read them before starting:

- `horizon-spec/products/DESIGN_SYSTEM.md` (501 lines) — Colors, typography, spacing, cards, states, accessibility
- `horizon-spec/products/UX_PHILOSOPHY.md` (487 lines) — UX principles, decision-first design, progressive disclosure
- `horizon-spec/products/NAVIGATION_SYSTEM.md` (525 lines) — Navigation architecture, search philosophy, deep linking
- `horizon-spec/products/PROJECT_CONTEXT.md` (165 lines) — Product vision, AI principles
- `horizon-spec/engineering/FRONTEND_ENGINEERING_GUIDE.md` (517 lines) — Engineering principles, state management, offline behavior

---

## Architecture Constraints

These apply to EVERY milestone:

- **Architecture is frozen** — Do not redesign. Follow existing patterns.
- **Specifications are frozen** — Do not modify specs. Implement as specified.
- **Technology stack is frozen** — Do not add new frameworks (Riverpod, GoRouter, Dio, Flutter Secure Storage are the stack).
- **No assumptions** — Verify the backend before writing Flutter code.
- **No placeholders** — Either implement fully or report why not.
- **No TODOs** — If a backend dependency is missing, stop and report it.
- **Build verification required** — `flutter analyze` with zero errors before marking complete.

---

## Current State (Before Phase 1)

- **Flutter app** at `C:\Kannan\Horizon\horizon\apps\mobile\`
- **Backend** at `C:\Kannan\Horizon\horizon\services\` and `C:\Kannan\Horizon\horizon\cmd\server\main.go` (monolith) + `services\infra\api-gateway\` (standalone gateway)
- **Design specs** at `C:\Kannan\Horizon\horizon-spec\`
- **Milestone 1 (Foundation)** was completed before these prompts — Auth, Settings, Networking, Error handling are stable
