---
name: ado-testcase-linker
description: >
  Creates Test Case work items in Azure DevOps Test Plans, defines structured
  test steps, and links them to parent user stories for full traceability.
  Trigger: When test cases need to be documented in Azure DevOps and linked to work items.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-test-designer
  phase: testing
---

# Azure DevOps Test Case Linker

## When to Use

- Test cases have been designed and need to be recorded in Azure DevOps
- User asks to "create test cases in ADO", "link tests to the user story", or "add test coverage in Azure DevOps"
- After test-plan-gen or e2e-test-gen produces test cases that need formal tracking
- Sprint planning requires visibility of test coverage per user story in ADO

## Critical Rules

1. Create test cases as "Test Case" work item type — ADO has native test case support with structured steps
2. Every test case MUST link back to at least one parent user story via "Tests / Tested By" link
3. Use ADO's native test step format: Action + Expected Result pairs — not free-text descriptions
4. Add the tag `qason-generated` to all created test cases
5. Use consistent naming: `TC-[Feature]: [Scenario description]`
6. Include traceability metadata — which requirement (REQ-ID or AC-ID) each test covers
7. NEVER create orphan test cases — every test must trace to a requirement
8. When a Test Plan and Test Suite exist, add test cases to the appropriate suite
9. Set automation status: "Not Automated", "Planned", or "Automated" with script path
10. Check for existing test cases on the story before creating duplicates
11. Apply the `format-mirror` skill BEFORE creating or updating any test case; fall back to the default format skill only if mirroring finds nothing to mirror

## Workflow

1. **Verify the parent work item** via REST API or MCP tools:
   - Confirm the user story exists and note its ID, title, and acceptance criteria
   - Check existing "Tested By" links for already-created test cases
2. **Check for existing Test Plans and Suites**:
   - Query: `GET {org}/{project}/_apis/test/plans?api-version=7.0`
   - If a relevant Test Plan exists, identify the correct Test Suite
   - If no suite exists for this feature, create one
3. **Create each test case** as a work item:
   - Endpoint: `POST {org}/{project}/_apis/wit/workitems/$Test Case?api-version=7.0`
   - Set fields: Title, Area Path, Iteration Path, Tags, Priority
   - Set `Microsoft.VSTS.Common.Priority` matching the parent story
4. **Add test steps** using the test step format:
   - Endpoint: `PATCH {org}/{project}/_apis/wit/workitems/{testCaseId}?api-version=7.0`
   - Field: `Microsoft.VSTS.TCM.Steps` — uses XML format for structured steps
   - Each step has: `<step>` with `<parameterizedString>` for action and expected result
5. **Link to parent user story**:
   - Add "Microsoft.VSTS.Common.TestedBy-Forward" relation
   - This creates the bidirectional "Tests / Tested By" link
6. **Add to Test Suite** (if Test Plan exists):
   - Endpoint: `POST {org}/{project}/_apis/test/plans/{planId}/suites/{suiteId}/testcases?api-version=7.0`
   - Body: `[{ "pointAssignments": [], "workItem": { "id": {testCaseId} } }]`
7. **Set automation status**:
   - Field: `Microsoft.VSTS.TCM.AutomationStatus`
   - Values: "Not Automated" | "Planned" | "Automated"
   - If automated, set `Microsoft.VSTS.TCM.AutomatedTestName` and `AutomatedTestStorage`
8. **Report results** — summarize all created test cases and their links

## Test Steps XML Format

ADO stores test steps as XML in `Microsoft.VSTS.TCM.Steps`:

```xml
<steps id="0" last="3">
  <step id="1" type="ValidateStep">
    <parameterizedString isformatted="true">[Action: what to do]</parameterizedString>
    <parameterizedString isformatted="true">[Expected: what should happen]</parameterizedString>
  </step>
  <step id="2" type="ValidateStep">
    <parameterizedString isformatted="true">[Action: next step]</parameterizedString>
    <parameterizedString isformatted="true">[Expected: next result]</parameterizedString>
  </step>
  <step id="3" type="ValidateStep">
    <parameterizedString isformatted="true">[Action: final step]</parameterizedString>
    <parameterizedString isformatted="true">[Expected: final verification]</parameterizedString>
  </step>
</steps>
```

## Test Case Description Template

**Do NOT hardcode the description here.** Format selection is delegated to the `format-mirror` skill, which resolves in this order:

1. Project override — `.qason/templates/testcase.md`
2. User override — `~/.qason/templates/testcase.md`
3. Mirror from an existing Test Case work item in the same ADO project (mandatory in Acceleration+). Use `wit_query_by_wiql` to find a recent Test Case, then `wit_get_work_item(id, expand="Fields")` to inspect `Microsoft.VSTS.Common.*` fields, `System.Tags`, `Description` (HTML), and the Area Path convention.
4. Fallback to `testcase-format-default` rendered as HTML (ADO's Description is HTML, not markdown — translate the template accordingly).

ADO specifics that `format-mirror` MUST discover from the reference TC:
- Whether Description is populated at all (some teams store everything in `Microsoft.VSTS.TCM.Steps`)
- Area Path + Iteration Path conventions
- Which tags are project-mandatory (beyond `qason-generated`)
- Whether Priority is a required field at create time

Provenance requirement: append `Format source: mirrored-from-#12345` / `project-override` / `default-template` to your output after creation.

## Batch Creation Summary

After creating multiple test cases, output:

```markdown
## Test Cases Created in Azure DevOps

**Parent Story**: #[WORK-ITEM-ID] — [Story title]
**Test Plan**: [Plan name] (if applicable)
**Test Suite**: [Suite name] (if applicable)
**Test Cases Created**: [count]

| Test Case ID | Scenario | Covers | Type | Automation |
|-------------|----------|--------|------|------------|
| #[ID] | [Description] | [AC/REQ IDs] | [Type] | [Not Automated/Planned/Automated] |
| #[ID] | [Description] | [AC/REQ IDs] | [Type] | [Not Automated/Planned/Automated] |

### Coverage Matrix
| Acceptance Criteria | Test Cases | Status |
|---------------------|------------|--------|
| AC-001: [criteria] | #[ID], #[ID] | Covered |
| AC-002: [criteria] | — | GAP |

### Gaps Identified
- [Any acceptance criteria not covered and recommended action]
```
