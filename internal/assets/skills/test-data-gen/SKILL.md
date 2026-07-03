---
name: test-data-gen
description: >
  Generates test data factories, fixtures, and builders from data models.
  Analyzes TypeScript interfaces, Go structs, Python classes, and DB schemas
  to produce realistic, referentially consistent test data with edge cases
  and boundary values.
  Trigger: When test data generation, factories, fixtures, or data builders are needed.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-test-designer
  phase: development
---

# Test Data Generator

## When to Use

- New data models or entities are added to the project
- Test setup is verbose with repetitive object construction
- Tests need realistic data instead of trivial placeholders
- User asks for "test data", "factory", "fixture", "builder", "fake data", or "seed data"
- Load tests need bulk data generation
- Tests require referentially consistent data across related entities

## Critical Rules

1. Analyze **actual data models** in the codebase — do not guess field names or types
2. Generate **realistic data**, not random gibberish (use locale-aware fakers)
3. Maintain **referential integrity** — foreign keys must point to valid parent records
4. Include **edge case values** for every field type (empty, null, max length, unicode, special chars)
5. NEVER use the project's existing test framework for anything else — detect and use the data generation library the project already has
6. If no data library exists, choose based on language:
   - JS/TS: `@faker-js/faker` + `fishery` (factories)
   - Python: `factory_boy` + `faker`
   - Go: `go-faker` or `gofakeit` + custom builders
   - Java: `java-faker` + custom builders
7. Factories produce **valid data by default** — invalid data must be an explicit override
8. Support **overrides** for every field — tests must be able to customize specific values
9. Generate a **boundary value set** for each field type — these catch off-by-one and validation bugs

## Workflow

1. **Analyze** data models in the project:
   - TypeScript: scan for `interface`, `type`, and class definitions
   - Go: scan for `struct` definitions with field tags
   - Python: scan for dataclass, Pydantic model, SQLAlchemy model, Django model definitions
   - Java: scan for entity classes with JPA annotations
   - Database: scan migration files or schema definitions
   - Map relationships: one-to-many, many-to-many, foreign keys

2. **Detect** existing data generation tools:
   - Check `package.json` for `@faker-js/faker`, `fishery`, `factory.ts`, `chance`
   - Check `go.mod` for `gofakeit`, `go-faker`
   - Check `requirements.txt`/`pyproject.toml` for `factory_boy`, `faker`, `model_bakery`
   - Check `pom.xml`/`build.gradle` for `java-faker`, `easy-random`
   - If none found: recommend and install the appropriate library

3. **Generate** factories (dynamic data, faker-based):
   - One factory per entity/model
   - Default values produce a **valid, complete** instance
   - Use faker for realistic values:
     ```
     name        → faker.person.fullName()
     email       → faker.internet.email()
     phone       → faker.phone.number()
     address     → faker.location.streetAddress()
     date        → faker.date.recent()
     url         → faker.internet.url()
     uuid        → faker.string.uuid()
     description → faker.lorem.paragraph()
     price       → faker.number.float({ min: 0.01, max: 9999.99, fractionDigits: 2 })
     ```
   - Support overrides:
     ```
     UserFactory.build()                        // valid defaults
     UserFactory.build({ email: 'specific@test.com' }) // override one field
     UserFactory.buildList(10)                   // bulk generation
     UserFactory.build({ role: 'admin' })        // specific variant
     ```
   - Include **traits/variants** for common scenarios:
     ```
     UserFactory.traits.admin     // user with admin role
     UserFactory.traits.inactive  // deactivated user
     UserFactory.traits.withPosts // user with related posts
     ```

4. **Generate** fixtures (static JSON/YAML data):
   - Minimal fixture set: one valid instance per entity
   - Named fixtures for specific test scenarios:
     ```
     fixtures/users/valid-user.json
     fixtures/users/admin-user.json
     fixtures/users/user-with-posts.json
     fixtures/orders/empty-cart.json
     fixtures/orders/completed-order.json
     ```
   - Maintain referential integrity across fixture files
   - Include a fixture loader utility if the project does not have one

5. **Generate** builders (fluent API pattern):
   - For languages/projects that prefer builder pattern:
     ```
     new UserBuilder()
       .withName("Test User")
       .withEmail("test@example.com")
       .withRole("admin")
       .withPosts(3)
       .build()
     ```
   - Builder validates required fields at build time
   - Builder supports nested builders for related entities

6. **Generate** boundary value sets per field type:

   ```
   | Type    | Boundary Values                                                    |
   |---------|--------------------------------------------------------------------|
   | string  | "", " ", null, "a", max_length, max_length+1, unicode: "日本語",   |
   |         | emoji: "🎉", XSS: "<script>", SQL: "'; DROP", newlines: "a\nb"    |
   | integer | 0, 1, -1, MIN_INT, MAX_INT, null                                  |
   | float   | 0.0, 0.01, -0.01, MAX_FLOAT, NaN, Infinity, null                  |
   | boolean | true, false, null                                                  |
   | date    | today, yesterday, epoch (1970-01-01), far future (2099-12-31),     |
   |         | leap day (2024-02-29), null                                        |
   | email   | valid@test.com, a@b.co, very+long@domain.com, @missing-local,     |
   |         | missing-domain@, null                                              |
   | array   | [], [one], [max_items], [max_items+1], null, [duplicates]          |
   | enum    | each valid value, invalid value, null, empty string                |
   | uuid    | valid v4, all zeros, invalid format, null                          |
   ```

7. **Generate** relationship-aware data:
   - Parent entities are created before children
   - Foreign keys reference existing parent IDs
   - Cascading creation: `OrderFactory.build()` automatically creates a valid `User` parent
   - Provide utilities for creating complete entity graphs:
     ```
     createScenario('user-with-orders', {
       user: 1,
       orders: 3,
       orderItems: { perOrder: 2 }
     })
     ```

8. **Generate** bulk data utilities for load testing:
   - Provide a script/function to generate N records with unique data
   - Include CSV/JSON export for external tools (k6, JMeter, Locust)
   - Ensure uniqueness constraints are respected (unique emails, usernames)
   - Support deterministic generation (same seed = same data) for reproducibility

## Output Template

```markdown
## Test Data Suite: [Application Name]

### Entities Analyzed
| Entity | Fields | Relationships | Factory | Fixture | Builder |
|--------|--------|---------------|---------|---------|---------|
| User | 8 | has_many: Post, Order | ✅ | ✅ | ✅ |
| Post | 5 | belongs_to: User | ✅ | ✅ | — |

### Files Generated
- `test-data/factories/[entity].factory.{ext}` — Dynamic factories per entity
- `test-data/fixtures/[entity]/` — Static fixture files (JSON/YAML)
- `test-data/builders/[entity].builder.{ext}` — Fluent builders (if applicable)
- `test-data/boundaries/[entity].boundaries.{ext}` — Boundary value sets
- `test-data/scenarios.{ext}` — Relationship-aware scenario builders
- `test-data/bulk-generator.{ext}` — Bulk data generation for load tests

### Usage Examples
| Scenario | Code |
|----------|------|
| Single valid entity | `UserFactory.build()` |
| With specific override | `UserFactory.build({ role: 'admin' })` |
| Bulk generation | `UserFactory.buildList(100)` |
| With related entities | `UserFactory.traits.withPosts.build()` |
| Boundary testing | `UserBoundaries.string.empty` |
```
