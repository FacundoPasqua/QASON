---
name: risk-matrix-gen
description: >
  Creates risk matrices (likelihood x impact) for test prioritization.
  Maps components and features to risk scores, derives test priority levels,
  and identifies areas safe to skip when time-constrained.
  Trigger: When planning test effort allocation, triaging test scope, or assessing feature risk.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-analyst
  phase: planning
---

# Risk Matrix Generator

## When to Use

- Before test planning to prioritize where to focus effort
- When time or resources are constrained and test scope must be reduced
- User asks to "assess risk", "prioritize testing", or "what should we test first"
- When a new release includes changes across multiple components
- After requirements analysis to rank features by test criticality

## Critical Rules

1. NEVER assign risk scores without justification — every score needs a rationale
2. Both likelihood AND impact must be evaluated independently on a 1-5 scale
3. Always consider historical defect data when available — past bugs predict future bugs
4. Business-critical paths (payments, auth, data integrity) default to High impact minimum
5. Risk scores must account for change velocity — frequently changed code is higher risk
6. When time-constrained, NEVER skip High-risk items — cut from Low first, then Medium
7. Include both functional and non-functional risk dimensions

## Workflow

1. **Inventory** all components, features, or user flows in scope
2. **Assess likelihood** (1-5) of defects for each item based on complexity, change frequency, historical defects, and developer familiarity
3. **Assess impact** (1-5) of failure for each item based on user reach, business criticality, data sensitivity, and recovery difficulty
4. **Calculate** risk score (likelihood x impact) and assign priority tier
5. **Rank** items by risk score descending
6. **Identify** the cut line — items below it can be deferred under time pressure
7. **Map** risk tiers to recommended test types (manual, automated, exploratory)

## Output Template

```markdown
## Risk Matrix: [Project/Release Name]

### Risk Scale Reference
| Score | Likelihood | Impact |
|-------|-----------|--------|
| 5 | Almost certain — changes every sprint, history of defects | Catastrophic — data loss, security breach, revenue loss |
| 4 | Likely — complex logic, recent refactoring | Major — core workflow blocked, significant UX degradation |
| 3 | Possible — moderate complexity, some change history | Moderate — workaround exists, partial functionality loss |
| 2 | Unlikely — stable code, well-tested | Minor — cosmetic issue, edge case only |
| 1 | Rare — no recent changes, simple logic | Negligible — no user-visible effect |

### Risk Assessment

| Component/Feature | Likelihood (1-5) | Impact (1-5) | Risk Score | Priority | Rationale |
|-------------------|-------------------|---------------|------------|----------|-----------|
| [component name] | [score] | [score] | [L x I] | Critical/High/Medium/Low | [why these scores] |

### Priority Tiers
- **Critical (20-25)**: Must test — block release if untested
- **High (12-19)**: Should test — automated + exploratory coverage
- **Medium (6-11)**: Test if time permits — targeted manual checks
- **Low (1-5)**: Safe to defer — monitor via production observability

### Time-Constrained Test Plan
| Available Time | Test Scope | Items Covered | Items Deferred |
|----------------|------------|---------------|----------------|
| Full | All tiers | All | None |
| 75% | Critical + High + partial Medium | [list] | [list] |
| 50% | Critical + High only | [list] | [list] |
| 25% | Critical only | [list] | [list] |

### Recommended Test Types per Tier
- **Critical**: Automated regression + manual exploratory + load testing
- **High**: Automated integration + targeted manual
- **Medium**: Automated unit/integration only
- **Low**: Existing automated coverage sufficient

### Risk Mitigation Notes
- [specific risks that need monitoring even if testing is deferred]
- [areas where observability/alerting can compensate for reduced testing]
```
