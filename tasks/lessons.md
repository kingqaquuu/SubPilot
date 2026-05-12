# Lessons

## Documentation Scope

- When `AGENTS.md` already contains project principles, stack rules, and workflow rules, do not repeat them in generated project planning documents.
- For a requested development specification with phases, focus on the concrete Phase Plan: goals, prerequisites, deliverables, task breakdown, dependencies, acceptance criteria, and exit gates.

## Port Defaults

- Avoid using `8080` as the default backend port because it commonly conflicts with other local services. Prefer `18080` unless the user specifies otherwise.

## Commit Messages

- Commit messages must include key details whenever possible: what changed, why it changed, and how it was verified.
- When existing commit messages are too terse, supplement or rewrite them instead of leaving unclear history.
