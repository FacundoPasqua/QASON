---
name: format-mirror
description: >
  Mandatory protocol for QASON sub-agents before creating test cases or bug reports
  in any tracker (Jira, ADO). Forces the agent to READ an existing artifact from the
  project and MIRROR its format, instead of applying a hardcoded template blindly.
  This prevents the "QASON format vs project format" mismatch that requires manual
  cleanup after every run.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator,qa-ops,qa-test-designer
  phase: formats
---

# Format Mirror — Read Before You Create

## Why This Skill Exists

Every tracker has its own conventions: custom fields, labels, summary prefixes, step structures, severity taxonomies. A hardcoded template will match those conventions for one project and miss for the next. Users end up fixing formatting by hand after every QASON run — the very cost QASON was meant to eliminate.

The rule is simple and non-negotiable: **read an existing artifact from the project first, mirror its format, then create yours.**

## Lookup Order — Priority from Highest to Lowest

When about to create a TC or bug, the sub-agent MUST resolve the format by walking this list top-to-bottom and stopping at the first hit:

1. **Project-local override** — `.qason/templates/testcase.md` or `.qason/templates/bug.md` in the current git repo root.
2. **User-global override** — `~/.qason/templates/testcase.md` or `~/.qason/templates/bug.md`.
3. **Mirror from an existing artifact** — live TC/bug from the target tracker (see protocol below).
4. **Embedded QASON default** — `testcase-format-default`, `testcase-format-xray`, `bug-format-default`, or `bug-format-incident-style`, selected by the tracker/platform.

If `format-mirror` is NOT active (Foundation preset), the agent skips steps 3 and uses steps 1–2 and 4 only.

## Protocol for Step 3 — Mirror from Existing

### Test Cases

Before calling `jira_create_issue` or `wit_create_work_item` for a new TC:

1. **Search for a reference TC**:
   - Jira: `jira_search` with JQL `project = "{PROJECT}" AND issuetype = "Test" ORDER BY updated DESC` — take the first result.
   - ADO: `wit_query_by_wiql` with `SELECT [System.Id] FROM workitems WHERE [System.WorkItemType] = 'Test Case' AND [System.TeamProject] = '{PROJECT}' ORDER BY [System.ChangedDate] DESC`.
   - If the project has multiple test types or components, prefer a TC linked to the same parent story or the same component as the one you're about to create.
2. **Fetch the full TC**:
   - Jira: `jira_get_issue(key, expand="renderedFields,names,schema")` — you need `names` and `schema` to discover custom field IDs.
   - ADO: `wit_get_work_item(id, expand="Fields")`.
3. **Inspect the structure**:
   - What custom fields are populated? (Xray's `customfield_10210`, priority custom field, test repository folder, etc.)
   - What labels / tags? What Components?
   - Is the description in ADF, markdown, or plain text?
   - How are steps structured? Xray custom field array? Body table? Numbered list?
4. **Discover the field IDs** — never hardcode `customfield_10200`. Read what the actual instance uses.
5. **Mirror**:
   - Replicate the field structure exactly.
   - Replicate the labels (plus add `qason-generated`).
   - Replicate the summary convention (prefix, casing, punctuation).
   - Keep your content in the same format (ADF vs markdown vs plain).
6. **Do NOT mirror**:
   - The **content** (the actual test being described) — that comes from the test plan.
   - The **author/assignee** — leave blank or use the current user, whatever the tracker expects.
   - **Dates or version fixVersion** that are specific to the reference TC.

### Bug Reports

Before calling `jira_create_issue` or `wit_create_work_item` for a bug:

1. **Search recently-closed bugs** (last 90 days) in the same project:
   - Jira: `jira_search` with JQL `project = "{PROJECT}" AND issuetype = Bug AND status in (Closed, Done, Resolved) AND resolved >= -90d ORDER BY resolved DESC`.
   - ADO: `wit_query_by_wiql` with an analogous filter on `[System.State]` and `[Microsoft.VSTS.Common.ResolvedDate]`.
2. **Fetch 2–3 bugs** to spot conventions (one sample is biased).
3. **Inspect**:
   - Severity field name and valid values.
   - Priority field name and valid values.
   - Summary prefix convention (`[SEV1][Component]`? `[BUG]`? `#{id}:`?).
   - Required labels / components / environment field.
   - Description structure (postmortem-style? simple repro? ADF with panels?).
4. **Decide between default and incident style**:
   - If the mirrored bugs have SEV levels and postmortem links → use `bug-format-incident-style`.
   - Otherwise → use `bug-format-default`.
5. **Mirror** summary prefix, label set, severity/priority vocabulary, description structure.

## Critical Rules

1. **Never skip the mirror step** in Acceleration or Optimization presets. If the tracker MCP is unavailable, fail the creation and tell the user — do NOT fall back to the default silently. Silent fallback is how we got here.
2. **Do not mirror content, only structure.** A reference TC about `login` has no bearing on the TC you're creating about `checkout`.
3. **Custom field IDs are instance-specific.** Re-discover them per project, cache nothing.
4. **Log what you mirrored** — in the sub-agent's output, include a 1-line note: `Mirrored format from ABC-123 (Xray, customfield_10210=steps, labels: qa,automation)`. This gives the user a trail for debugging.
5. **When in doubt, ask** — if the mirror results are contradictory (two TCs with different structures in the same project), surface the ambiguity. Do not pick one silently.

## Interaction with Override Files

If `.qason/templates/testcase.md` exists in the project, that file **wins over mirror**. The rationale: a team that went to the trouble of checking in a template has already decided the format, and mirror may surprise them by following a legacy TC.

The override file is plain markdown with `{{placeholders}}` matching the templates shipped in `testcase-format-*` and `bug-format-*`. No YAML frontmatter required for override files — just the template body.

## Output Requirement

After every creation, the sub-agent MUST append a 1-line provenance note to its output:

```
Format source: {project-override | mirrored-from-ABC-123 | default-template-{name}}
```

This makes the source auditable without reading the agent's internal reasoning.
