---
name: visual-regression-gen
description: >
  Generates visual regression test setup with screenshot comparison.
  Detects the UI framework and test runner, configures baseline management,
  and generates tests for key pages across responsive breakpoints.
  Trigger: When visual regression testing, screenshot comparison, or UI consistency checks are needed.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: testing
---

# Visual Regression Test Generator

## When to Use

- UI redesign or theme changes need validation
- Component library updates need visual consistency checks
- CSS refactoring or framework migration
- User asks for "visual test", "screenshot test", "visual regression", or "UI diff"
- Pre-release visual validation across browsers or breakpoints
- Dark mode or theming support is added

## Critical Rules

1. ALWAYS detect the existing test runner before choosing a visual testing approach
2. Playwright's built-in `toHaveScreenshot()` is the **preferred approach** — it requires no extra dependencies
3. If Cypress: suggest `cypress-image-diff-js` (open source, actively maintained) or Percy (cloud-based)
4. If no test framework exists: scaffold Playwright visual tests as the default
5. Set threshold tolerance to **0.1%** by default — this catches real changes while ignoring anti-aliasing noise
6. NEVER commit baseline screenshots to the repository without explicit user confirmation
7. Generate tests for **at least 3 breakpoints**: mobile (375px), tablet (768px), desktop (1280px)
8. Disable animations and transitions in visual tests — they cause flaky comparisons
9. Use consistent viewport sizes and font rendering settings across test runs
10. Include CI configuration — visual tests are most valuable in automated pipelines

## Workflow

1. **Detect** UI framework and test runner:
   - Check `package.json` for Playwright, Cypress, Puppeteer, or Storybook
   - Check for existing visual test configuration files
   - If Playwright found: use built-in `toHaveScreenshot()` API
   - If Cypress found: check for `cypress-image-diff-js` or `percy`; install if missing
   - If Storybook found: suggest Chromatic or Storybook test-runner with screenshots
   - If nothing found: scaffold Playwright project with visual test config

2. **Identify** pages and components to capture:
   - Scan route definitions for key pages
   - Prioritize: home/landing, login, dashboard, settings, forms, data tables
   - Identify component variants if using Storybook
   - Ask user to confirm the critical visual surfaces

3. **Configure** visual test infrastructure:
   - Set viewport sizes for each breakpoint:
     ```
     mobile:  { width: 375, height: 812 }
     tablet:  { width: 768, height: 1024 }
     desktop: { width: 1280, height: 720 }
     ```
   - Disable animations globally:
     ```css
     *, *::before, *::after {
       animation-duration: 0s !important;
       transition-duration: 0s !important;
     }
     ```
   - Configure consistent font rendering (if cross-platform CI)
   - Set screenshot comparison threshold (0.1% default)
   - Configure baseline directory structure

4. **Generate** full-page screenshot tests:
   - One test per page per breakpoint
   - Wait for network idle and fonts loaded before capture
   - Hide dynamic content (timestamps, random data) with CSS or element masking
   - Name screenshots consistently: `[page]-[breakpoint].png`

5. **Generate** component-level screenshot tests (if applicable):
   - Capture individual components in isolation
   - Test component states: default, hover, focus, active, disabled, error, loading
   - Test with different data: empty, single item, full, overflow

6. **Generate** theme variant tests (if applicable):
   - Light mode vs dark mode
   - High contrast mode
   - Custom theme variations
   - Test each theme at desktop breakpoint minimum

7. **Configure** baseline management:
   - Generate initial baselines with a dedicated script/command
   - Document the baseline update workflow:
     - Run tests → review diffs → approve or reject → update baselines
   - Add `.gitattributes` for screenshot binary handling (Git LFS if repo is large)
   - Create a `.gitignore` entry for diff output images

8. **Add** CI integration:
   - Generate CI config snippet (GitHub Actions, GitLab CI)
   - Upload visual diff artifacts on failure
   - Use Docker or consistent OS image for deterministic rendering
   - Add a step to compare against main branch baselines

## Breakpoint Reference

```
| Breakpoint | Width  | Height | Device Example       |
|------------|--------|--------|----------------------|
| mobile-sm  | 320px  | 568px  | iPhone SE            |
| mobile     | 375px  | 812px  | iPhone 13            |
| tablet     | 768px  | 1024px | iPad                 |
| desktop    | 1280px | 720px  | Standard laptop      |
| desktop-lg | 1920px | 1080px | Full HD monitor      |
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
## Visual Regression Test Suite: [Application Name]

### Tool: [Playwright toHaveScreenshot / cypress-image-diff-js / Percy]
### Breakpoints: mobile (375px), tablet (768px), desktop (1280px)

### Files Generated
- `visual/[app]-pages.visual.{ext}` — Full-page screenshot tests
- `visual/[app]-components.visual.{ext}` — Component screenshot tests (if applicable)
- `visual/config.{ext}` — Visual test configuration
- `visual/baselines/` — Baseline screenshot directory
- `.github/workflows/visual-tests.yml` — CI pipeline config (if GitHub)

### Pages Captured
| Page | Mobile | Tablet | Desktop | Baseline |
|------|--------|--------|---------|----------|
| Home | ⬜ | ⬜ | ⬜ | ⬜ |
| Login | ⬜ | ⬜ | ⬜ | ⬜ |

### Baseline Management
1. Generate baselines: `[command]`
2. Run comparison: `[command]`
3. Update baselines after approved changes: `[command]`
4. Review diffs: [location of diff output]

### Configuration
- Threshold: 0.1% (pixel difference tolerance)
- Animations: disabled via global CSS injection
- Dynamic content: masked elements listed in config
```
