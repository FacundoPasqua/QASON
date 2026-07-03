---
name: performance-test-gen
description: >
  Generates performance and load test scripts targeting key application
  endpoints. Detects available tools (k6, Artillery, JMeter, Gatling, Locust),
  generates multi-scenario scripts with thresholds, and produces test data
  for realistic load simulation.
  Trigger: When performance testing, load testing, or stress testing is needed.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: testing
---

# Performance Test Generator

## When to Use

- New API endpoints need load testing before release
- Performance regression testing is required for critical paths
- Capacity planning requires understanding system limits
- User asks for "load test", "stress test", "performance test", or "benchmark"
- Pre-release validation of SLAs and response time requirements

## Critical Rules

1. ALWAYS detect available performance tools before generating scripts
2. If no tool is found, default to **k6** and include installation instructions
3. Generate scripts from OpenAPI specs when available — fall back to endpoint descriptions
4. Every script MUST include defined **thresholds** (p95, error rate, throughput)
5. NEVER hardcode authentication tokens — use environment variables or config files
6. Include realistic **think time** between requests (1-5 seconds) to simulate real users
7. Generate **test data** separately — do not embed large datasets in the script
8. Always include a **smoke scenario** (1 VU) to verify the script works before scaling
9. Document how to run the test and interpret results

## Workflow

1. **Detect** available performance testing tools:
   - Check for `k6` binary in PATH
   - Check `package.json` for `artillery`
   - Check for `jmeter` binary or `.jmx` files in the project
   - Check `build.gradle`/`build.sbt` for Gatling dependencies
   - Check `requirements.txt`/`pyproject.toml` for `locust`
   - If none found: recommend k6, provide install command (`brew install k6` / `choco install k6`)

2. **Identify** target endpoints:
   - Parse OpenAPI/Swagger spec if available
   - Scan existing test files for API calls
   - Ask user to confirm critical paths to test
   - Identify authentication requirements

3. **Generate** test data:
   - Create CSV/JSON files with realistic test data
   - Include enough records for the maximum VU count (at least 100 unique entries)
   - Cover edge cases: long strings, unicode, special characters
   - Generate auth tokens/credentials for test users

4. **Create** test scenarios (ALL four must be included):

   **Smoke Test** (validation):
   - 1 virtual user, 1 minute duration
   - Purpose: verify script correctness before scaling
   - Thresholds: p95 < 2s, error rate = 0%

   **Load Test** (expected traffic):
   - Ramp from 0 to expected concurrent users over 2 minutes
   - Hold at target for 5 minutes
   - Ramp down over 1 minute
   - Thresholds: p95 < 500ms, error rate < 1%, throughput > X rps

   **Stress Test** (2x expected traffic):
   - Ramp from 0 to 2x expected concurrent users over 3 minutes
   - Hold at target for 5 minutes
   - Ramp down over 2 minutes
   - Thresholds: p95 < 1s, error rate < 5%

   **Spike Test** (sudden burst):
   - Start at normal load for 2 minutes
   - Spike to 5x-10x in 10 seconds
   - Hold spike for 1 minute
   - Drop to normal for 2 minutes
   - Thresholds: recovery time < 30s, no crashes

5. **Define** thresholds per scenario:
   - p95 response time (ms)
   - Error rate (%)
   - Throughput (requests/sec)
   - Custom checks (response body validation, status codes)

6. **Add** execution instructions and CI integration config

## Threshold Reference

```
| Metric          | Smoke    | Load     | Stress   | Spike    |
|-----------------|----------|----------|----------|----------|
| p95 latency     | < 2000ms | < 500ms  | < 1000ms | < 3000ms |
| Error rate      | 0%       | < 1%     | < 5%     | < 10%    |
| Throughput      | baseline | > target | > target | recovery |
| p99 latency     | < 3000ms | < 1000ms | < 2000ms | < 5000ms |
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
## Performance Test Suite: [Service Name]

### Tool: [detected tool]
### Target Endpoints: [list]

### Files Generated
- `perf/[service]-smoke.{ext}` — Smoke test (1 VU, validation)
- `perf/[service]-load.{ext}` — Load test (expected traffic)
- `perf/[service]-stress.{ext}` — Stress test (2x capacity)
- `perf/[service]-spike.{ext}` — Spike test (sudden burst)
- `perf/test-data.csv` — Generated test data

### How to Run
1. Smoke: `[command to run smoke test]`
2. Load: `[command to run load test]`
3. Stress: `[command to run stress test]`
4. Spike: `[command to run spike test]`

### Thresholds
| Scenario | p95 Latency | Error Rate | Throughput |
|----------|-------------|------------|------------|
| Smoke    | < 2000ms    | 0%         | baseline   |
| Load     | < 500ms     | < 1%       | > X rps    |
| Stress   | < 1000ms    | < 5%       | > X rps    |
| Spike    | < 3000ms    | < 10%      | recovery   |

### CI Integration
[Config snippet for running in CI pipeline]
```
