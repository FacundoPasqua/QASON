---
name: integration-test-gen
description: >
  Generates integration tests when API or contract changes are detected.
  Tests real interactions between components, verifies request/response contracts,
  and validates error propagation across service boundaries.
  Trigger: When API contracts change, new endpoints are added, or service interactions are modified.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: development
---

# Integration Test Generator

## When to Use

- API endpoints are added, modified, or deprecated
- Service-to-service contracts change (request/response schemas)
- New dependencies between components are introduced
- User asks to "test the integration", "verify the API", or "check contracts"
- After detecting breaking changes in shared interfaces or DTOs
- When error handling paths between services need validation

## Critical Rules

1. NEVER mock the component under test — mocks are for external boundaries only
2. Each test must verify a REAL interaction between at least two components
3. Always test both the happy path AND at least 3 error scenarios per endpoint
4. Request/response schemas must be validated against the actual contract (OpenAPI, GraphQL schema, protobuf)
5. Test error propagation end-to-end — verify that upstream errors produce correct downstream responses
6. Include timeout and retry behavior tests for async interactions
7. Tests must be idempotent — running them twice produces the same result
8. Always clean up test data in teardown, even if the test fails
9. HTTP status codes must be explicitly asserted, not just response body

## Workflow

1. **Detect** the changed API surface — new/modified endpoints, altered contracts, changed DTOs
2. **Map** the interaction chain — which components call what, in which order
3. **Identify** contract boundaries — request schemas, response schemas, error formats
4. **Generate** happy path tests covering the full request-response cycle
5. **Generate** error path tests — invalid input, auth failures, downstream unavailability, timeouts
6. **Generate** contract validation tests — schema compliance, required fields, type correctness
7. **Add** setup/teardown for test data and external dependencies
8. **Verify** tests are idempotent and isolated from each other

## Validation (MANDATORY)

After writing test files, you MUST run them:

1. Detect the test runner from `package.json` scripts, `Makefile`, or framework defaults
2. Run ONLY the generated test files (not the full suite)
3. If tests fail: read the error, fix the test, re-run
4. Report: X passed, Y failed, Z skipped

Never deliver tests you haven't run. A test that doesn't execute is not a test.

## Output Format

```markdown
## Integration Tests: [API/Service Name]

### Scope
- **Changed endpoints**: [list of endpoints affected]
- **Interaction chain**: [ServiceA] -> [ServiceB] -> [Database/External]
- **Contract source**: [OpenAPI spec path / GraphQL schema / protobuf file]

### Test Cases

#### Happy Path
| ID | Test Name | Method | Endpoint | Input | Expected Status | Expected Response | Validates |
|----|-----------|--------|----------|-------|-----------------|-------------------|-----------|
| INT-001 | [descriptive name] | [GET/POST/...] | [path] | [key input data] | [200/201/...] | [key response fields] | [what contract aspect] |

#### Error Paths
| ID | Test Name | Scenario | Input | Expected Status | Expected Error Response | Error Propagation |
|----|-----------|----------|-------|-----------------|------------------------|-------------------|
| INT-ERR-001 | [descriptive name] | [invalid input / auth failure / timeout / ...] | [input that triggers error] | [400/401/500/...] | [expected error body] | [how error surfaces to caller] |

#### Contract Validation
| ID | Test Name | Contract Aspect | Validation |
|----|-----------|-----------------|------------|
| INT-CTR-001 | [descriptive name] | [required fields / types / enums / ...] | [what is checked] |

### Setup & Teardown
- **Pre-conditions**: [required test data, service state, auth tokens]
- **Setup steps**: [ordered list of setup actions]
- **Teardown steps**: [cleanup actions — must run even on failure]

### Test Code Template

\`\`\`[language]
describe("[API/Service] Integration", () => {
  beforeAll(async () => {
    // Setup: seed test data, obtain auth tokens
  });

  afterAll(async () => {
    // Teardown: clean up test data
  });

  describe("[Endpoint] - Happy Path", () => {
    it("should [expected behavior]", async () => {
      // Arrange: prepare request
      // Act: call the real endpoint
      // Assert: verify status, response body, and side effects
    });
  });

  describe("[Endpoint] - Error Handling", () => {
    it("should return [status] when [error condition]", async () => {
      // Arrange: prepare invalid request or error state
      // Act: call the endpoint
      // Assert: verify error status, error body format, no side effects
    });
  });
});
\`\`\`

### Dependencies
- [services that must be running for these tests]
- [test databases or containers required]
- [environment variables or configuration needed]
```
