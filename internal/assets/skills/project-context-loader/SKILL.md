---
name: project-context-loader
description: >
  Loads QA-specific project context from persistent memory at the start of any task.
  Retrieves stack info, conventions, domain knowledge, auth flows, and historical gotchas
  so the agent does not re-discover what is already known.
  Trigger: At the start of ANY QA task, before analyzing requirements, generating tests, or triaging bugs.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-analyst, qa-automator, qa-test-designer, qa-reviewer, qa-ops
  phase: pre-task
  requires_mcp: engram
---

# Project Context Loader

## When to Use

- At the START of any QA task (before doing any work)
- When user mentions a project or repo you may have worked on before
- After a context compaction or new session begins
- Before making recommendations that depend on project conventions

## Critical Rules

1. ALWAYS run this BEFORE any other skill on a new task
2. If engram MCP is NOT available, skip silently and continue — do NOT block the task
3. Use the current project name (detect from `git remote -v`, package.json, or working directory basename)
4. Use these topic keys — do NOT invent new ones
5. If memory is empty for this project, flag it as a "cold start" and enable auto-discovery mode

## Topic Keys to Load

| Topic Key | Contains |
|-----------|----------|
| `qason/project/{name}/stack` | Framework, versions, test tools, build system |
| `qason/project/{name}/conventions` | Naming, folder structure, test patterns |
| `qason/project/{name}/domain` | Business rules, glossary, critical flows |
| `qason/project/{name}/auth-flow` | Authentication model, tokens, roles |
| `qason/project/{name}/bug-patterns` | Recurring bug types, hot spots |
| `qason/project/{name}/flaky-tests` | Known flaky tests and root causes |
| `qason/project/{name}/test-data` | Test users, fixtures, seed data |
| `qason/project/{name}/coverage-map` | What is tested, what is not |

## Workflow

1. **Detect project name**:
   - Try `git remote get-url origin` → extract repo name
   - Fallback: basename of working directory
2. **Load core context** (stack + conventions + domain):
   - `mem_search(query: "qason/project/{name}/stack", project: "{name}")`
   - `mem_search(query: "qason/project/{name}/conventions", project: "{name}")`
   - `mem_search(query: "qason/project/{name}/domain", project: "{name}")`
3. **For each result found**:
   - `mem_get_observation(id)` to get the FULL untruncated content
4. **Load task-specific context** based on what the user is asking:
   - Bug triage task → load `bug-patterns` + `flaky-tests`
   - Test generation task → load `test-data` + `coverage-map`
   - Review task → load `coverage-map` + `conventions`
5. **If nothing found** (cold start):
   - Tell the user: "This looks like the first time I am working on {name}. I will discover the context as I go and save it for next time."
   - Trigger `project-context-saver` skill to save discoveries as you make them
6. **If found**:
   - Summarize in 1-2 sentences what you already know about the project
   - Then proceed with the actual task

## Output Template

```markdown
## Project Context: {name}

**Status**: ✅ Loaded from memory | 🆕 Cold start (no prior context)

### Known
- **Stack**: [frameworks + versions from memory]
- **Conventions**: [key patterns from memory]
- **Domain**: [business rules from memory]

### Relevant for this task
- [bullet points of memory entries that apply to the current task]

### Will discover (cold start only)
- [list of things to auto-save during this task]
```

## Fallback: no engram available

QASON memory is 3-layered: engram preferred, file store fallback, install gate
at Acceleration+ presets. If engram MCP is not in your tool list:

1. Do NOT error out or mention memory to the user.
2. Load from `~/.qason/memory/qason/project/{name}/*.md` via the Read/Glob
   tools. See the `memory-fallback-protocol` skill for the exact layout.
3. If neither engram nor the file store have context for this project,
   treat it as a cold start — document discoveries through
   `project-context-saver`, which will choose the same backend on write.

Only if BOTH backends are unreachable (engram unavailable AND `~/.qason/memory/`
does not exist) do we proceed as pre-memory QASON did — work the task, no
persistence. At `acceleration` and above this state is a broken install; the
install gate should have caught it.
