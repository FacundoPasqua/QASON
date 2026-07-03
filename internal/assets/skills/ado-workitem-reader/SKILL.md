---
name: ado-workitem-reader
description: >
  Reads Azure DevOps work items (User Stories, Bugs, Features, Tasks) to extract
  testable requirements, acceptance criteria, and linked context for QA analysis.
  Trigger: When an Azure DevOps work item ID is provided or user asks to analyze ADO items.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-analyst
  phase: planning
---

# Azure DevOps Work Item Reader

## When to Use

- User provides an Azure DevOps work item ID (e.g., #12345) for QA analysis
- User asks to "read this work item", "analyze this user story from ADO", or "get requirements from Azure DevOps"
- Before generating a test plan and the source requirements live in Azure DevOps
- When triaging a backlog of ADO work items for testability assessment

## Critical Rules

1. Use Azure DevOps REST API or available MCP tools — never ask the user to copy-paste work item content
2. Fetch the FULL work item: title, description, acceptance criteria, repro steps (for bugs), related work items
3. Handle ALL work item types: User Story, Bug, Feature, Epic, Task — each has different relevant fields
4. Parse acceptance criteria into discrete, testable statements — one per line
5. Flag any User Story that has NO acceptance criteria as **INCOMPLETE**
6. Follow related work items (Parent, Child, Related, Predecessor, Successor) up to 2 levels deep
7. Extract iteration path, area path, and tags for context on team ownership and sprint scope
8. Check the Discussion/History for requirement clarifications that extend the original description
9. Preserve the original work item ID as a traceability reference in all outputs
10. Output format MUST match jira-ticket-reader structure for cross-platform consistency

## Work Item Type Field Mapping

| Work Item Type | Key Fields |
|----------------|------------|
| **User Story** | Title, Description, Acceptance Criteria, Story Points, Iteration, Area Path |
| **Bug** | Title, Repro Steps, System Info, Acceptance Criteria, Severity, Priority |
| **Feature** | Title, Description, Acceptance Criteria, Target Date, Business Value |
| **Epic** | Title, Description, Business Value, Start/Target Date |
| **Task** | Title, Description, Remaining Work, Activity |

## Workflow

1. **Fetch the work item** via Azure DevOps REST API:
   - Endpoint: `GET {org}/{project}/_apis/wit/workitems/{id}?$expand=relations&api-version=7.0`
   - Or use available MCP tools for Azure DevOps if configured
   - Include all fields and relations in the response
2. **Parse by work item type**:
   - User Story: Extract Description + Acceptance Criteria fields
   - Bug: Extract Repro Steps + System Info + Acceptance Criteria
   - Feature/Epic: Extract Description + child work items for full scope
3. **Extract acceptance criteria** as discrete testable statements:
   - ADO stores acceptance criteria in `Microsoft.VSTS.Common.AcceptanceCriteria`
   - Parse HTML content — strip tags, preserve structure (lists, tables)
   - If missing, flag as INCOMPLETE and generate suggested criteria
4. **Fetch related work items**:
   - Parse the `relations` array from the work item response
   - Categorize: Parent, Child, Related, Predecessor, Successor, Tested By
   - Fetch summary details for each related item (batch request if many)
5. **Check discussion history**:
   - Endpoint: `GET {org}/{project}/_apis/wit/workitems/{id}/comments?api-version=7.0-preview`
   - Extract clarifications, scope changes, and PO decisions
6. **Assess completeness** — identify gaps, ambiguities, and missing information
7. **Produce structured output** matching jira-ticket-reader format for consistency

## Handling Features and Epics

When the provided work item is a Feature or Epic:

1. Query child work items using WIQL:
   ```
   SELECT [System.Id], [System.Title], [System.State]
   FROM WorkItemLinks
   WHERE [Source].[System.Id] = {id}
   AND [System.Links.LinkType] = 'System.LinkTypes.Hierarchy-Forward'
   ```
2. Fetch each child's summary, state, and acceptance criteria
3. Group children by state (New, Active, Resolved, Closed)
4. Produce a consolidated requirements summary across all children

## Output Template

```markdown
## Requirements Summary: #[WORK-ITEM-ID] — [Title]

**Source**: Azure DevOps (#[WORK-ITEM-ID])
**Type**: [User Story | Bug | Feature | Epic | Task]
**State**: [Current state]
**Priority**: [Priority]
**Iteration**: [Iteration path]
**Area Path**: [Area path]
**Tags**: [Tag list]

### Acceptance Criteria
- [ ] AC-001: [Testable acceptance criterion]
- [ ] AC-002: [Testable acceptance criterion]
- [ ] AC-003: ...

### Requirements Extracted from Description
- [ ] REQ-001: [Explicit requirement from description]
- [ ] REQ-002: [Implicit requirement inferred from context]

### Clarifications from Discussion
| Date | Author | Clarification |
|------|--------|---------------|
| [date] | [author] | [What was clarified or changed] |

### Related Work Items
| ID | Type | Relationship | State | Title |
|----|------|-------------|-------|-------|
| #[ID] | [User Story/Bug] | [Parent/Child/Related] | [State] | [Title] |

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
