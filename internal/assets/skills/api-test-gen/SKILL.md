---
name: api-test-gen
description: >
  Generates comprehensive API test suites from OpenAPI specs or endpoint
  descriptions. Covers happy paths, error handling, boundary values,
  auth flows, and contract validation.
  Trigger: When given an OpenAPI spec, Swagger doc, or API endpoint description.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-test-designer
  phase: testing
---

# API Test Generator

## When to Use

- User provides an OpenAPI/Swagger specification
- User describes API endpoints to test
- New API endpoints are added or existing ones modified
- Integration testing for microservice communication

## Critical Rules

1. ALWAYS test both success AND failure paths for every endpoint
2. Include auth scenarios: valid token, expired token, missing token, wrong role
3. Test boundary values for ALL input fields (min, max, empty, null, special chars)
4. Validate response schemas strictly — extra fields are bugs too
5. Test idempotency where applicable (PUT, DELETE)
6. Include rate limiting and concurrent request scenarios
7. Never hardcode test data — use variables/factories

## Workflow

1. **Parse** the OpenAPI spec or endpoint description
2. **Map** each endpoint to test categories
3. **Generate** test cases per endpoint:
   - Happy path (valid request → expected response)
   - Error cases (4xx for each validation rule)
   - Auth permutations
   - Boundary values per field
   - Schema validation
4. **Add** cross-endpoint scenarios (CRUD lifecycle, dependent endpoints)
5. **Include** setup/teardown (create test data, clean up after)

## Test Categories Per Endpoint

```
GET /resource
  ✓ 200 — returns expected data
  ✓ 200 — pagination works (first page, last page, out of range)
  ✓ 200 — filtering works per filter parameter
  ✓ 200 — sorting works per sort parameter
  ✓ 401 — missing/invalid auth
  ✓ 403 — insufficient permissions
  ✓ 404 — resource not found

POST /resource
  ✓ 201 — valid payload creates resource
  ✓ 201 — response matches request + generated fields
  ✓ 400 — missing required fields (one per field)
  ✓ 400 — invalid field types
  ✓ 400 — boundary violations (too long, too short, negative)
  ✓ 409 — duplicate resource (if uniqueness constraint)
  ✓ 401/403 — auth scenarios

PUT /resource/:id
  ✓ 200 — valid update
  ✓ 200 — partial update (if supported)
  ✓ 200 — idempotent (same PUT twice = same result)
  ✓ 404 — non-existent ID
  ✓ 400 — validation errors

DELETE /resource/:id
  ✓ 204 — successful deletion
  ✓ 204 — idempotent (delete twice = no error)
  ✓ 404 — non-existent ID (if not idempotent)
  ✓ 403 — cannot delete others' resources
```

## Validation (MANDATORY)

After writing test files, you MUST run them:

1. Detect the test runner from `package.json` scripts, `Makefile`, or framework defaults
2. Run ONLY the generated test files (not the full suite)
3. If tests fail: read the error, fix the test, re-run
4. Report: X passed, Y failed, Z skipped

Never deliver tests you haven't run. A test that doesn't execute is not a test.

## Output Format

Generate tests in the detected framework. If no framework detected, use the project's language with standard HTTP testing patterns.

```markdown
## API Test Suite: [API Name]

### Endpoint: [METHOD] [path]
| ID | Category | Description | Request | Expected | Priority |
|----|----------|-------------|---------|----------|----------|
| API-001 | Happy Path | [desc] | [req details] | [status + body] | High |
```
