# Horizon AI Agent Instructions

This file provides context for AI coding agents working on this repository.

## Repository Structure

| Path | Purpose |
|---|---|
| `packages/` | Shared Go libraries (errors, events, types, etc.) |
| `services/` | Backend services (domains, engines, infra, import) |
| `apps/` | Frontend applications (Flutter) |
| `proto/` | Protobuf definitions |
| `deploy/` | Docker and deployment configurations |

## Key Documents (read before implementing)

- `MASTER_BUILD_GUIDE.md` — Implementation discipline
- `MASTER_BUILD_PLAN.md` — Build sequencing
- `AI_IMPLEMENTATION_PLAYBOOK.md` — Coding standards
- `TECH_STACK.md` — Approved technologies

## Rules

- Never invent business rules
- Never redesign architecture
- Never modify specifications
- Every field must trace to a specification
- Every business rule must be in the domain layer
- Every event must be published on state change
