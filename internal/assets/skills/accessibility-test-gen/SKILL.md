---
name: accessibility-test-gen
description: >
  Generates accessibility audit tests targeting WCAG 2.1 AA compliance.
  Integrates axe-core with the project's existing test framework to
  automate checks for color contrast, keyboard navigation, ARIA attributes,
  screen reader compatibility, and focus management.
  Trigger: When accessibility testing, WCAG compliance, or a11y audit is needed.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: testing
---

# Accessibility Test Generator

## When to Use

- New UI pages or components are added
- Accessibility audit or WCAG compliance is required
- User asks for "a11y test", "accessibility check", "WCAG audit", or "screen reader test"
- Pre-release accessibility validation
- Redesign or theme changes that affect color, layout, or interaction patterns

## Critical Rules

1. Target **WCAG 2.1 Level AA** — this is the most widely adopted legal and industry standard
2. Use **axe-core** for automated checks — it is the de facto standard and integrates with all major frameworks
3. NEVER rely solely on automated checks — include a **manual testing checklist** for things automation cannot catch
4. Automated tools catch only ~30-40% of accessibility issues — always state this limitation
5. Use the project's existing test framework (Playwright, Cypress, etc.) — do not introduce a separate tool
6. Categorize findings by **impact**: Critical (blocks users), Major (significant barrier), Minor (enhancement)
7. Test with **real assistive technology patterns**, not just DOM checks
8. ALWAYS test keyboard navigation independently — axe-core does not fully cover it

## Workflow

1. **Detect** UI framework and test runner:
   - Playwright: use `@axe-core/playwright` integration
   - Cypress: use `cypress-axe` plugin
   - Puppeteer: use `@axe-core/puppeteer`
   - If no framework: scaffold Playwright with axe-core (lightest setup)
   - Check if axe-core dependency exists, add it if missing

2. **Identify** pages and components to test:
   - Scan route definitions or page components
   - Prioritize: landing page, login/signup, main navigation, forms, data tables, modals/dialogs
   - Include responsive variants if layout changes significantly

3. **Generate** axe-core automated tests:
   - Run axe scan on each page/route
   - Configure rules to target WCAG 2.1 AA (`runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag21aa'] }`)
   - Test with default axe rules PLUS custom checks for:
     - Color contrast (minimum 4.5:1 for normal text, 3:1 for large text)
     - Image alt text presence and quality
     - Form label association
     - ARIA attribute validity
     - Heading hierarchy (no skipped levels)
     - Link purpose (descriptive text, not "click here")
     - Language attribute on html element

4. **Generate** keyboard navigation tests:
   - Tab order follows visual layout (logical sequence)
   - All interactive elements are focusable via Tab
   - Focus is visible on every focusable element (focus indicator meets 3:1 contrast)
   - Escape closes modals/dropdowns and returns focus to trigger
   - Enter/Space activates buttons and links
   - Arrow keys work in custom widgets (tabs, menus, comboboxes)
   - No keyboard traps (user can always Tab away)
   - Skip navigation link is first focusable element

5. **Generate** screen reader compatibility checks:
   - Landmark regions: header, nav, main, footer
   - ARIA live regions for dynamic content updates
   - Form error messages announced on validation
   - Modal focus trap with screen reader announcement
   - Table headers associated with data cells
   - Decorative images have empty alt (`alt=""`)
   - Status messages use `role="status"` or `aria-live="polite"`

6. **Generate** focus management tests:
   - Focus moves to new content after navigation (SPA route changes)
   - Focus moves into modals when opened
   - Focus returns to trigger when modal closes
   - Focus is not lost after dynamic content updates
   - Autofocus is used sparingly and appropriately

7. **Generate** form accessibility tests:
   - Every input has a visible, associated label
   - Required fields are indicated programmatically (`aria-required`)
   - Error messages are associated with inputs (`aria-describedby`)
   - Error summary appears on form submission with links to fields
   - Input purpose is identified (`autocomplete` attributes)
   - Fieldset/legend groups related controls (radio buttons, checkboxes)

## Impact Classification

```
| Impact   | Criteria                                              | Examples                              |
|----------|-------------------------------------------------------|---------------------------------------|
| Critical | Completely blocks users from completing a task         | No keyboard access, missing form labels, keyboard trap |
| Major    | Significant barrier, workaround is difficult           | Poor contrast, missing alt text, no focus indicator |
| Minor    | Usable but not ideal, enhancement opportunity          | Heading hierarchy skip, redundant ARIA, missing skip nav |
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
## Accessibility Test Suite: [Application Name]

### WCAG 2.1 AA Coverage
| Principle | Guidelines Tested | Automated | Manual |
|-----------|-------------------|-----------|--------|
| Perceivable (1.x) | [count] | [count] | [count] |
| Operable (2.x) | [count] | [count] | [count] |
| Understandable (3.x) | [count] | [count] | [count] |
| Robust (4.x) | [count] | [count] | [count] |

### Files Generated
- `a11y/[app]-axe-scan.test.{ext}` — Automated axe-core scans per page
- `a11y/[app]-keyboard.test.{ext}` — Keyboard navigation tests
- `a11y/[app]-forms.test.{ext}` — Form accessibility tests
- `a11y/manual-checklist.md` — Manual testing checklist

### Findings Summary
| ID | Impact | WCAG | Description | Element | Page |
|----|--------|------|-------------|---------|------|
| A11Y-001 | Critical | 1.1.1 | [finding] | [selector] | [page] |

### Manual Testing Checklist
- [ ] Screen reader announces page title on navigation
- [ ] Screen reader reads form errors when they appear
- [ ] Content is readable at 200% zoom without horizontal scroll
- [ ] Animations respect prefers-reduced-motion
- [ ] Touch targets are at least 44x44 CSS pixels (mobile)
```
