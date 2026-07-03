---
name: cypress-scaffold
description: >
  Scaffolds a Cypress testing project with custom commands, fixtures, plugins,
  interceptors, configuration, and CI integration.
  Trigger: When setting up end-to-end or component testing with Cypress.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: frameworks
---

# Cypress Scaffold

## When to Use

- Setting up e2e or component testing with Cypress for a new project
- Migrating from another testing framework to Cypress
- User asks to "set up Cypress" or "create Cypress tests"

## Critical Rules

1. Use **custom commands** for reusable actions (login, API setup) — not helper functions
2. Use **cy.intercept()** for network stubbing — NEVER rely on real backend for unit/component tests
3. Fixtures hold **static test data** — dynamic data uses factories in `support/`
4. All selectors must use **data-cy** attributes — no CSS class or tag selectors
5. NEVER use `cy.wait(ms)` — use `cy.intercept()` + `cy.wait('@alias')` for network, assertions for DOM
6. Configure **retries** in `cypress.config.ts` for CI stability (runMode: 2, openMode: 0)
7. Separate **e2e** and **component** test configurations
8. Include **TypeScript** support and typed custom commands

## Scaffold Structure

```
cypress/
├── cypress.config.ts                  # Base config, e2e + component settings
├── e2e/
│   ├── auth/
│   │   ├── login.cy.ts
│   │   └── logout.cy.ts
│   └── [feature]/
│       └── [feature].cy.ts
├── component/
│   └── [Component].cy.ts             # Component tests (if applicable)
├── fixtures/
│   ├── users.json                     # Static test data
│   └── [resource].json
├── support/
│   ├── commands.ts                    # Custom commands (login, setup)
│   ├── e2e.ts                         # e2e support file (imports commands)
│   ├── component.ts                   # Component support file
│   └── index.d.ts                     # TypeScript declarations for commands
├── plugins/
│   └── index.ts                       # Task definitions (DB seed, file ops)
├── downloads/                         # Downloaded files (gitignored)
├── screenshots/                       # Failure screenshots (gitignored)
├── videos/                            # Test run videos (gitignored)
├── tsconfig.json
└── package.json
```

## Custom Commands

```typescript
// support/commands.ts
Cypress.Commands.add('login', (email: string, password: string) => {
  cy.session([email, password], () => {
    cy.request('POST', '/api/auth/login', { email, password }).then((res) => {
      window.localStorage.setItem('token', res.body.token);
    });
  });
});

Cypress.Commands.add('getByCy', (selector: string) => {
  return cy.get(`[data-cy="${selector}"]`);
});

// support/index.d.ts
declare namespace Cypress {
  interface Chainable {
    login(email: string, password: string): Chainable<void>;
    getByCy(selector: string): Chainable<JQuery<HTMLElement>>;
  }
}
```

## Intercept Pattern

```typescript
// Stub API responses for deterministic tests
cy.intercept('GET', '/api/users', { fixture: 'users.json' }).as('getUsers');
cy.visit('/users');
cy.wait('@getUsers');
cy.getByCy('user-list').should('have.length.greaterThan', 0);
```

## CI Integration

```yaml
# GitHub Actions
- name: Run Cypress tests
  uses: cypress-io/github-action@v6
  with:
    build: npm run build
    start: npm run start
    wait-on: 'http://localhost:3000'
    browser: chrome
    record: true
  env:
    CYPRESS_RECORD_KEY: ${{ secrets.CYPRESS_RECORD_KEY }}
```
