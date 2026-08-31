---
name: playwright-mcp-runtime
description: >
  Uses the Playwright MCP server (@playwright/mcp) to drive a real browser
  at design and debug time — captures accessibility snapshots, validates
  selectors before committing them to spec files, and debugs failures by
  replaying flows interactively. Trigger: whenever the agent is about to
  write or fix a Playwright test and the project has the Playwright MCP
  server installed.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: frameworks
---

# Playwright MCP Runtime

## When to Use

- You are about to **write a new Playwright test** and need real selectors from the live app — do not hand-author selectors from memory or from the HTML alone
- You are **debugging a failing test** and need to see the page state the test sees (DOM + accessibility tree)
- You are **reviewing a PR** that adds Playwright tests and want to verify the selectors resolve against the deployed preview
- You are **fixing a flaky test** and need to confirm whether the target element is actually stable, loaded, and reachable

Prerequisite: Playwright MCP is installed. QASON auto-installs it in the `acceleration` and `optimization` presets (Claude Code only today). If `qason doctor` shows `✗ playwright` under the MCP section, re-run `qason install --preset acceleration` or install manually with `claude mcp add -s user playwright -- npx -y @playwright/mcp@latest`.

## Critical Rules

1. **MCP is for exploration, not execution.** Production tests ALWAYS run through `npx playwright test`, not through MCP. MCP is how you *discover* selectors; `@playwright/test` is how you *assert* them in CI.
2. **Capture a snapshot BEFORE writing any `getByRole`/`getByTestId` call.** `browser_snapshot` returns the accessibility tree with stable refs — copy those refs directly into your POM rather than inventing a selector that "looks right".
3. **Prefer accessibility-first selectors.** The snapshot gives you roles, names, and refs. Use `page.getByRole('button', { name: 'Submit' })` over CSS selectors whenever the role is present.
4. **Never commit MCP session IDs or refs into test files.** Refs are session-scoped — they're valid for the duration of the MCP session, not across test runs. Translate them to role/name/testid selectors before writing the spec.
5. **Always close browser resources at the end of a design session.** Call `browser_close` so stale tabs do not leak into the next agent's session.

## Core Tools You Will Use

| Tool | When to call it |
|------|-----------------|
| `browser_navigate(url)` | First action — load the page under test |
| `browser_snapshot()` | Get the accessibility tree with role + name + ref for every interactive element. The single most valuable tool — call it liberally |
| `browser_click(ref)` | Drive the flow forward to the state the test needs to assert |
| `browser_fill(ref, text)` | Populate form fields |
| `browser_select_option(ref, values)` | Pick from native selects |
| `browser_wait_for(text or selector)` | Block until async UI settles before snapshotting |
| `browser_take_screenshot()` | Attach a visual to bug reports or PR review comments |
| `browser_console_messages()` | Capture errors the UI logged during the flow — report them alongside the test results |
| `browser_network_requests()` | See XHR/fetch traffic; useful for verifying API mocking and for contract tests |
| `browser_close()` | Release resources at end of session |

Tool names follow `@playwright/mcp` v1.x. If the project pins an older version, adapt the names from its README — the same patterns apply.

## Workflow: Authoring a new e2e test

1. **Navigate to the starting URL** with `browser_navigate`.
2. **Snapshot** the landing state. Find the elements your test needs to interact with. Note the `role`, `name`, and `ref` for each.
3. **Drive the flow**: click / fill / wait through the scenario the test covers.
4. **Snapshot at each assertion point.** Copy stable text / role / aria-label values into the `expect(...)` calls you will write.
5. **Open the existing POM** (`pages/*.page.ts`) and add selectors using **roles and testids**, not refs. Example:
   ```ts
   // ✅ Good — translated from snapshot into a stable selector
   readonly submit = this.page.getByRole('button', { name: 'Submit order' });

   // ❌ Bad — ref is session-scoped, will not work in CI
   readonly submit = this.page.locator('[data-mcp-ref="e4"]');
   ```
6. **Close the MCP session** with `browser_close`.
7. **Run the test you just wrote** with `npx playwright test <path>` to verify the selectors resolve against a fresh browser. If it fails, repeat from step 2 — the snapshot will tell you what changed.

## Workflow: Debugging a failing test

1. Read the failing assertion + the Playwright trace (`--trace on-first-retry` output).
2. Navigate to the URL where the test failed.
3. Drive the flow with MCP to the state just before the failing assertion.
4. Snapshot — what element / state is actually there? Compare against the test's expectation.
5. Common diagnoses:
   - **Selector drift**: the role/name/testid changed in the UI. Update the POM.
   - **Timing**: the element exists but is mid-transition. Add `browser_wait_for` / `page.waitForLoadState('networkidle')`.
   - **Environment**: the test ran against a different build. Confirm base URL and auth state.
   - **Flaky assertion**: the assertion depends on order-of-appearance. Use a more specific selector with `.first()` / `.nth()` justified by the snapshot.

## Workflow: reviewing existing Playwright tests

When a PR adds or modifies Playwright tests:
1. Check out the PR branch; deploy or point to the preview URL.
2. For each modified `.spec.ts`, navigate to the scenario start and snapshot.
3. Verify every `getByRole`/`getByTestId` in the spec resolves to exactly one element in the snapshot. Flag ambiguous selectors (multiple matches) in the review.
4. Confirm assertions reference values that are actually in the accessibility tree — don't trust comments, trust the snapshot.

## What Not to Use MCP For

- **Do not** run the full test suite through MCP — `npx playwright test` is vastly more efficient and parallelizable.
- **Do not** use MCP as a production test runner. Sessions are not reproducible; traces from the real Playwright runner are.
- **Do not** hand MCP output to a visual-regression step. Use Playwright's built-in `toHaveScreenshot()` for baselines.
- **Do not** use MCP in CI. It's a design/debug tool; CI runs the deterministic `@playwright/test` binary.

## Handoff

When you finish a design session, the output should contain:
- A list of **stable selectors** chosen (role + name pairs, or testids) with a one-line justification per selector based on the snapshot
- The **POM updates** written with those selectors
- The **spec file** that uses the POM
- The **command used to run the new test**, with the outcome (pass/fail + flakiness notes)

That handoff is what the next reader — human or agent — will rely on. Without it, the MCP session's value evaporates when the session closes.
