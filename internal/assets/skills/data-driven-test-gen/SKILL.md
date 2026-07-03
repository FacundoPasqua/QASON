---
name: data-driven-test-gen
description: >
  Creates parameterized and data-driven test specifications using
  equivalence class partitioning, boundary value analysis, and
  structured data tables. Generates comprehensive input combinations
  while keeping the test count manageable.
  Trigger: When a feature has multiple input parameters, complex validation rules, or combinatorial test needs.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-test-designer
  phase: testing
---

# Data-Driven Test Generator

## When to Use

- Feature has multiple input fields with validation rules
- Business logic varies based on input combinations
- Same test flow needs to run with many different data sets
- Equivalence classes and boundary values need systematic identification
- Pairwise/combinatorial testing is needed to manage test explosion

## Critical Rules

1. ALWAYS identify equivalence classes BEFORE generating test data — random data is not data-driven testing
2. Boundary values MUST include: min, min-1, min+1, max, max-1, max+1, zero, empty, null
3. Never generate more combinations than necessary — use pairwise testing for 3+ parameters
4. Each data row MUST have an expected result — data without expectations is not a test
5. Separate test data from test logic — data tables must be independently maintainable
6. Include at least one NEGATIVE test per equivalence class boundary
7. Mark destructive or state-changing test data clearly — execution order may matter

## Equivalence Class Partitioning

```
For each input parameter, identify:

VALID CLASSES (should succeed)
  VC1: Typical valid value (middle of range)
  VC2: Minimum valid value
  VC3: Maximum valid value
  VC4: Special valid formats (if applicable)

INVALID CLASSES (should fail with specific error)
  IC1: Below minimum
  IC2: Above maximum
  IC3: Wrong type (string where number expected)
  IC4: Empty / blank
  IC5: Null / undefined / missing
  IC6: Special characters / injection attempts
  IC7: Invalid format (wrong date format, bad email)

BOUNDARY VALUES (test the edges)
  BV1: min - 1     (invalid)
  BV2: min         (valid)
  BV3: min + 1     (valid)
  BV4: max - 1     (valid)
  BV5: max         (valid)
  BV6: max + 1     (invalid)
  BV7: zero / empty (depends on rules)
```

## Workflow

1. **Identify** all input parameters and their validation rules
2. **Partition** each parameter into equivalence classes:
   - List valid classes with representative values
   - List invalid classes with representative values
   - Identify boundary values
3. **Combine** parameters:
   - For 1-2 parameters: exhaustive combinations (all pairs)
   - For 3+ parameters: pairwise testing (every pair of values appears at least once)
   - For dependent parameters: decision table approach
4. **Build** the data table:
   - One row per test case
   - All input values + expected output/behavior
   - Tags for priority, category, and execution requirements
5. **Add** special data sets:
   - Localization: Unicode, RTL text, multi-byte characters
   - Performance: large strings, maximum-size payloads
   - Security: SQL injection, XSS, path traversal patterns
6. **Validate** the table: no duplicate rows, no missing expected results

## Decision Table Pattern

```
Use when business rules create conditional logic across parameters:

CONDITION STUB          | R1  | R2  | R3  | R4  | R5  |
─────────────────────────────────────────────────────────
User is premium?        | Y   | Y   | N   | N   | N   |
Order total > $100?     | Y   | N   | Y   | Y   | N   |
Has coupon?             | -   | -   | Y   | N   | -   |
─────────────────────────────────────────────────────────
ACTION STUB             |     |     |     |     |     |
─────────────────────────────────────────────────────────
Free shipping?          | Y   | Y   | Y   | N   | N   |
Discount applied?       | 20% | 10% | 15% | 0%  | 0%  |
Loyalty points earned?  | 2x  | 1x  | 1x  | 1x  | 1x  |

"-" means the condition is irrelevant (don't care) for that rule.
```

## Pairwise Testing Guide

```
When exhaustive combinations are too many (e.g., 4 params x 5 values = 625),
use pairwise to ensure every pair of parameter values appears at least once.

Example: 3 parameters with 3 values each
  Exhaustive: 27 combinations
  Pairwise:   9 combinations (covers all pairs)

Pairwise guarantees:
  ✓ Every value of Parameter A appears with every value of Parameter B
  ✓ Every value of Parameter B appears with every value of Parameter C
  ✓ Every value of Parameter A appears with every value of Parameter C

Pairwise does NOT guarantee:
  ✗ Every triple (A, B, C) is tested — some 3-way interactions may be missed
  → For critical 3-way interactions, add explicit rows
```

## Validation (MANDATORY)

After writing test files, you MUST run them:

1. Detect the test runner from `package.json` scripts, `Makefile`, or framework defaults
2. Run ONLY the generated test files (not the full suite)
3. If tests fail: read the error, fix the test, re-run
4. Report: X passed, Y failed, Z skipped

Never deliver tests you haven't run. A test that doesn't execute is not a test.

## Output Format

```markdown
## Data-Driven Test Specification: [Feature Name]

### Parameters
| Parameter | Type | Valid Range | Required | Default |
|-----------|------|-------------|----------|---------|
| [name] | [type] | [constraints] | [Y/N] | [value or N/A] |

### Equivalence Classes
| Parameter | Class ID | Type | Representative Value | Description |
|-----------|----------|------|---------------------|-------------|
| [name] | VC1 | Valid | [value] | [description] |
| [name] | IC1 | Invalid | [value] | [description] |
| [name] | BV1 | Boundary | [value] | [description] |

### Data Table
| ID | [Param1] | [Param2] | [ParamN] | Expected Result | Category | Priority |
|----|----------|----------|----------|-----------------|----------|----------|
| DDT-001 | [value] | [value] | [value] | [expected] | Happy Path | High |
| DDT-002 | [value] | [value] | [value] | [error message] | Negative | High |
| DDT-003 | [value] | [value] | [value] | [expected] | Boundary | Medium |

### Special Data Sets
| ID | Category | Input Data | Purpose | Expected |
|----|----------|-----------|---------|----------|
| SDS-001 | Localization | [Unicode string] | Verify encoding | [expected] |
| SDS-002 | Security | [injection pattern] | Verify sanitization | [expected] |
| SDS-003 | Performance | [large payload] | Verify limits | [expected] |

### Combination Strategy
- Total parameters: [N]
- Combination method: [exhaustive / pairwise / decision table]
- Total test cases: [count]
- Coverage: [what is guaranteed to be covered]
```
