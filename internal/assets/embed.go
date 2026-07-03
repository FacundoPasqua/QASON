// Package assets embeds everything QASON installs: the three QA agent
// prompts, the orchestrator prompt, and the skills they use. The
// content is the educational subset of QATES — same prompts, same
// skills, rebranded, so what students learn transfers 1:1 to the full
// product.
package assets

import "embed"

// QAAgentPrompts embeds the agent system prompts:
// qaagents/{analyst,designer,automator,orchestrator}.md
//
//go:embed all:qaagents
var QAAgentPrompts embed.FS

// Skills embeds every skill as skills/<id>/SKILL.md.
//
//go:embed all:skills
var Skills embed.FS
