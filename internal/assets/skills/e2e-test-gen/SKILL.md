---
name: e2e-test-gen
description: >
  Designs end-to-end test cases covering complete user journeys through the
  application. Uses Given/When/Then format, grouped by feature area, with
  setup/teardown steps and full CRUD lifecycle coverage.
  Trigger: When designing e2e tests, user journey tests, or acceptance test suites.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-test-designer
  phase: testing
---

# E2E Test Generator

## When to Use

- Designing acceptance tests for a new feature or user story
- User asks to "write e2e tests", "test the user journey", or "create acceptance tests"
- Before a release to ensure critical user flows are covered
- When translating requirements or specs into executable test scenarios
- After feature completion to validate the full workflow end-to-end

## Critical Rules

1. NEVER test implementation details — test what the USER sees and does
2. Every test must represent a complete user journey from entry point to outcome
3. Always use Given/When/Then format — no exceptions
4. Always cover the full CRUD lifecycle when the feature involves data management
5. Include BOTH the happy path and at least 2 negative paths per feature
6. Setup must create ALL preconditions explicitly — never assume prior test state
7. Teardown must restore the system to its original state — tests must be independent
8. Group tests by feature area, NOT by page or component
9. Each test must be runnable in isolation — no dependency on test execution order
10. Include realistic test data, not placeholder values like "test123" or "foo@bar.com"

## Workflow

1. **Identify** the user journeys — what are the complete flows a user performs?
2. **Map** each journey to entry point, steps, decision points, and expected outcomes
3. **Group** journeys by feature area (e.g., Authentication, Checkout, Profile Management)
4. **Design** happy path scenarios first — the primary success flow
5. **Design** negative scenarios — invalid input, unauthorized access, edge cases
6. **Design** CRUD lifecycle — create, read, update, delete in sequence where applicable
7. **Specify** setup and teardown for each test group
8. **Add** assertions for each step — what should the user see at every point?
9. **Review** for independence — can each test run alone without prior tests?

## Validation (MANDATORY)

After writing test files, you MUST run them:

1. Detect the test runner from `package.json` scripts, `Makefile`, or framework defaults
2. Run ONLY the generated test files (not the full suite)
3. If tests fail: read the error, fix the test, re-run
4. Report: X passed, Y failed, Z skipped

Never deliver tests you haven't run. A test that doesn't execute is not a test.

## Output Format

```markdown
## E2E Test Suite: [Feature Area]

### Overview
- **Feature**: [what this test suite covers]
- **User roles**: [which user types are involved]
- **Prerequisites**: [system state required before any test runs]

### Setup (runs before each test group)
1. [setup step — e.g., create test user with specific role]
2. [setup step — e.g., seed required reference data]
3. [setup step — e.g., navigate to starting page]

### Teardown (runs after each test group)
1. [cleanup step — e.g., delete created test data]
2. [cleanup step — e.g., reset user state]

---

### Feature Group: [Feature Name]

#### E2E-001: [Descriptive test name — happy path]
**Priority**: Critical / High / Medium

**Given** [precondition — the initial state of the system]
  - [additional context if needed]

**When** [user action — what the user does, step by step]
  1. [action step]
  2. [action step]
  3. [action step]

**Then** [expected outcome — what the user should see]
  - [ ] [assertion: specific UI element, message, or state change]
  - [ ] [assertion: data persisted correctly]
  - [ ] [assertion: navigation landed on correct page]

---

#### E2E-002: [Descriptive test name — negative path]
**Priority**: High / Medium

**Given** [precondition]

**When** [user performs invalid or unauthorized action]
  1. [action step]

**Then** [expected error handling]
  - [ ] [assertion: error message displayed]
  - [ ] [assertion: no data was modified]
  - [ ] [assertion: user remains on current page]

---

### CRUD Lifecycle: [Entity Name]

#### E2E-CRUD-001: Create [entity]
**Given** [user is authenticated with appropriate permissions]

**When** the user creates a new [entity]
  1. [navigate to creation form]
  2. [fill in required fields with valid data]
  3. [submit the form]

**Then** the [entity] is created successfully
  - [ ] [success confirmation displayed]
  - [ ] [entity appears in the list view]
  - [ ] [entity details match input data]

#### E2E-CRUD-002: Read [entity]
**Given** [an entity exists in the system]

**When** the user views the [entity]
  1. [navigate to list view]
  2. [select the entity]

**Then** all [entity] details are displayed correctly
  - [ ] [all fields are visible and accurate]
  - [ ] [related data is loaded]

#### E2E-CRUD-003: Update [entity]
**Given** [an entity exists in the system]

**When** the user modifies the [entity]
  1. [navigate to entity detail/edit page]
  2. [modify specific fields]
  3. [save changes]

**Then** the [entity] is updated successfully
  - [ ] [success confirmation displayed]
  - [ ] [updated values reflected in detail view]
  - [ ] [updated values reflected in list view]

#### E2E-CRUD-004: Delete [entity]
**Given** [an entity exists in the system with no blocking dependencies]

**When** the user deletes the [entity]
  1. [navigate to entity]
  2. [initiate delete action]
  3. [confirm deletion]

**Then** the [entity] is removed
  - [ ] [confirmation of deletion displayed]
  - [ ] [entity no longer appears in list view]
  - [ ] [related data handled appropriately (cascaded or orphan-protected)]

---

### Test Data Requirements
| Entity | Field | Value | Notes |
|--------|-------|-------|-------|
| [entity] | [field name] | [realistic test value] | [constraints or format] |

### Environment Requirements
- [browser/device requirements]
- [backend services that must be running]
- [test database state or seed data]
- [feature flags or configuration required]
```
