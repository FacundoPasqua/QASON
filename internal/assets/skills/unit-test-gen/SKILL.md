---
name: unit-test-gen
description: >
  Generates unit tests for new or modified functions. Follows AAA pattern,
  mocks external dependencies, and matches project conventions.
  Trigger: When code is added/modified and unit tests are needed.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: development
---

# Unit Test Generator

## When to Use

- New functions or methods are added
- Existing functions are modified (add tests for new behavior)
- PR diff contains untested code
- User asks to "add unit tests" or "test this function"

## Critical Rules

1. Follow the **AAA pattern**: Arrange → Act → Assert
2. Mock external dependencies (DB, HTTP, file system), NOT the system under test
3. Test the **public API**, not implementation details
4. Each test must have a **single assertion focus** (one logical thing being verified)
5. Test names must describe the scenario: `test_[function]_[scenario]_[expected]`
6. NEVER introduce a new test framework — use what the project already has
7. Match the project's existing test file structure and naming conventions
8. Include: happy path, error cases, edge cases, boundary values

## Workflow

1. **Read** the source file and any existing tests
2. **Identify** the project's test framework and conventions
3. **Analyze** each function: inputs, outputs, side effects, error paths
4. **Generate** tests following this priority:
   - Happy path (normal operation)
   - Error handling (what happens when things fail)
   - Edge cases (empty input, null, boundary values)
   - State transitions (if applicable)
5. **Verify** all imports are resolved and tests follow project conventions

## Test Structure

```
// For each function under test:

// Happy path — normal operation
test("[function] returns [expected] when [condition]")

// Error cases — what happens when things fail
test("[function] throws [error] when [invalid condition]")
test("[function] returns [fallback] when [dependency] fails")

// Edge cases — boundaries and special values
test("[function] handles empty input")
test("[function] handles null/undefined")
test("[function] handles maximum value")
test("[function] handles concurrent calls") // if relevant

// State — if the function modifies state
test("[function] updates [state] from [before] to [after]")
test("[function] does not modify [state] when [condition]")
```

## Framework Detection

Before writing tests, detect and match:
- **JavaScript/TypeScript**: jest, vitest, mocha → check `package.json` + existing test files
- **Go**: testing stdlib + testify → check `go.mod` + existing `_test.go` files
- **Python**: pytest, unittest → check `pyproject.toml`/`requirements.txt` + existing test files
- **Java**: JUnit 5, TestNG → check `pom.xml`/`build.gradle`

## Validation (MANDATORY)

After writing test files, you MUST run them:

1. Detect the test runner from `package.json` scripts, `Makefile`, or framework defaults
2. Run ONLY the generated test files (not the full suite)
3. If tests fail: read the error, fix the test, re-run
4. Report: X passed, Y failed, Z skipped

Never deliver tests you haven't run. A test that doesn't execute is not a test.

## Output

Generate test files that:
- Compile and run without modifications
- Follow the project's directory structure (`__tests__/`, `*_test.go`, `test_*.py`, etc.)
- Include all necessary imports
- Use the project's assertion library
- Have meaningful failure messages
