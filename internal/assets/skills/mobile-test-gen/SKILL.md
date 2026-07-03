---
name: mobile-test-gen
description: >
  Generates mobile test suites for Appium and Detox frameworks with
  platform-aware selectors, gesture handling, and async operation support.
  Covers iOS and Android differences, device-specific behaviors, and
  mobile-specific test patterns.
  Trigger: When mobile app testing is needed, or user provides mobile screens/flows to test.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: testing
---

# Mobile Test Generator

## When to Use

- User provides mobile app screens, flows, or requirements to test
- New mobile feature needs automated test coverage
- Existing mobile tests need to be extended or migrated between frameworks
- Cross-platform testing (iOS + Android) requires platform-aware test generation

## Critical Rules

1. ALWAYS generate platform-aware selectors — accessibility IDs first, then platform-specific fallbacks
2. Never use XPath selectors in mobile tests — they are brittle and slow
3. EVERY test MUST handle async operations explicitly — no arbitrary sleeps, use waitFor/waitForElement
4. Account for platform differences: iOS and Android handle navigation, permissions, keyboards, and gestures differently
5. Include app lifecycle tests: backgrounding, foregrounding, low memory, orientation changes
6. Test on real device characteristics: slow network, no network, small screen, large text (accessibility)
7. Never assume element visibility — always wait for the element to be visible before interacting

## Selector Strategy (Priority Order)

```
PRIORITY 1 — Accessibility ID (cross-platform, stable)
  Appium:  driver.findElement(By.accessibilityId('login-button'))
  Detox:   element(by.id('login-button'))
  → ALWAYS prefer this. Add testID/accessibilityIdentifier to the app if missing.

PRIORITY 2 — Resource ID / Test ID (platform-specific, stable)
  Android: driver.findElement(By.id('com.app:id/login_button'))
  iOS:     predicateString: "name == 'loginButton'"
  → Use when accessibility ID is not available.

PRIORITY 3 — Text (fragile, locale-dependent)
  Appium:  driver.findElement(By.xpath("//*[@text='Login']"))
  Detox:   element(by.text('Login'))
  → Last resort. Breaks on localization. ALWAYS note the fragility.

NEVER USE — XPath with structural paths
  //android.widget.LinearLayout[2]/android.widget.Button[1]
  → Breaks on ANY layout change. Absolutely forbidden.
```

## Gesture Patterns

```
TAP
  Appium:  element.click()
  Detox:   element(by.id('btn')).tap()

LONG PRESS
  Appium:  new Actions(driver).press(element).waitAction(2000).release().perform()
  Detox:   element(by.id('item')).longPress(2000)

SWIPE / SCROLL
  Appium:  new Actions(driver).press({x, y}).moveTo({x, y}).release().perform()
  Detox:   element(by.id('list')).scroll(200, 'down')

PINCH / ZOOM
  Appium:  driver.executeScript('mobile: pinch', {scale: 0.5, ...})
  Detox:   element(by.id('map')).pinch(1.5)  // scale factor

PULL TO REFRESH
  Appium:  swipe from top to bottom on scrollable container
  Detox:   element(by.id('scrollView')).swipe('down', 'slow', 0.5)
```

## Workflow

1. **Identify** the target screens and user flows
2. **Map** each flow to test scenarios:
   - Happy path through the flow
   - Input validation on each form field
   - Navigation (forward, back, deep link)
   - Platform-specific behavior differences
3. **Generate** tests per scenario:
   - Setup: app state, test data, permissions
   - Actions: using the selector strategy above
   - Waits: explicit waits for async operations
   - Assertions: visual state, data state, navigation state
4. **Add** mobile-specific scenarios:
   - Orientation change (portrait ↔ landscape)
   - Background/foreground app lifecycle
   - Permission dialogs (camera, location, notifications)
   - Keyboard handling (show, hide, input types)
   - Network conditions (offline, slow 3G, WiFi)
5. **Include** device matrix considerations

## Mobile-Specific Test Categories

```
NAVIGATION
  ✓ Forward navigation through each screen
  ✓ Back button / back gesture behavior
  ✓ Deep link opens correct screen with correct state
  ✓ Tab bar / bottom navigation switches correctly
  ✓ Modal presentation and dismissal

FORMS & INPUT
  ✓ Keyboard appears for each input field
  ✓ Keyboard type matches input (email, numeric, password)
  ✓ Field validation on blur and submit
  ✓ Scroll to field when keyboard covers it
  ✓ Copy/paste works in text fields

PERMISSIONS
  ✓ Grant permission → feature works
  ✓ Deny permission → graceful fallback
  ✓ Revoke permission mid-session → app handles it
  ✓ "Don't ask again" state → guide user to settings

LIFECYCLE
  ✓ Background → foreground preserves state
  ✓ Low memory warning → no crash
  ✓ Process death → restore or re-auth gracefully
  ✓ Orientation change preserves form data and scroll position

NETWORK
  ✓ Offline → show cached data or meaningful error
  ✓ Request in flight + network lost → timeout + retry option
  ✓ Slow network → loading indicators shown, no duplicate submissions
```

## Validation (MANDATORY)

After writing test files, you MUST run them:

1. Detect the test runner from `package.json` scripts, `Makefile`, or framework defaults
2. Run ONLY the generated test files (not the full suite)
3. If tests fail: read the error, fix the test, re-run
4. Report: X passed, Y failed, Z skipped

Never deliver tests you haven't run. A test that doesn't execute is not a test.

## Output Format

```markdown
## Mobile Test Suite: [Feature/Flow Name]

### Target
| Field | Value |
|-------|-------|
| Platforms | [iOS, Android, or both] |
| Framework | [Appium / Detox] |
| Min OS versions | [iOS X+ / Android API Y+] |
| Device matrix | [list target devices] |

### Test Cases
| ID | Category | Description | Platform | Selectors Used | Priority |
|----|----------|-------------|----------|---------------|----------|
| MOB-001 | Happy Path | [description] | Both | [selector IDs] | High |
| MOB-002 | Gesture | [description] | iOS-specific | [selector IDs] | Medium |

### Platform Differences
| Scenario | iOS Behavior | Android Behavior | Test Implication |
|----------|-------------|-----------------|-----------------|
| [scenario] | [behavior] | [behavior] | [what to assert differently] |

### Async Wait Strategy
| Operation | Wait Condition | Timeout | Fallback |
|-----------|---------------|---------|----------|
| [API call] | [element visible / text changes] | [ms] | [retry / fail] |
```
