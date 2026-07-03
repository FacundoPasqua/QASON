---
name: test-plan-gen
description: >
  Generates structured test plans from analyzed requirements, including
  scope, categories, scenarios, dependencies, and entry/exit criteria.
  Trigger: After requirements analysis or when asked to plan testing for a feature.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-analyst
  phase: planning
---

# Test Plan Generator

## When to Use

- After a PRD/requirements analysis is complete
- User asks to "create a test plan" or "plan testing for X"
- Starting a new sprint/release with new features to test

## Critical Rules

1. Always include SCOPE — what IS and what IS NOT being tested
2. Scenarios must cover: happy path, edge cases, negative cases, error recovery
3. Include entry/exit criteria — when does testing start and when is it "done"
4. Estimate effort per test category when possible
5. Link scenarios to requirements (traceability)

## Workflow

1. **Review** requirements analysis (or perform one if not available)
2. **Define** test scope and boundaries
3. **Categorize** test types needed (functional, integration, performance, security, etc.)
4. **Design** scenarios per category with priority and risk
5. **Identify** dependencies (test data, environments, external services)
6. **Set** entry/exit criteria
7. **Estimate** effort and suggest automation candidates

## Output Template

```markdown
## Test Plan: [Feature Name]
**Version**: 1.0
**Date**: [date]
**Author**: QA Analyst Agent

### 1. Scope
**In Scope**:
- [what will be tested]

**Out of Scope**:
- [what will NOT be tested + why]

### 2. Test Strategy
| Category | Approach | Priority | Automation Candidate |
|----------|----------|----------|---------------------|
| Functional | Manual + Automated | High | Yes |
| Integration | Automated | High | Yes |
| Performance | Automated | Medium | Yes |
| Security | Manual review + Automated scans | Medium | Partial |
| Accessibility | Manual + axe-core | Low | Partial |

### 3. Test Scenarios
| ID | Requirement | Category | Scenario | Steps | Expected Result | Priority |
|----|-------------|----------|----------|-------|-----------------|----------|
| TS-001 | REQ-001 | Happy Path | [scenario] | [steps] | [result] | High |
| TS-002 | REQ-001 | Edge Case | [scenario] | [steps] | [result] | Medium |

### 4. Test Data Requirements
- [data needed + how to generate/obtain it]

### 5. Environment Requirements
- [environments needed, their configurations]

### 6. Dependencies
- [external services, APIs, databases, third-party tools]

### 7. Entry Criteria
- [ ] Requirements reviewed and approved
- [ ] Test environment available
- [ ] Test data prepared

### 8. Exit Criteria
- [ ] All High priority tests executed
- [ ] No Critical/High severity bugs open
- [ ] Coverage targets met ([X]% functional, [Y]% integration)

### 9. Risks
| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
```
