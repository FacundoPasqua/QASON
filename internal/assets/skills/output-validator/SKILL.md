---
name: output-validator
description: >
  Validates generated test code by invoking `qason verify` against the files
  just produced. Catches hallucinated imports, syntax errors, lint violations,
  and (at level=full) Playwright discovery failures before handing output
  back to the user.
  Trigger: Immediately after any skill that writes .ts/.tsx/.js/.jsx files
  (e.g. playwright-scaffold, unit-test-gen, api-test-gen, e2e-test-gen).
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: development
---

# Output Validator

## Why this skill exists

The QA Transformation Plan (Section 7 — Riesgos y Mitigaciones) explicitly
calls out two risks this skill mitigates:

1. **"Tests generados con falsos positivos/negativos"** → mitigation requires
   *"métricas de calidad de tests generados"*.
2. **"Hallucinations en generación de tests"** → mitigation requires
   *"Review process, validación contra specs, ejecución en ambientes
   controlados"*.

This skill is the **first line of defense** against hallucinations: it runs
the actual TypeScript/JavaScript toolchain (tsc, eslint, optionally
playwright) and reports any failure before the code leaves the agent's hands.

## When to Use

- Immediately after writing any test file (unit, integration, e2e, API)
- After scaffolding a framework (Playwright, Cypress, Postman)
- Before claiming "tests generated" to the user
- In the final step of workflows like `spec-to-test`, `pr-guardian`

## Critical Rules

1. **NEVER** claim generated code is ready without running this skill first
2. Run with `--level medium` by default (syntax + lint)
3. Run with `--level full` when the target project has `@playwright/test`
   installed (dry-runs the spec discovery)
4. If the report has failures, **DO NOT DELIVER**: fix the issues and re-run
5. Report tool-missing status as informational — don't treat it as a failure
6. Pass the EXACT file paths that were just written (not the whole project)

## Workflow

1. **Collect** the list of files written in the current task
2. **Run** `qason verify --level medium <file-paths>`
3. **Read** the report:
   - `Result: OK` → proceed to delivery
   - `Result: FAIL` → parse `FAILURES` section, fix, goto step 2
   - `TOOLS MISSING` → log warning, continue (not a blocker)
4. **Report back** with the summary line from the report:
   `"Validation: X files checked, Y passed, Z failed (level=medium)"`

## Invocation

Use the Bash tool:

```bash
qason verify --level medium path/to/file.spec.ts path/to/other.ts
```

For JSON output (useful when you want to parse):

```bash
qason verify --level medium --json path/to/file.spec.ts
```

For the deepest check on Playwright specs (requires project to have
`@playwright/test` installed):

```bash
qason verify --level full path/to/file.spec.ts
```

## Handling Failures

When `qason verify` exits non-zero:

1. **Read each failure's file, line, column, rule, and message**
2. **Categorize** the failure type:
   - Syntax error (TS1xxx, Node SyntaxError) → rewrite the offending construct
   - Type error (TS2xxx) → fix imports, types, or assertions
   - Lint error (eslint rule) → adjust to match project conventions
   - Playwright discovery error → usually a bad `test(...)` signature or
     import path
3. **Re-read** the generated file and apply the fix in place
4. **Re-run** `qason verify` on the same paths
5. **Give up after 3 attempts** — escalate to the user with the remaining
   errors rather than looping indefinitely

## Output Template

Always include this section at the end of your response when this skill
fired:

```
Validation Report (qason verify --level medium)
-----------------------------------------------
Files checked: <N>
Pass: <N> | Fail: <N> | Tool-missing: <N>
<if failures: list them with file:line and rule>
<if tool-missing: "Note: install <tool> for deeper coverage">
```

## Why not just "trust the LLM"

The PDF is explicit: hallucinations in test generation are a **medium-
probability, medium-impact risk**. A 2-minute verification step catches
~80% of those failures cheaply — without this, the 401 lines of Playwright
we generated for EAAAAC-206 could have contained bad imports and no one
would have noticed until a human ran them.

The rule is simple: **we don't ship tests we haven't verified.**
