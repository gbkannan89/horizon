# Contributing to Horizon

## Development Workflow

1. Branch from `develop`
2. Implement your changes
3. Run `make test` and `make lint`
4. Submit a pull request to `develop`

## Branch Naming

- `feat/PKG-NNN-description` — New features per package
- `fix/description` — Bug fixes
- `chore/description` — Maintenance

## Commit Messages

```
PKG-NNN: Short description

Longer description explaining what and why, not how.

References: HZN-DOM-001 §4, AC-03
```

## Code Review

All PRs require:
- CI passing
- At least one approval from a code owner
- Architecture compliance verified
- No AI Boundary violations
