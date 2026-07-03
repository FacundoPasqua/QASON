---
name: project-context-saver
description: >
  Saves QA-specific project discoveries to persistent memory with standardized topic keys.
  Called automatically when the agent learns something non-obvious about the project:
  stack details, conventions, business rules, auth flows, or patterns.
  Trigger: After discovering any non-obvious project fact during a QA task.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-analyst, qa-automator, qa-test-designer, qa-reviewer, qa-ops
  phase: post-discovery
  requires_mcp: engram
---

# Project Context Saver

## When to Use

- After discovering the project uses a specific framework pattern (POM, AAA, container-presentational)
- After identifying a non-obvious business rule from the code or docs
- After finding an authentication mechanism (JWT with refresh, session-based, OAuth)
- After learning a team convention (test file naming, fixture location, data-testid strategy)
- After identifying a hot spot (file with recurring bugs)
- After discovering existing test data (fixtures, seed users, mock endpoints)

## Critical Rules

1. NEVER overwrite a topic key with less information than what is already there — upsert, merge, preserve
2. Use the EXACT topic key format: `qason/project/{name}/{category}` — do NOT improvise
3. Include the `project: "{name}"` parameter in every `mem_save` call so scoping works
4. Save in structured format (What/Why/Where/Learned) — future agents need to parse it
5. If engram MCP is NOT available, skip silently — do NOT block the task

## Topic Key Categories

| Category | What goes here |
|----------|----------------|
| `stack` | "Playwright, TypeScript, pnpm, Vitest for unit" |
| `conventions` | "Tests live next to source (*.test.ts), use describe/it, beforeEach for setup" |
| `domain` | "Users have roles: admin/editor/viewer. Orders have 3 states: pending/confirmed/shipped" |
| `auth-flow` | "JWT in httpOnly cookie, 15min access + 7d refresh. Login at /api/auth/login" |
| `bug-patterns` | "Null pointer in checkout flow when user has no shipping address" |
| `flaky-tests` | "tests/e2e/payment.spec.ts is flaky — race condition with Stripe webhook" |
| `test-data` | "Test user: qa@test.com / Test123!. Fixture in tests/fixtures/users.json" |
| `coverage-map` | "High: auth, user CRUD. Low: reports module, notifications" |

## Workflow

1. **Detect the project name** (same detection as project-context-loader)
2. **Categorize** the discovery: which topic key bucket does it belong to?
3. **Check existing content**:
   - `mem_search(query: "qason/project/{name}/{category}")`
   - `mem_get_observation(id)` to see current state
4. **Merge**, do NOT overwrite:
   - If topic exists, add the new fact to the existing entry
   - Preserve older entries that are still valid
5. **Save** with the structured format:
   ```
   mem_save(
     title: "Project context: {name}/{category}",
     type: "discovery",
     project: "{name}",
     topic_key: "qason/project/{name}/{category}",
     content: "
       **What**: [the discovered fact]
       **Why**: [how this affects QA work]
       **Where**: [files/paths where this is visible]
       **Learned**: [gotchas, edge cases]
     "
   )
   ```

## Output

Do NOT produce user-facing output for saves — they happen silently as a side effect of doing QA work. Just confirm to yourself that the save succeeded and continue.

## Fallback: no engram available

Follow the `memory-fallback-protocol` skill to write the discovery as a
markdown file at `~/.qason/memory/qason/project/{name}/{category}.md`. The
save still happens silently — the only difference is the backend.

The merge step above still applies: Read the existing file (if any), combine
its body with your new content using the same **What/Why/Where/Learned**
template, bump the version counter, Write the whole file.

Only if BOTH backends are unreachable do we drop the save silently. At
`acceleration` and above, that state is a broken install — the install gate
should have caught it.
