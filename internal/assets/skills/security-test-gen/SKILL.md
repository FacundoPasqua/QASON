---
name: security-test-gen
description: >
  Generates security test checks based on OWASP Top 10 methodology.
  Analyzes API endpoints and UI forms for common vulnerabilities including
  injection, XSS, auth bypass, and sensitive data exposure. Uses the
  project's existing test framework.
  Trigger: When security testing, vulnerability scanning, or OWASP compliance is needed.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: testing
---

# Security Test Generator

## When to Use

- New API endpoints handle user input or authentication
- UI forms accept and process user data
- Pre-release security validation is required
- User asks for "security test", "OWASP check", "pen test", or "vulnerability scan"
- Changes affect authentication, authorization, or data handling

## Critical Rules

1. ALWAYS reference the specific **OWASP Top 10 category** for each test
2. Use the project's **existing test framework** — do not introduce new dependencies for security tests
3. NEVER include actual exploit payloads that could cause damage — use detection payloads only
4. Test against a **local/staging environment only** — include safeguards against running in production
5. Mark all security tests with a dedicated tag/label (e.g., `@security`, `security:true`)
6. Include both **automated checks** (script-verifiable) and a **manual checklist** for findings that require human judgment
7. Categorize every finding by severity: Critical, High, Medium, Low
8. NEVER log or store sensitive data (passwords, tokens) in test output

## Workflow

1. **Analyze** the application surface:
   - Scan for API endpoint definitions (routes, controllers, handlers)
   - Identify UI forms and input fields
   - Map authentication and authorization flows
   - Identify data storage and transmission patterns
   - Check for existing security headers configuration

2. **Generate** injection tests (OWASP A03 - Injection):
   - SQL injection: `' OR 1=1--`, `'; DROP TABLE--`, parameterized query bypass
   - NoSQL injection: `{"$gt": ""}`, `{"$ne": null}` in JSON payloads
   - Command injection: `; ls`, `| cat /etc/passwd` in text fields
   - LDAP injection if applicable
   - For each: send payload, verify it is rejected or sanitized (status 400, no data leak)

3. **Generate** XSS tests (OWASP A03 - Injection):
   - Reflected XSS: inject `<script>alert(1)</script>` in query params and form fields
   - Stored XSS: submit payloads via forms, verify sanitization on retrieval
   - DOM XSS: test URL fragments and client-side rendering
   - Verify: Content-Security-Policy header, output encoding, input sanitization

4. **Generate** authentication tests (OWASP A07 - Identification and Auth Failures):
   - Missing auth: access protected endpoints without token
   - Expired tokens: use expired JWT/session
   - Tampered tokens: modify JWT payload without re-signing
   - Brute force: verify rate limiting on login endpoint
   - Default credentials: check common admin/admin, test/test
   - Password policy: test weak password acceptance

5. **Generate** authorization tests (OWASP A01 - Broken Access Control):
   - IDOR: access resources using another user's ID
   - Privilege escalation: perform admin actions with regular user token
   - Horizontal access: access peer user's data
   - Missing function-level access control: call admin APIs with user role
   - Directory traversal: `../../../etc/passwd` in file paths

6. **Generate** data protection tests (OWASP A02 - Cryptographic Failures):
   - Sensitive data in responses: check for passwords, tokens, PII in API responses
   - HTTPS enforcement: verify HTTP redirects to HTTPS
   - Security headers: Strict-Transport-Security, X-Content-Type-Options, X-Frame-Options
   - Cookie flags: Secure, HttpOnly, SameSite on session cookies
   - Error messages: verify no stack traces or internal details in error responses

7. **Generate** CSRF tests (OWASP A01 - Broken Access Control):
   - State-changing requests without CSRF token
   - CSRF token reuse/replay
   - Cross-origin request handling (CORS misconfiguration)

8. **Generate** rate limiting tests:
   - Send rapid requests to login endpoint — verify 429 response
   - Send rapid requests to sensitive operations — verify throttling
   - Test rate limit reset behavior

## Severity Classification

```
| Severity | Criteria                                          | Examples                           |
|----------|---------------------------------------------------|------------------------------------|
| Critical | Direct data breach, RCE, auth bypass              | SQL injection, admin access bypass |
| High     | Significant data exposure, privilege escalation    | IDOR, stored XSS, CSRF            |
| Medium   | Limited exposure, requires specific conditions     | Reflected XSS, missing headers     |
| Low      | Informational, defense-in-depth                    | Verbose errors, missing CSP        |
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
## Security Test Suite: [Application Name]

### OWASP Coverage
| Category | ID | Tests | Status |
|----------|----|-------|--------|
| Broken Access Control | A01 | [count] | ⬜ |
| Cryptographic Failures | A02 | [count] | ⬜ |
| Injection | A03 | [count] | ⬜ |
| Insecure Design | A04 | [count] | ⬜ |
| Security Misconfiguration | A05 | [count] | ⬜ |
| Vulnerable Components | A06 | [count] | ⬜ |
| Auth Failures | A07 | [count] | ⬜ |
| Data Integrity Failures | A08 | [count] | ⬜ |
| Logging Failures | A09 | [count] | ⬜ |
| SSRF | A10 | [count] | ⬜ |

### Test File: `security/[app]-security.test.{ext}`
### Manual Checklist: `security/[app]-security-checklist.md`

### Findings Summary
| ID | Severity | OWASP | Description | Endpoint | Status |
|----|----------|-------|-------------|----------|--------|
| SEC-001 | Critical | A03 | [finding] | [endpoint] | ⬜ |
```
