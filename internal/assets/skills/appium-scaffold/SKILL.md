---
name: appium-scaffold
description: >
  Scaffolds an Appium mobile testing project with capabilities config, page
  objects, platform helpers for iOS and Android, and gesture utilities.
  Trigger: When setting up mobile testing with Appium.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: frameworks
---

# Appium Scaffold

## When to Use

- Setting up mobile e2e testing for iOS and/or Android
- Migrating from manual mobile testing to automated Appium tests
- User asks to "set up Appium" or "create mobile tests"

## Critical Rules

1. Use **Page Object Model** — no raw selectors in test files
2. Capabilities must be **platform-specific** with shared base config
3. Use **accessibility IDs** as the primary locator strategy (cross-platform)
4. Gesture utilities must abstract platform differences (swipe, scroll, tap coordinates)
5. Tests must declare their **platform target** (iOS, Android, or both)
6. NEVER hardcode device names or OS versions — use capabilities profiles
7. Include **wait helpers** for element visibility — mobile rendering is slower than web
8. Separate **platform-specific** page objects only when behavior diverges

## Scaffold Structure

```
mobile-tests/
├── config/
│   ├── capabilities.ts                # Shared + platform-specific caps
│   ├── android.caps.ts                # Android device profiles
│   ├── ios.caps.ts                    # iOS device/simulator profiles
│   └── appium.server.ts              # Server connection config
├── pages/
│   ├── base.page.ts                   # Abstract base (find, wait, tap helpers)
│   ├── login.page.ts                  # Cross-platform page object
│   ├── android/                       # Android-only page overrides
│   │   └── [feature].page.ts
│   └── ios/                           # iOS-only page overrides
│       └── [feature].page.ts
├── tests/
│   ├── auth/
│   │   └── login.spec.ts
│   └── [feature]/
│       └── [feature].spec.ts
├── helpers/
│   ├── gestures.ts                    # Swipe, scroll, pinch, long press
│   ├── platform.ts                    # Platform detection and branching
│   └── wait.ts                        # Explicit wait utilities
├── data/
│   ├── test-users.ts                  # Test account credentials
│   └── [resource].data.ts
├── wdio.conf.ts                       # WebdriverIO + Appium config
├── package.json
└── tsconfig.json
```

## Capabilities Config

```typescript
// config/capabilities.ts
const sharedCaps = {
  'appium:automationName': 'UiAutomator2', // Android
  'appium:newCommandTimeout': 240,
  'appium:noReset': false,
};

export const androidCaps = {
  ...sharedCaps,
  platformName: 'Android',
  'appium:automationName': 'UiAutomator2',
  'appium:deviceName': process.env.ANDROID_DEVICE || 'Pixel_7_API_34',
  'appium:app': process.env.ANDROID_APP_PATH || './apps/app-debug.apk',
};

export const iosCaps = {
  ...sharedCaps,
  platformName: 'iOS',
  'appium:automationName': 'XCUITest',
  'appium:deviceName': process.env.IOS_DEVICE || 'iPhone 15',
  'appium:platformVersion': process.env.IOS_VERSION || '17.0',
  'appium:app': process.env.IOS_APP_PATH || './apps/App.app',
};
```

## Page Object Pattern

```typescript
// pages/base.page.ts
export abstract class BasePage {
  constructor(protected driver: WebdriverIO.Browser) {}

  protected async findByAccessibilityId(id: string) {
    return this.driver.$(`~${id}`);
  }

  protected async tapElement(id: string): Promise<void> {
    const el = await this.findByAccessibilityId(id);
    await el.waitForDisplayed({ timeout: 10000 });
    await el.click();
  }

  protected async setText(id: string, text: string): Promise<void> {
    const el = await this.findByAccessibilityId(id);
    await el.waitForDisplayed({ timeout: 10000 });
    await el.setValue(text);
  }
}
```

## Gesture Utilities

```typescript
// helpers/gestures.ts
export async function swipeUp(driver: WebdriverIO.Browser): Promise<void> {
  const { width, height } = await driver.getWindowSize();
  await driver.touchAction([
    { action: 'press', x: width / 2, y: height * 0.8 },
    { action: 'wait', ms: 200 },
    { action: 'moveTo', x: width / 2, y: height * 0.2 },
    { action: 'release' },
  ]);
}

export async function scrollToElement(
  driver: WebdriverIO.Browser,
  accessibilityId: string,
  maxScrolls: number = 5,
): Promise<void> {
  for (let i = 0; i < maxScrolls; i++) {
    const el = await driver.$(`~${accessibilityId}`);
    if (await el.isDisplayed()) return;
    await swipeUp(driver);
  }
  throw new Error(`Element ~${accessibilityId} not found after ${maxScrolls} scrolls`);
}
```
