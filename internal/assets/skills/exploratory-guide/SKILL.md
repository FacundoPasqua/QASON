---
name: exploratory-guide
description: >
  Generates structured exploratory testing charters with time boxes,
  heuristics (SFDIPOT), and note templates. Guides testers through
  systematic exploration while preserving the creative freedom of ET.
  Trigger: When exploratory testing is needed for a feature, area, or risk zone.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-test-designer
  phase: testing
---

# Exploratory Testing Guide

## When to Use

- New feature delivered without detailed requirements or with ambiguous specs
- Existing feature with low test coverage or high defect history
- Pre-release exploration to find issues scripted tests might miss
- User reports vague symptoms ("it feels slow", "something is off")
- After a major refactor to validate behavior preservation

## Critical Rules

1. ALWAYS define a charter BEFORE starting exploration — undirected testing is not exploratory testing
2. Every charter MUST have a time box — 30 to 90 minutes maximum per session
3. Use SFDIPOT heuristics as the primary lens for test idea generation
4. Document observations AS YOU GO, not from memory after the session
5. Separate observations (what you saw) from interpretations (what you think it means)
6. A bug found during ET is not "done" until it has a reproduction path documented
7. NEVER skip debriefing — the charter output is the deliverable, not just the bugs

## SFDIPOT Heuristic Framework

```
S — Structure
    What is it made of? Components, modules, files, data structures.
    → Test internal consistency, naming, organization, dead code.

F — Function
    What does it do? Features, operations, user workflows.
    → Test each function works, edge cases, error handling.

D — Data
    What data does it process? Inputs, outputs, persistence, transformations.
    → Test boundaries, invalid data, encoding, large datasets, empty states.

I — Interfaces
    How does it connect? APIs, UI elements, integrations, protocols.
    → Test interaction points, handoffs, format mismatches, timeouts.

P — Platform
    What does it run on? OS, browser, device, network, dependencies.
    → Test cross-platform, responsive, offline, low bandwidth, version compat.

O — Operations
    How is it managed? Deploy, monitor, backup, config, logging.
    → Test install/upgrade, config changes, log output, recovery scenarios.

T — Time
    How does it behave over time? Concurrency, scheduling, aging, sequences.
    → Test race conditions, session expiry, date boundaries, long-running ops.
```

## Workflow

1. **Define** the charter:
   - Target: what area/feature to explore
   - Mission: what you are trying to learn or find
   - Heuristics: which SFDIPOT categories to focus on
   - Time box: session duration (default: 60 min)
2. **Prepare** the environment:
   - Test data prerequisites
   - Tools needed (browser devtools, proxy, logs)
   - Access and permissions
3. **Explore** systematically:
   - Follow the chosen heuristics as guides, not scripts
   - When you find something interesting, DIVE DEEPER before moving on
   - Take timestamped notes using the note template
   - Capture screenshots/recordings for anything unexpected
4. **Debrief** immediately after the session:
   - Summarize findings per heuristic category
   - Classify: Bug / Risk / Question / Observation
   - Identify areas that need deeper investigation
5. **Report** using the output template

## Note Template (During Session)

```
[HH:MM] HEURISTIC: [S/F/D/I/P/O/T]
ACTION: What I did
OBSERVATION: What happened
EXPECTED: What I expected (if different)
CLASSIFICATION: Bug | Risk | Question | Observation
SEVERITY: Critical | High | Medium | Low
SCREENSHOT: [filename or N/A]
FOLLOW-UP: [next action or N/A]
```

## Output Template

```markdown
## Exploratory Testing Charter Report

### Charter
| Field | Value |
|-------|-------|
| Target | [feature/area under test] |
| Mission | [what we aimed to discover] |
| Heuristics | [SFDIPOT categories used] |
| Time box | [planned duration] |
| Actual duration | [actual duration] |
| Tester | [name/role] |

### Session Summary
- Areas explored: [list]
- Areas NOT explored (deferred): [list]
- Overall confidence: [Low/Medium/High]

### Findings
| ID | Time | Heuristic | Type | Summary | Severity | Status |
|----|------|-----------|------|---------|----------|--------|
| ET-001 | [HH:MM] | [S/F/D/I/P/O/T] | [Bug/Risk/Question/Observation] | [description] | [severity] | [New/Reported/Deferred] |

### Bugs Found
| ID | Summary | Reproduction Steps | Severity | Priority |
|----|---------|-------------------|----------|----------|
| BUG-001 | [summary] | [steps] | [severity] | [priority] |

### Risks Identified
| ID | Risk Description | Likelihood | Impact | Suggested Mitigation |
|----|-----------------|------------|--------|---------------------|
| RISK-001 | [description] | [H/M/L] | [H/M/L] | [action] |

### Recommended Follow-Up
- [ ] [Additional charter needed for area X]
- [ ] [Scripted test needed for scenario Y]
- [ ] [Bug ticket for finding Z]
```
