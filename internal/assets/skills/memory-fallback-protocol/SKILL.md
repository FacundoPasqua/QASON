---
name: memory-fallback-protocol
description: >
  Canonical protocol for persisting QA project context when the engram MCP server
  is NOT available. Defines the on-disk layout, read/write steps, and topic-key
  schema the other six memory skills fall back to. Activated by every memory skill
  as the last resort so agents always have a place to save discoveries.
  Trigger: Whenever a memory skill wants to save/load context and engram is unavailable.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-analyst, qa-automator, qa-test-designer, qa-reviewer, qa-ops
  phase: memory-infrastructure
  requires_mcp: none
---

# Memory Fallback Protocol

## Why this exists

QASON treats persistent memory as a first-class capability: without it, every
QA agent starts cold, re-discovers the same conventions, and burns tokens on
work it already did last week. The preferred backend is **engram** (MCP server
with semantic search, auto-compaction, cross-session survival).

But engram is not always installed. Rather than fail silently — the old
"optional engram" behaviour, which made memory feel unreliable — QASON falls
back to a **file-based** store that works on any machine, no MCP required.

The three layers:

| Layer | Backend | When it kicks in |
|-------|---------|------------------|
| 1 | engram MCP | Always preferred if the `mem_*` tools are available |
| 2 | File store (this skill) | Whenever engram is NOT available |
| 3 | Install gate | `qason install --preset acceleration/optimization` aborts if neither is reachable — memory is mandatory at scale |

## Detection

At the start of every memory task, check whether `mem_context` (or any
`mem_*` tool) appears in your available tools. If present, use engram. If
absent, use the file store fallback — do not attempt the call: when engram
is not installed the tool is simply not in the tool list, so there is
nothing to "try" or time out on.

NEVER tell the user engram is missing unless they explicitly ask — the
fallback is designed to be invisible.

## File Layout

All memory lives under the QASON state directory in the user's home:

```
~/.qason/memory/
└── qason/
    └── project/
        └── {project-name}/
            ├── stack.md
            ├── conventions.md
            ├── domain.md
            ├── auth-flow.md
            ├── bug-patterns.md
            ├── flaky-tests.md
            ├── test-data.md
            └── coverage-map.md
```

- `~/.qason/` is the same directory used by `state.json` — we reuse the same
  root so uninstall can clean it up in one sweep.
- `{project-name}` is the detected project slug:
  1. Try `git remote get-url origin` and take the repo basename (lowercase,
     strip `.git`).
  2. Fallback: basename of the working directory (lowercase).
- Each topic key maps to ONE file. A topic key of the form
  `qason/project/myapp/stack` becomes `~/.qason/memory/qason/project/myapp/stack.md`.

## File Format

Every memory file is markdown with a small YAML frontmatter header so we can
roundtrip metadata without parsing the body:

```markdown
---
topic_key: qason/project/myapp/stack
project: myapp
updated_at: 2026-04-17T19:42:11Z
version: 1
---

# stack

<free-form content the agent saved — typically a short bulleted summary>
```

Rules:

- The `topic_key` and `project` fields MUST match the path, so a stray file
  is self-describing.
- `updated_at` is RFC3339 UTC. Refresh it on every write.
- `version` starts at 1 and increments on every overwrite. Old versions are
  NOT kept — this is a scratchpad, not a git log.
- The body is free-form. Agents save a short summary (bullets, tables) —
  aim for 100–500 words per file.

## Operations

### Save (equivalent to `mem_save`)

```
1. Compute path = "~/.qason/memory/" + topic_key + ".md"
2. Ensure parent directory exists (mkdir -p).
3. If the file exists, read it to preserve the `version` counter and
   increment by 1. Otherwise start at version 1.
4. Write frontmatter + content as shown above, overwriting the whole file.
```

Use the `Write` tool. Do NOT use `Edit` — memory files are overwritten wholesale.

### Load (equivalent to `mem_search` + `mem_get_observation`)

```
1. Compute path as above.
2. Use the `Read` tool on that path.
3. If the file does not exist, return "no memory for this topic key".
4. If it exists, return the body (everything after the closing `---`).
```

For broader searches (e.g. "all memory for project X"):

```
1. Use the `Glob` tool with pattern "~/.qason/memory/qason/project/{name}/**/*.md"
2. Read each match and aggregate results.
```

### Context summary (equivalent to `mem_context`)

```
1. Detect project name.
2. Glob "~/.qason/memory/qason/project/{name}/*.md".
3. Read each file, keep only the body (skip frontmatter).
4. Concatenate with `## {topic}` headers derived from the filename.
```

## Topic-Key Schema

Same schema as the engram-backed skills — memory files mirror it 1:1:

| Topic key | Contains |
|-----------|----------|
| `qason/project/{name}/stack` | Framework, versions, test tools, build system |
| `qason/project/{name}/conventions` | Naming, folder structure, test patterns |
| `qason/project/{name}/domain` | Business rules, glossary, critical flows |
| `qason/project/{name}/auth-flow` | Authentication model, tokens, roles |
| `qason/project/{name}/bug-patterns` | Recurring bug types, hot spots |
| `qason/project/{name}/flaky-tests` | Known flaky tests + root causes |
| `qason/project/{name}/test-data` | Test users, fixtures, seed data |
| `qason/project/{name}/coverage-map` | What is tested, what is not |

## Concurrency & Consistency

- Writes are atomic because Go/Node/Python `os.WriteFile` primitives are
  atomic on the platforms we target (the host agent chooses). QASON does
  NOT guarantee atomicity across parallel agent runs on the same project.
- If two agents save to the same topic key concurrently, last write wins.
  This is acceptable: memory is a cache of recent discoveries, not a
  transactional database.
- Use distinct project names if you run two agents on the same codebase
  with different contexts.

## When to use engram vs. files

| Situation | Backend |
|-----------|---------|
| engram MCP is installed and responsive | engram |
| engram is installed but timing out / erroring | engram (let it fail; do NOT quietly fall through mid-task) |
| engram tool is not in your tool list | files |
| User explicitly asked for "offline mode" | files |

The install gate ensures that for presets `acceleration` and above, at least
one of the two backends is functioning before QASON considers itself
installed. Skills SHOULD assume a backend exists and fail loudly if both
are gone — that signals a broken install.
