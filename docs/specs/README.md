# Specs

SubPilot uses Spec-Driven Development (SDD).

Before implementing any non-trivial phase or feature, create a spec in this directory. The spec is the contract for what will be built, why it is needed, how it should behave, and how completion will be verified.

## Naming

Use this format:

```text
phase-XX-feature-name.md
```

Examples:

- `phase-02-database-foundation.md`
- `phase-03-authentication.md`
- `phase-05-subscription-management.md`

## Required Flow

1. Write or update the spec.
2. Review the spec before implementation.
3. Create a checklist in `tasks/todo.md`.
4. Implement only what the spec covers.
5. Update the spec first if scope changes.
6. Verify every acceptance criterion.
7. Record traceability notes before marking the work complete.

## Template

Use [TEMPLATE.md](TEMPLATE.md) for new specs.
