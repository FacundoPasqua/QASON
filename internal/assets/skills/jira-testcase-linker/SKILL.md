---
name: jira-testcase-linker
description: >
  Creates test case issues in Jira and links them to parent stories or features,
  establishing full traceability between requirements and test coverage.
  Trigger: When test cases need to be documented in Jira and linked to stories.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-test-designer
  phase: testing
---

# Jira Test Case Linker

## When to Use

- Test cases have been designed and need to be recorded in Jira
- User asks to "link test cases to the story", "create test cases in Jira", or "add traceability"
- After test-plan-gen or e2e-test-gen produces test cases that need formal tracking
- Sprint planning requires visibility of test coverage per story

## Critical Rules

1. ALWAYS create test cases as Jira issues (type: Test or Sub-task) — not just comments
2. If the project does not have a "Test" issue type, fall back to Sub-task under the parent story
3. Every test case MUST link back to at least one parent story via issue link
4. Include preconditions, step-by-step actions, and expected results in a structured format
5. Add the label `qason-generated` to all created test cases
6. Use a consistent naming convention: `TC-[Feature]: [Scenario description]`
7. Include traceability metadata — which requirement (REQ-ID or AC-ID) each test covers
8. NEVER create orphan test cases — every test must trace to a requirement
9. When creating multiple test cases for one story, batch them and report a summary
10. Check for existing test cases on the story before creating duplicates
11. Apply the `format-mirror` skill BEFORE creating or updating any test case; fall back to the default format skill only if mirroring finds nothing to mirror

## Workflow

1. **Verify the parent ticket** using `jira_get_issue`:
   - Confirm the story exists and note its key, summary, and acceptance criteria
   - Check for any existing linked test cases to avoid duplication
2. **Check project capabilities**:
   - Determine if the project supports "Test" issue type
   - If not, plan to create Sub-tasks or add structured comments
3. **For each test case**, create an issue using `jira_create_issue`:
   - Type: Test (preferred) or Sub-task (fallback)
   - Summary: `TC-[Feature]: [Scenario]`
   - Description: Structured test case (see template below)
   - Labels: `qason-generated`, test category (e.g., `functional`, `edge-case`, `negative`)
   - Priority: Match the parent story priority unless risk dictates otherwise
4. **Link to parent story** using `jira_create_issue_link`:
   - Link type: "tests" or "is tested by" (check available link types first with `jira_get_link_types`)
   - If "tests" link type is unavailable, use "relates to"
5. **Add traceability comment** on the parent story using `jira_add_comment`:
   - List all created test case keys with their scenario summaries
   - Include a coverage matrix showing which AC each test covers
6. **Report results** — summarize all created test cases and their links

## Test Case Description Format

**Do NOT hardcode a template here.** Format selection is delegated to the `format-mirror` skill, which resolves in this order:

1. Project override — `.qason/templates/testcase.md`
2. User override — `~/.qason/templates/testcase.md`
3. Mirror from an existing TC in the same Jira project (mandatory in Acceleration+)
4. Fallback to `testcase-format-default` (plain Jira) or `testcase-format-xray` (if the project uses Xray, detected from the mirrored TC)

Always call `format-mirror` BEFORE drafting the description. If the mirror target exposes Xray custom fields (`customfield_*` with structured steps), send the steps via those fields — not as a table in the description body.

Provenance requirement: after creating each TC, append a one-line note to your output indicating which source the format came from (`Format source: mirrored-from-ABC-123` / `project-override` / `default-template-xray`).

## Batch Creation Summary

After creating multiple test cases for a story, add this comment to the parent:

```markdown
h2. QASON Test Coverage Summary

*Story*: [STORY-KEY]
*Test Cases Created*: [count]
*Date*: [today]

||Test Case||Scenario||Covers||Type||
|[TC-KEY]|[Brief scenario]|AC-001|Functional|
|[TC-KEY]|[Brief scenario]|AC-001, AC-002|Edge Case|
|[TC-KEY]|[Brief scenario]|AC-003|Negative|

*Coverage*: [X of Y] acceptance criteria covered
*Gaps*: [Any AC not covered and why]
```

## Output Template

```markdown
## Test Cases Created

**Parent Story**: [STORY-KEY] — [Story summary]
**Test Cases Created**: [count]

| Test Case Key | Scenario | Covers | Type |
|---------------|----------|--------|------|
| [TC-KEY] | [Description] | [AC/REQ IDs] | [Type] |
| [TC-KEY] | [Description] | [AC/REQ IDs] | [Type] |

### Coverage Matrix
| Acceptance Criteria | Test Cases | Status |
|---------------------|------------|--------|
| AC-001: [criteria] | [TC-KEY, TC-KEY] | Covered |
| AC-002: [criteria] | — | GAP |

### Gaps Identified
- [Any acceptance criteria not covered and recommended action]
```
