# Horizon Flutter — Phase 2 Implementation Roadmap

**10 milestones** completing all remaining gaps. No Redis, no Kubernetes, no Grafana, no testing.

---

## How to Use

Execute milestones in order. Each is self-contained with context reminders, backend audit, implementation plan, verification checklist, and architect message.

---

## Milestone Overview

| # | File | Focus | Area |
|---|------|-------|------|
| 1 | `01-docker-infrastructure.md` | Add monolith to docker-compose, auto-migrate on startup | Infra |
| 2 | `02-auth-profile-endpoints.md` | Change Password, User Profile/Preferences/Privacy HTTP routes | Backend |
| 3 | `03-transaction-edit-delete.md` | PUT update + DELETE for transactions | Backend |
| 4 | `04-planning-endpoints.md` | Budget, Retirement, Emergency Fund, Debt Payoff, Investment | Backend |
| 5 | `05-data-management.md` | Export, Import, Backup, Restore | Backend + Flutter |
| 6 | `06-ai-streaming.md` | SSE streaming for AI chat | Backend |
| 7 | `07-notifications-complete.md` | Push registration, delete notification | Backend |
| 8 | `08-flutter-polish.md` | Animations, haptic, responsive layouts | Flutter |
| 9 | `09-flutter-branding.md` | App icon, splash, illustrations, empty states | Flutter |
| 10 | `10-final-cleanup.md` | Remove unused packages, lint, final verification | Both |

---

## Key Constraints

- **No Redis** — removed from codebase already
- **No Kubernetes** — not in scope
- **No Grafana/Observability** — not in scope
- **No testing** — testing strategy is deferred
- **Monorepo** — already organized as one
- **Architecture is frozen** — follow existing patterns
