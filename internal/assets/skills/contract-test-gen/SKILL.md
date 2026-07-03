---
name: contract-test-gen
description: >
  Generates consumer-driven contract tests using the Pact framework.
  Identifies consumer-provider relationships, creates consumer tests with
  expected interactions, and generates provider verification. Supports
  REST, GraphQL, and basic async messaging.
  Trigger: When contract testing, consumer-driven testing, or API compatibility checks are needed.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-test-designer
  phase: testing
---

# Contract Test Generator

## When to Use

- Microservices communicate via REST or GraphQL APIs
- Multiple teams own different services that integrate with each other
- API changes risk breaking downstream consumers
- User asks for "contract test", "Pact test", "consumer-driven test", or "API compatibility"
- New service-to-service integration is being built
- API versioning or migration is planned

## Critical Rules

1. Use **Pact** as the contract testing framework — it is the industry standard for consumer-driven contracts
2. Consumer tests define the contract — the provider verifies it. NEVER reverse this flow
3. Each consumer-provider pair gets its **own Pact file** — do not combine multiple providers
4. Use **matchers** for flexible matching (`like()`, `eachLike()`, `term()`) — do not assert exact values for dynamic fields
5. Provider states MUST be idempotent — running them twice produces the same result
6. NEVER include implementation details in contracts — only the interface (request shape + response shape)
7. Test the **minimum required fields**, not every field the provider returns
8. Include contract tests in the **CI pipeline** — they are most valuable as automated gates
9. Detect the project language and use the appropriate Pact library

## Workflow

1. **Detect** the project stack and Pact compatibility:
   - JavaScript/TypeScript: `@pact-foundation/pact` (check `package.json`)
   - Go: `github.com/pact-foundation/pact-go`
   - Python: `pact-python`
   - Java: `au.com.dius.pact` (check `pom.xml`/`build.gradle`)
   - If not installed: add the dependency and note it in the output

2. **Identify** consumer-provider relationships:
   - Scan for HTTP client calls (fetch, axios, http.Client, requests)
   - Scan for service configuration (base URLs, service registry)
   - Parse OpenAPI specs if available to understand provider capabilities
   - Map: which service (consumer) calls which service (provider) and what endpoints

3. **Generate** consumer tests (one per provider):

   For each consumer-provider interaction:
   - Define the **provider state** (precondition): e.g., "a user with ID 123 exists"
   - Define the **request**: method, path, headers, query params, body
   - Define the **expected response**: status, headers, body shape with matchers
   - Use matchers for dynamic values:
     ```
     like(value)       — matches type, not exact value
     eachLike(example) — array with at least one element matching shape
     term(regex, gen)  — matches regex, generates example value
     integer()         — any integer
     uuid()            — any UUID format
     iso8601DateTime() — any ISO 8601 datetime
     ```

   Structure each test:
   ```
   interaction: "[consumer] fetches [resource] from [provider]"
   given: "[provider state description]"
   uponReceiving: "[human-readable interaction name]"
   withRequest:
     method: [HTTP method]
     path: [endpoint path]
     headers: [required headers]
     body: [request body if applicable]
   willRespondWith:
     status: [expected status]
     headers: [expected headers]
     body: [response shape with matchers]
   ```

4. **Generate** provider verification tests:
   - Set up provider states handler (create/reset test data per state)
   - Configure Pact verification to:
     - Load contracts from local pact files OR Pact Broker
     - Run each interaction against the real provider (test instance)
     - Verify all consumer contracts are satisfied
   - Include state management:
     ```
     providerStates:
       "a user with ID 123 exists":
         setup: [create test user]
         teardown: [delete test user]
     ```

5. **Handle** different API styles:

   **REST APIs**:
   - Standard request/response contract
   - Include content-type negotiation
   - Test both success and error responses

   **GraphQL endpoints**:
   - Contract on the query/mutation shape and response structure
   - Use body matchers for the GraphQL query string
   - Verify response data shape matches consumer expectations

   **Async messaging** (basic):
   - Define message shape (producer → consumer)
   - Consumer test: verify message handler processes the expected message shape
   - Producer test: verify published messages match the contract

6. **Configure** Pact Broker (if applicable):
   - Generate broker configuration
   - Add publish step to consumer CI
   - Add verification step to provider CI
   - Configure webhook for provider verification on new contracts
   - Set up can-i-deploy checks before deployment

## Contract Structure Reference

```
Consumer Test (defines the contract):
  ┌──────────────────────────────────┐
  │ Provider State (precondition)    │
  │ Request (what consumer sends)    │
  │ Expected Response (what it needs)│
  └──────────────────────────────────┘
          │
          ▼ generates
  ┌──────────────────────────────────┐
  │ Pact File (JSON contract)        │
  │ - consumer name                  │
  │ - provider name                  │
  │ - interactions[]                 │
  └──────────────────────────────────┘
          │
          ▼ verified by
  ┌──────────────────────────────────┐
  │ Provider Verification            │
  │ - loads pact file                │
  │ - sets up provider states        │
  │ - replays each interaction       │
  │ - asserts response matches       │
  └──────────────────────────────────┘
```

## Validation (MANDATORY)

After writing test files, you MUST run them:

1. Detect the test runner from `package.json` scripts, `Makefile`, or framework defaults
2. Run ONLY the generated test files (not the full suite)
3. If tests fail: read the error, fix the test, re-run
4. Report: X passed, Y failed, Z skipped

Never deliver tests you haven't run. A test that doesn't execute is not a test.

## Output Template

```markdown
## Contract Test Suite: [Consumer] → [Provider]

### Relationship
- Consumer: [service name] — [what it needs from the provider]
- Provider: [service name] — [what it exposes]

### Files Generated
- `contracts/consumer/[consumer]-[provider].pact.test.{ext}` — Consumer contract tests
- `contracts/provider/[provider]-verification.test.{ext}` — Provider verification tests
- `contracts/provider/states.{ext}` — Provider state handlers
- `contracts/pact-config.{ext}` — Pact Broker configuration (if applicable)

### Interactions Covered
| ID | Consumer Action | Provider State | Method | Path | Status |
|----|----------------|----------------|--------|------|--------|
| CT-001 | Fetch user profile | User 123 exists | GET | /api/users/123 | 200 |
| CT-002 | Fetch missing user | User 999 does not exist | GET | /api/users/999 | 404 |

### Provider States Required
| State | Setup Action | Teardown Action |
|-------|-------------|-----------------|
| "User 123 exists" | Create test user | Delete test user |

### CI Integration
- Consumer pipeline: run contract tests → publish pact to broker
- Provider pipeline: verify pacts → report results to broker
- Deployment gate: `can-i-deploy --pacticipant [consumer] --version [sha]`
```
