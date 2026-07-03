---
name: jira-thin-batch-update
description: >
  Comments on or transitions N Jira tickets efficiently. Uses the bug-thin
  format by default. Reports successes by key only, failures verbosely.
  Trigger: When the user wants to update multiple Jira tickets in one
  operation (N>1) — typically after a CI run, smoke suite, regression
  sweep, or batch validation.
license: MIT
metadata:
  author: QASON
  version: "0.1.0"
  agent: qa-ops
  phase: release
---

# Jira Thin Batch Update

## When to Use

- User wants to comment on **multiple Jira tickets** with the same operation (e.g. "comentá en estos 6 tickets con el resultado del último run").
- User wants to transition **multiple tickets** through the same workflow step (e.g. "movelos a 'Ready for QA'").
- After a CI run, smoke suite, or regression sweep that produced per-ticket pass/fail signals.

Do NOT use this skill when:

- Filing a **new bug ticket** — that's `jira-bug-reporter`. This skill updates existing tickets, doesn't create them.
- Commenting on a **single ticket** with full context — `jira_add_comment` directly is cleaner. This skill is for N>1.
- The per-ticket update needs **different content per ticket and full repro steps**. Verbose-per-ticket is `jira-bug-reporter`'s job, not this one.

## Important — what this skill does NOT do

It does NOT batch the API calls. **Atlassian Cloud has no bulk-comment endpoint**, so each ticket still gets its own `jira_add_comment` (or `jira_transition_issue`) MCP call. The token optimization is **in the OUTPUT format** the agent emits to the user — not in the underlying API surface.

If a future Atlassian release adds a true bulk endpoint, the workflow below collapses naturally to a single call.

## Critical Rules

1. **Confirm the ticket list before any MCP call.** Echo back the keys to the user: "Voy a actualizar EAAAAC-203, 204, 205, 206. Confirmá." This is the scope check the orchestrator demands.
2. **Use bug-thin format for every comment.** Per ticket: `[EMOJI] {key}: {outcome} ({reason})`. Reference: `bug-thin` skill.
3. **One MCP call per ticket.** Atlassian cloud limit, not ours. Loop sequentially.
4. **Stop on systemic failure.** If 3 consecutive `jira_add_comment` calls fail with the same error (auth, rate limit, network), stop and report. Don't retry the remaining N-3 — the issue is not per-ticket.
5. **Aggregate the response.** Success-by-key, failure-with-reason. Do NOT echo full comment bodies back. The user already gave them.

## Workflow

1. Parse the ticket key list from the user's request. If ambiguous (e.g. "los tickets de hoy"), ask for explicit keys before any call.
2. Build the per-ticket comment using the bug-thin format. The user-provided context (test result, validation outcome) goes in the `outcome` and `reason` slots — one line each.
3. For each key in the list:
   - Call `jira_add_comment(issue_key=KEY, comment=THIN_COMMENT)` (or `jira_transition_issue` for status changes).
   - Track the outcome (success / fail-reason) — do NOT print after each call.
   - On 3rd consecutive failure with same error: STOP, report what completed.
4. After the loop completes (or stops), emit a single aggregated response.

## Output Template (the agent's reply to the user)

```
Comments added to N tickets:
✅ KEY-1, KEY-2, KEY-3, ... (N successful)
⚠️ KEY-X (passed but flagged: response shape changed)
❌ KEY-Y (failed: 401 Unauthorized)
🚧 KEY-Z (blocked: ticket not found in project)

Stopped early: <reason if applicable, e.g. "3 consecutive 401s suggest expired token">
```

If everything succeeded:

```
Comments added to N tickets:
✅ KEY-1, KEY-2, KEY-3, ..., KEY-N (N successful)
```

That's it. No per-ticket detail, no echo of comment bodies, no markdown tables.

## Token Cost Awareness

Compared to using `jira-bug-reporter` for N tickets (each producing the full bug template in the comment + the agent's verbose summary in chat):

- **MCP calls**: same count (N, can't reduce).
- **Comment body sent to Jira**: one short line per ticket vs full bug template per ticket. ~10x reduction.
- **Agent output to user**: aggregated summary vs N verbose summaries. ~15x reduction in chat tokens.

The skill exists because the chat-output reduction was the dominant cost reported by the QA pilot 2026-04-27 (22 tickets × verbose summaries dominated the session's token budget).
