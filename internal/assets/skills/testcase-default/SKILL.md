---
name: testcase-format-default
description: >
  Default, platform-agnostic test case format for QASON. Used by jira-testcase-linker,
  ado-testcase-linker, and any sub-agent that produces a TC without a project-specific
  template. Optimized for Jira (classic) and ADO Test Plans without plugins.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-test-designer,qa-automator
  phase: formats
---

# Test Case Format — Default

## When to Use

Use this format when ALL of the following are true:

- The target system is Jira (classic, no Xray plugin) or Azure DevOps Test Plans.
- The project does NOT ship a `.qason/templates/testcase.md` override.
- `format-mirror` did not find an existing TC in the project to mirror.

If any of those assumptions fail, prefer the override OR `testcase-format-xray` OR the existing-TC mirror — not this one.

## Required Fields

Every TC created with this format MUST include these fields. Missing fields are a bug, not a stylistic choice.

| Field | Where it lives | Why it's required |
|---|---|---|
| Summary | Issue title | Searchable, must start with `TC-[Component]:` |
| Parent traceability | Linked issue | Orphan TCs are banned |
| Preconditions | Body | Reproducibility |
| Numbered steps (action + expected) | Body | Executable by any tester |
| Test type label | Labels | `functional`, `edge-case`, `negative`, `regression` |
| Automation status | Body | So coverage reports know what's manual vs automated |
| `qason-generated` label | Labels | Identifies our output, enables bulk ops |

## Template

Copy the block below, replace `{{placeholders}}`, send as the description field.

```markdown
h2. Test Case: {{scenario_name}}

h2. Traceability
* *Parent Story*: {{parent_issue_key}}
* *Requirement Covered*: {{ac_ids_comma_separated}}
* *Test Type*: {{functional|edge|negative|regression}}

h2. Preconditions
* {{state_or_data_required}}
* {{user_role_or_permissions}}

h2. Test Steps
||Step||Action||Expected Result||
|1|{{action_1}}|{{expected_1}}|
|2|{{action_2}}|{{expected_2}}|
|3|{{action_3}}|{{expected_3}}|

h2. Test Data
* {{specific_input_values_or_fixtures}}

h2. Automation Status
* *Automated*: {{Yes|No|Planned}}
* *Framework*: {{Playwright|Cypress|Jest|...|N/A}}
* *Script Path*: {{relative/path/to/test.spec.ts|N/A}}
```

## Naming Convention

- Summary MUST follow: `TC-[Component]: [Scenario in present tense, <80 chars]`
- Good: `TC-Auth: User with expired token gets 401 on protected route`
- Bad: `Testing login` / `Login test case 5` / `TC: auth works`

## Rules

1. **One purpose per TC** — if the steps verify two distinct behaviors, split into two TCs.
2. **Steps are numbered and specific** — "click the button" is wrong; "click the 'Save' button in the footer" is right.
3. **Expected results are observable** — no "the system works correctly"; write what a human or script can check.
4. **No environment-coupled assumptions** — if the TC assumes staging data, call it out in Preconditions.
5. **Repeatability** — running the TC twice must produce the same result or the TC must declare its idempotency assumption.
