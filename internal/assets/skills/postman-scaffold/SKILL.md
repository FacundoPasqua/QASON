---
name: postman-scaffold
description: >
  Scaffolds a Postman/Newman API testing workspace with collections,
  environments, pre-request scripts, and test assertions.
  Trigger: When setting up API testing with Postman or Newman.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-automator
  phase: frameworks
---

# Postman Scaffold

## When to Use

- Setting up API testing for a new project
- Migrating manual Postman collections to automated Newman pipelines
- User asks to "create Postman tests" or "set up Newman"

## Critical Rules

1. Collections must be organized by resource/domain, not by method
2. Use environment variables for ALL URLs, tokens, and dynamic data
3. Pre-request scripts handle auth token generation/refresh
4. Every request must have at least one test assertion
5. Use collection variables for values shared across requests in a flow
6. Newman-ready: collections must work headless in CI/CD

## Scaffold Structure

```
api-tests/
├── collections/
│   ├── [service-name].postman_collection.json
│   └── ...
├── environments/
│   ├── local.postman_environment.json
│   ├── staging.postman_environment.json
│   └── production.postman_environment.json
├── data/
│   └── test-data.csv                    # For data-driven runs
├── scripts/
│   ├── pre-request/
│   │   ├── auth.js                      # Auth token generation
│   │   └── setup.js                     # Common setup
│   └── tests/
│       ├── common-assertions.js         # Reusable test scripts
│       └── schema-validation.js         # JSON schema validators
├── newman.config.js                     # Newman runner config
├── package.json                         # Newman + dependencies
└── README.md
```

## Collection Organization

```
Collection: [Service Name]
├── Auth
│   ├── Login
│   ├── Refresh Token
│   └── Logout
├── [Resource A]
│   ├── Create [Resource A]
│   ├── Get [Resource A]
│   ├── List [Resource A]
│   ├── Update [Resource A]
│   └── Delete [Resource A]
├── [Resource B]
│   └── ...
└── Workflows
    ├── Happy Path: [Complete user flow]
    └── Error Flow: [Error scenario]
```

## Test Assertions Per Request

Every request MUST include:

```javascript
// 1. Status code
pm.test("Status code is 200", () => {
    pm.response.to.have.status(200);
});

// 2. Response time
pm.test("Response time < 500ms", () => {
    pm.expect(pm.response.responseTime).to.be.below(500);
});

// 3. Schema validation (for JSON responses)
pm.test("Response matches schema", () => {
    const schema = { /* JSON Schema */ };
    pm.response.to.have.jsonSchema(schema);
});

// 4. Business logic (specific to the endpoint)
pm.test("Returns created resource with ID", () => {
    const json = pm.response.json();
    pm.expect(json.id).to.exist;
    pm.expect(json.name).to.eql(pm.variables.get("resourceName"));
});
```

## CI/CD Integration

```json
{
  "scripts": {
    "test": "newman run collections/*.json -e environments/staging.json --reporters cli,junit",
    "test:local": "newman run collections/*.json -e environments/local.json",
    "test:data": "newman run collections/*.json -e environments/staging.json -d data/test-data.csv"
  }
}
```
