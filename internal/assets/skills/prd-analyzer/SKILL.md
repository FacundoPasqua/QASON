---
name: prd-analyzer
description: >
  Analyzes PRDs, user stories, and tickets to extract requirements,
  identify gaps, generate edge cases, and produce clarification questions.
  Trigger: When given a PRD, user story, epic, or feature specification.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-analyst
  phase: planning
---

# PRD Analyzer

## When to Use

- User provides a PRD, user story, ticket, or feature specification
- User asks to "analyze requirements" or "review this spec"
- Before starting test planning for a new feature

## Critical Rules

1. NEVER assume the requirements are complete — always identify gaps
2. Extract BOTH explicit and implicit requirements
3. Always generate at least 5 edge cases per feature
4. Include non-functional requirements even if not mentioned (performance, security, accessibility, i18n)
5. Questions for the PO must explain WHY the answer matters for testing

## Example: Explicit vs Implicit Extraction

> Spec says: "Users can sign up with email and password."

- **Explicit**: REQ-001 — a sign-up form accepts email + password and creates an account
- **Implicit**: IREQ-001 — email format validation (what counts as a valid email?)
- **Implicit**: IREQ-002 — duplicate-account handling (sign up twice with the same email)
- **Implicit**: IREQ-003 — password policy (min length, complexity, max length)
- **Missing criteria**: confirmation email? rate limiting? what does the error state look like?

One sentence of spec routinely hides 3-5 implicit requirements — extract them every time.

## Workflow

1. **Read** the full PRD/ticket carefully
2. **Extract** explicit requirements (what the document says)
3. **Infer** implicit requirements (what it assumes but doesn't state)
4. **Identify** missing acceptance criteria
5. **Generate** edge cases and boundary scenarios
6. **Flag** contradictions, ambiguities, or under-specified areas
7. **Produce** clarification questions ranked by impact

## Output Template

```markdown
## Requirements Analysis: [Feature Name]

### Explicit Requirements
- [ ] REQ-001: [requirement from the spec]
- [ ] REQ-002: ...

### Implicit Requirements
- [ ] IREQ-001: [inferred requirement + why it's necessary]
- [ ] IREQ-002: ...

### Missing Acceptance Criteria
- [criteria that should exist but don't]

### Edge Cases
| ID | Scenario | Risk | Priority |
|----|----------|------|----------|
| EC-001 | [edge case description] | [what could go wrong] | High/Med/Low |

### Ambiguities & Contradictions
- [specific contradiction + both interpretations]

### Clarification Questions
1. **[Question]** — Impact: [why this matters for testing]
2. ...

### Non-Functional Considerations
- **Performance**: [relevant performance scenarios]
- **Security**: [relevant security scenarios]
- **Accessibility**: [relevant a11y scenarios]
- **i18n**: [relevant internationalization scenarios]
```
