---
name: playwright-scaffold
description: >
  Scaffolds a Playwright testing project with Page Object Model, configuration,
  fixtures, custom reporters, and CI integration.
  Trigger: When setting up end-to-end testing with Playwright.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: frameworks
---

# Playwright Scaffold

## When to Use

- Setting up e2e testing for a new web application
- Migrating from another e2e framework to Playwright
- User asks to "set up Playwright" or "create e2e tests"

## Critical Rules

1. Always use **Page Object Model** — no raw selectors in test files
2. Use **data-testid** as the primary selector strategy
3. Configure **multiple browsers** (chromium, firefox, webkit) in `playwright.config.ts`
4. Use **fixtures** for authentication state and shared setup — not `beforeAll` hacks
5. Custom reporters must output both **human-readable** and **CI-parseable** formats (JUnit XML)
6. All waits must be **explicit** (waitForSelector, waitForResponse) — NEVER use `waitForTimeout`
7. Tests must be **parallelizable** — no shared mutable state between test files
8. Include **visual regression** setup with screenshot comparison baseline

## Scaffold Structure

```
e2e/
├── playwright.config.ts              # Multi-browser, base URL, timeouts, reporter
├── fixtures/
│   ├── base.fixture.ts               # Extended test with custom fixtures
│   ├── auth.fixture.ts               # Authenticated state fixture
│   └── test-data.fixture.ts          # Test data factory fixture
├── pages/
│   ├── base.page.ts                  # Abstract base page (common helpers)
│   ├── login.page.ts                 # Login page object
│   └── [feature].page.ts             # One page object per page/component
├── tests/
│   ├── auth/
│   │   ├── login.spec.ts
│   │   └── logout.spec.ts
│   └── [feature]/
│       └── [feature].spec.ts
├── helpers/
│   ├── api.helper.ts                 # API shortcuts for test setup/teardown
│   └── data.helper.ts                # Test data generators
├── reporters/
│   └── custom-reporter.ts            # Custom result formatter
├── .auth/                            # Stored auth state (gitignored)
│   └── user.json
├── package.json
└── tsconfig.json
```

## Page Object Pattern

```typescript
// pages/base.page.ts
import { Page, Locator } from '@playwright/test';

export abstract class BasePage {
  constructor(protected readonly page: Page) {}

  protected getByTestId(testId: string): Locator {
    return this.page.getByTestId(testId);
  }

  async waitForPageLoad(): Promise<void> {
    await this.page.waitForLoadState('networkidle');
  }
}

// pages/login.page.ts
import { BasePage } from './base.page';

export class LoginPage extends BasePage {
  readonly emailInput = this.getByTestId('email-input');
  readonly passwordInput = this.getByTestId('password-input');
  readonly submitButton = this.getByTestId('login-submit');

  async goto(): Promise<void> {
    await this.page.goto('/login');
  }

  async login(email: string, password: string): Promise<void> {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.submitButton.click();
  }
}
```

## Fixture Pattern

```typescript
// fixtures/base.fixture.ts
import { test as base } from '@playwright/test';
import { LoginPage } from '../pages/login.page';

type Fixtures = {
  loginPage: LoginPage;
};

export const test = base.extend<Fixtures>({
  loginPage: async ({ page }, use) => {
    await use(new LoginPage(page));
  },
});

export { expect } from '@playwright/test';
```

## CI Integration

```yaml
# GitHub Actions
- name: Run Playwright tests
  run: npx playwright test
- name: Upload report
  if: always()
  uses: actions/upload-artifact@v4
  with:
    name: playwright-report
    path: playwright-report/
```

## Playwright Test CLI — Flags That Actually Matter

Authoring tests without knowing these flags is expensive: you run 3-minute suites to validate a one-line selector change, debug flakiness without a trace, or push a broken shard to CI. Reference: <https://playwright.dev/docs/test-cli>.

### Authoring loop (run what you're writing, not the whole suite)

| Flag | When to use it |
|------|----------------|
| `--grep "<title>"` | Run only tests whose title matches the regex — fastest way to iterate on a single spec |
| `--grep-invert "<title>"` | Inverse of above; useful for skipping a known-broken test while you fix another |
| `<path>` (positional) | Pass the spec path directly: `npx playwright test tests/checkout.spec.ts` |
| `--headed` | Render the browser visibly; essential when reproducing a flaky flow by eye |
| `--debug` | Open Playwright Inspector — pauses on `page.pause()` and lets you step through actions |
| `--ui` | Open the Playwright UI mode (interactive time-travel, watch mode) — the single best debugging tool |
| `--project=<name>` | Restrict to one browser project from `playwright.config.ts` (e.g. `--project=chromium`) |

### Diagnostic loop (when tests pass locally and fail in CI, or vice versa)

| Flag | When to use it |
|------|----------------|
| `--trace=on-first-retry` | Record a full trace on the first retry of any failing test — the default we recommend for CI |
| `--trace=on` | Trace every test (expensive, use when debugging a specific flaky test) |
| `--trace=retain-on-failure` | Keep trace only for failed tests — good compromise for nightly runs |
| `--reporter=list` | Per-test live output; easier to read than the default dot reporter when you're triaging |
| `--reporter=html` | Produce the browsable HTML report under `playwright-report/` |
| `--reporter=junit` | JUnit XML for CI dashboards (gitlab, jenkins, ADO test tab) |

View a trace with `npx playwright show-trace path/to/trace.zip` — it's a time-travel debugger, not a log. Teach the team to open traces before they open issues.

### Flakiness & efficiency loop (CI discipline)

| Flag | When to use it |
|------|----------------|
| `--retries=<n>` | Retry failing tests up to n times; use 2 in CI, 0 locally so you actually see flakiness |
| `--workers=<n>` | Parallelism cap; CI usually auto-detects, but pin it on shared runners to avoid noisy-neighbor flakiness |
| `--shard=<i>/<n>` | Split the suite across n CI jobs; each runs shard i. Keeps wall-clock time flat as the suite grows |
| `--last-failed` | Re-run only the tests that failed in the previous run — 10× faster feedback while fixing a regression cluster |
| `--only-changed` | Run only specs affected by the git diff (compared to `main` by default). Use in pre-push hooks |
| `--update-snapshots` | Regenerate visual / text snapshots. Flag aware: also updates `toHaveScreenshot()` baselines |
| `--timeout=<ms>` | Override per-test timeout; bump when legitimately slow, don't bump to paper over flakiness |

### Rules of thumb for the qa-automator

1. **After writing a test, run it with `--grep "<your test title>" --headed --reporter=list`.** Confirm it passes AND fails when you intentionally break an assertion.
2. **In CI, always run with `--trace=on-first-retry --retries=2 --reporter=html --reporter=junit`.** The trace pays for itself the first time it explains a flaky failure in 30 seconds.
3. **While fixing a failing suite, run `--last-failed` in a tight loop until green.** Then re-run the full suite once to confirm no regressions.
4. **Never land a test without running `--only-changed` locally.** If you broke something adjacent, you'll see it before the CI job does.
5. **When updating visuals, run `--update-snapshots --project=chromium` per project explicitly.** A bare `--update-snapshots` with multi-browser config will blow away baselines you didn't mean to touch.
