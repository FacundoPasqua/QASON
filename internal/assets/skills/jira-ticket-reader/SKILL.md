---
name: jira-ticket-reader
description: >
  Reads Jira tickets, stories, and epics via MCP tools to extract testable
  requirements, acceptance criteria, and linked context for QA analysis.
  Trigger: When a Jira ticket key is provided or user asks to analyze a Jira issue.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-analyst
  phase: planning
---

# Jira Ticket Reader

## When to Use

- User provides a Jira ticket key (e.g., PROJ-123) for QA analysis
- User asks to "read this ticket", "analyze this story", or "get requirements from Jira"
- Before generating a test plan and you need the source requirements
- When triaging a backlog of tickets for testability assessment

## Critical Rules

1. ALWAYS use MCP Jira tools — never ask the user to copy-paste ticket content
2. Fetch the FULL ticket: summary, description, acceptance criteria, comments, attachments, and linked issues
3. Parse acceptance criteria into discrete, testable statements — one per line
4. Flag any ticket that has NO acceptance criteria as **INCOMPLETE**
5. Follow linked issues (blocks, is-blocked-by, relates-to) to understand the full scope — up to 2 levels deep
6. Extract custom fields that commonly hold QA-relevant data (story points, sprint, components, fix version)
7. NEVER skip the comments — they often contain clarifications that contradict or extend the description
8. If the ticket is an Epic, fetch its child stories and summarize scope across all children
9. Preserve the original Jira ticket key as a traceability reference in all outputs

## Workflow

1. **Fetch the ticket** using `jira_get_issue` with the provided ticket key
   - Include fields: summary, description, acceptance criteria, status, priority, labels, components, fix version
   - Request comment history and attachment metadata
2. **Parse the description** for testable requirements
   - Look for Given/When/Then patterns, bullet lists, numbered steps
   - Extract any embedded tables, diagrams references, or linked Figma/design URLs
3. **Extract acceptance criteria** as discrete testable statements
   - If stored in a custom field, parse that field specifically
   - If embedded in description, extract and separate them
   - If missing entirely, flag as INCOMPLETE and generate suggested criteria from the description
4. **Fetch linked issues** using `jira_get_issue` for each linked ticket
   - Categorize links: blocks, is-blocked-by, relates-to, duplicates, is-child-of
   - Summarize linked ticket context (title + status) for dependency mapping
5. **Check comments** for requirement clarifications, scope changes, or PO decisions
   - Extract any comment that modifies or extends the original requirements
   - Note the author and date of clarifying comments
6. **Assess completeness** — identify gaps, ambiguities, and missing information
7. **Produce the structured output** ready for consumption by test-plan-gen or risk-matrix-gen

## Handling Epics and Multi-Ticket Scope

When the provided ticket is an Epic:

1. Use `jira_search` with JQL: `"Epic Link" = EPIC-KEY OR parent = EPIC-KEY`
2. Fetch each child story's summary, status, and acceptance criteria
3. Group children by status (To Do, In Progress, Done)
4. Produce a consolidated requirements summary across all children

## Output Template

```markdown
## Requirements Summary: [TICKET-KEY] — [Summary]

**Source**: Jira ([TICKET-KEY])
**Type**: [Story | Bug | Epic | Task]
**Status**: [Current status]
**Priority**: [Priority]
**Sprint**: [Sprint name if assigned]
**Components**: [Component list]
**Fix Version**: [Version if set]

### Acceptance Criteria
- [ ] AC-001: [Testable acceptance criterion]
- [ ] AC-002: [Testable acceptance criterion]
- [ ] AC-003: ...

### Requirements Extracted from Description
- [ ] REQ-001: [Explicit requirement from description]
- [ ] REQ-002: [Implicit requirement inferred from context]

### Clarifications from Comments
| Date | Author | Clarification |
|------|--------|---------------|
| [date] | [author] | [What was clarified or changed] |

### Linked Issues
| Key | Type | Relationship | Status | Summary |
|-----|------|-------------|--------|---------|
| [KEY] | [Story/Bug] | [blocks/relates-to] | [Status] | [Title] |

### Completeness Assessment
- **Acceptance Criteria**: [Complete | Partial | Missing]
- **Description Quality**: [Detailed | Adequate | Vague]
- **Design References**: [Linked | Not provided]
- **Dependencies Identified**: [Yes | No]

### Gaps & Ambiguities
1. [Specific gap or ambiguity — why it matters for testing]
2. ...

### Suggested Questions for PO
1. **[Question]** — Impact: [How the answer affects test scope]
2. ...
```
