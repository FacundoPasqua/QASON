---
name: validation-memory
description: >
  Learns from recurring output-validator failures. Before generating tests,
  checks memory for past validation failures in this project (e.g. "the LLM
  keeps inventing the wrong Playwright import path here") so the same mistake
  isn't re-committed. After a validation cycle, saves the failure fingerprint.
  Trigger: Before running output-validator (read mode) and after it reports
  failures that were fixed (save mode).
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: development
  requires_mcp: engram
---

# Validation Memory

## Why this skill exists

`output-validator` catches hallucinations *reactively* — it tells the agent
"you wrote a bad import" after the fact. `validation-memory` closes the loop:
it remembers **which hallucinations recur in which project** so the next
generation avoids them *proactively*.

Pattern matches the other memory skills (see `bug-pattern-memory`,
`flaky-test-memory`): Engram is preferred and the QASON file fallback is mandatory when Engram tools are absent.

## When to Use

- **Before** generating a test file in a known project (read mode) — preload
  the "things the LLM gets wrong here" list into the current task prompt.
- **After** `output-validator` reports failures and the agent fixes them
  (save mode) — persist the fingerprint so the next run learns.

## Critical Rules

1. If engram MCP is not available → skip silently, never mention to user
2. Search BEFORE generating — applying a remembered fix up front is
   ~10x cheaper than running verify → fail → retry
3. Save ONLY confirmed fixes (failure + the edit that resolved it) — do NOT
   save raw failures, they're noise without the resolution
4. Key by **project + file-pattern + rule-id** — not by individual file path,
   otherwise memory can't generalize

## Topic Key

`qason/project/{name}/validation-failures`

## Workflow (Read Mode)

1. **At the start** of a test-generation task:
   - `mem_search(query: "qason/project/{name}/validation-failures {file-pattern}", project: "{name}")`
2. **Inject findings** into the working prompt:
   > Known validation pitfalls in this project:
   > - `@playwright/test` must come from `@playwright/test` not `playwright/test` (saved after 3 recurrences)
   > - `test.describe.configure({ mode: 'parallel' })` requires top-level placement
3. **Generate** avoiding the listed pitfalls.

## Workflow (Save Mode)

After `output-validator` reports `FAIL` and the agent fixes it:

1. **Extract the fingerprint**:
   ```
   Rule/error: <TS2345 | no-unused-vars | Node SyntaxError | ...>
   File pattern: <glob that matched, e.g. tests/**/*.spec.ts>
   Before fix: <the offending snippet>
   After fix: <the corrected snippet>
   Why: <one-line explanation>
   ```
2. **Merge with existing memory**:
   - `mem_search` the topic
   - Append the new fingerprint (dedupe by `rule + before`)
   - `mem_save(topic: "qason/project/{name}/validation-failures", content: <combined>, project: "{name}")`
3. **Increment recurrence counter** if same fingerprint seen ≥3 times:
   promote it to the "must-check" section at the top of the memory entry.

## Output Template (Read Mode)

Include this at the top of your internal reasoning (NOT shown to user
unless they ask):

```
Validation-memory preload (N hits):
- [TS2345] tests/**/*.spec.ts — wrong Playwright import (3 recurrences)
- [no-unused-vars] tests/fixtures/*.ts — unused imports from copy-paste (1)
```

## Fallback: no engram available

Follow the `memory-fallback-protocol` skill:

- Validation fingerprints live under
  `~/.qason/memory/qason/project/{name}/validation-fingerprints.md` in the
  file backend (same schema as above, appended section by section).
- Read → search for the current rule+file-glob in the file body.
- Write → read the file, merge the new fingerprint (or bump its counter if
  it already exists), rewrite with an updated version number.

Only if BOTH backends are unreachable, skip the search silently and proceed
with the normal generation + output-validator cycle. ADR-009 originally
called memory OPTIONAL; under the 3-layer design it is REQUIRED at
`acceleration`+ presets but always downgradable to the file backend.

## Why not store the raw `qason verify` JSON?

Because the JSON report is volatile — file paths, line numbers, timestamps.
Fingerprints abstract away the noise and keep only the reusable insight:
*this rule, on this kind of file, needs this fix.*
