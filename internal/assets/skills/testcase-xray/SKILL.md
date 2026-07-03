---
name: testcase-format-xray
description: >
  Test case format for Jira projects that use the Xray plugin. Uses Xray's "Test"
  issue type, the "Manual" Test Type by default, and the dedicated Test Steps
  custom field instead of a body table.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-test-designer,qa-automator
  phase: formats
---

# Test Case Format — Xray

## When to Use

Use this format when:

- The target Jira project has the **Xray** plugin installed (check for issue type "Test" with a "Test Type" field).
- OR `format-mirror` found an existing TC in the project that exposes Xray custom fields.
- OR `.qason/templates/testcase.md` in the project declares `platform: xray`.

If the project is Jira classic without Xray → use `testcase-format-default` instead. Using Xray fields on a project without Xray will fail the create API call or silently drop fields.

## Xray-Specific Fields (do NOT treat as description body)

These are **custom fields on the issue**, populated via `fields.customfield_xxxxx` at create time, not markdown in the description. The exact `customfield_xxxxx` ID varies per Jira instance — `format-mirror` MUST read one existing Xray TC from the project to discover the real IDs before creating new ones.

| Xray field | Typical API name | Purpose |
|---|---|---|
| Test Type | `customfield_10200` (varies) | `Manual` / `Cucumber` / `Generic` — default to `Manual` |
| Test Steps | `customfield_10210` (varies) | Structured array of step objects (see below) |
| Pre-Conditions | Linked issue (type "Pre-Condition") | Reusable setup, not inline text |
| Tests | Linked issue (`tests` / `tested by`) | Link to the parent story |

**Do not assume the IDs above.** Always call `jira_get_issue` on an existing Xray TC first and grep the customfield values. Mirror what's there.

## Test Steps — Xray Structured Format

Xray expects an **array of step objects**, not a table. Each step has three rich-text fields:

```json
{
  "steps": [
    { "action": "Navigate to /login", "data": "", "result": "Login form is displayed with email and password inputs" },
    { "action": "Enter an expired token in the Authorization header", "data": "Bearer <expired-jwt>", "result": "Response is HTTP 401 with body { \"error\": \"token_expired\" }" },
    { "action": "Observe UI", "data": "", "result": "User is redirected to /login with error toast 'Session expired'" }
  ]
}
```

Some Xray versions accept these fields as ADF (Atlassian Document Format) instead of plain text. `format-mirror` MUST confirm the format from an existing TC.

## Summary and Description

Summary and description follow the same conventions as `testcase-format-default`:

- Summary: `TC-[Component]: [Scenario]`
- Description (reduced, because steps live in Xray's custom field):

```markdown
h2. Objective
{{one_line_goal_of_this_test}}

h2. Traceability
* *Parent Story*: {{parent_issue_key}}
* *Requirement Covered*: {{ac_ids_comma_separated}}
* *Test Type*: Manual
* *Automation Status*: {{Automated|Not automated|Planned}}
* *Script*: {{path|N/A}}

h2. Test Data
* {{fixtures_or_input_data}}

h2. Pre-Conditions
(Pre-Conditions are linked Xray issues, not inline text. See the "Pre-Condition" links on this ticket.)
```

## Rules Specific to Xray

1. **Never** put the steps in the description body — the Test Steps custom field is where Xray renders them and where Test Execution tracks pass/fail per step. Steps in the description are invisible to Xray reports.
2. **Label** `qason-generated` still applies (identifies QASON-authored TCs).
3. **Test Type** defaults to `Manual`. If generating a Cucumber scenario, use `Cucumber` AND put the Gherkin in the Cucumber Scenario field, not in Test Steps.
4. **Preconditions** as separate Xray issues, not body text — if none exist, create or skip; document the decision in the description.
5. If the instance requires **Test Repository** folder assignment, mirror the folder of the existing TC you read; do not invent a new folder.
