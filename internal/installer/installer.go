// Package installer writes the QASON content into Claude Code's config
// directory: three QA sub-agents, their skills, and the orchestrator
// prompt that routes work between them.
//
// Everything here is deliberately simple and readable — QASON is course
// material, and this installer is part of the lesson: agent files are
// markdown with YAML frontmatter, skills are directories of SKILL.md
// files, and the orchestrator is a managed block inside ~/.claude/CLAUDE.md.
package installer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/FacundoPasqua/qason/internal/assets"
)

const (
	// beginMarker / endMarker delimit the QASON-owned block in
	// CLAUDE.md. Everything outside the markers belongs to the user
	// and is never touched.
	beginMarker = "<!-- QASON:BEGIN -->"
	endMarker   = "<!-- QASON:END -->"
	// stagingSuffix builds the temporary sibling directory used for
	// the atomic skills swap: write everything to staging, then one
	// rename. An interrupted install can never leave a half-populated
	// skills directory.
	stagingSuffix = ".qason-staging"
)

// Options configures where QASON installs.
type Options struct {
	ClaudeDir string // Claude Code config dir, e.g. ~/.claude
	// CopilotDir, when set, is a project root that also gets the
	// orchestrator written to .github/copilot-instructions.md. See
	// ensureCopilotBlock for why only that one file is needed.
	CopilotDir string
}

// DefaultClaudeDir returns ~/.claude, the standard Claude Code config
// directory. Shared by the CLI flags and the TUI so both default to
// the same place.
func DefaultClaudeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".claude"
	}
	return filepath.Join(home, ".claude")
}

// Result reports what Install wrote.
type Result struct {
	Agents  []string // absolute paths of the agent files written
	Skills  int      // number of skills installed
	Copilot string   // path of the Copilot instructions file, empty if not requested
}

// agentSpecs binds each sub-agent to its embedded prompt, the one-line
// description Claude Code shows in the Agent tool, and its display color.
//
// The colors are teaching aids, not decoration: when the Spec-to-Test
// pipeline runs, the task list shows three distinctly colored agents
// handing work to each other, which makes specialization visible at a
// glance. Red is deliberately unused — it reads as failure in a QA
// context, and no agent should look broken just by running.
//
// Claude Code accepts eight named colors: red, blue, green, yellow,
// purple, orange, pink, cyan.
var agentSpecs = []struct {
	ID          string
	Asset       string
	Description string
	Color       string
}{
	{"qa-analyst", "analyst.md", "Analyzes requirements, generates test plans and risk matrices", "cyan"},
	{"qa-test-designer", "designer.md", "Designs test cases: functional, edge, negative, exploratory, contract, data-driven", "green"},
	{"qa-automator", "automator.md", "Generates automation scripts: unit, integration, e2e, performance, security, a11y", "purple"},
}

// Install writes the QASON agents, skills, and orchestrator block into
// opts.ClaudeDir. Safe to run repeatedly: agent files are overwritten,
// skills are swapped atomically, and the CLAUDE.md block is replaced,
// never duplicated.
func Install(opts Options) (Result, error) {
	var res Result
	if opts.ClaudeDir == "" {
		return res, fmt.Errorf("claude dir is required")
	}

	agentsDir := filepath.Join(opts.ClaudeDir, "agents")
	if err := os.MkdirAll(agentsDir, 0o755); err != nil {
		return res, fmt.Errorf("create agents dir: %w", err)
	}
	for _, a := range agentSpecs {
		body, err := assets.QAAgentPrompts.ReadFile("qaagents/" + a.Asset)
		if err != nil {
			return res, fmt.Errorf("embedded prompt %s: %w", a.Asset, err)
		}
		path := filepath.Join(agentsDir, a.ID+".md")
		if err := os.WriteFile(path, renderAgentFile(a.ID, a.Description, a.Color, body), 0o644); err != nil {
			return res, fmt.Errorf("write %s: %w", path, err)
		}
		res.Agents = append(res.Agents, path)
	}

	n, err := installSkills(opts.ClaudeDir)
	if err != nil {
		return res, err
	}
	res.Skills = n

	if err := ensureOrchestratorBlock(opts.ClaudeDir); err != nil {
		return res, err
	}

	if opts.CopilotDir != "" {
		path, err := ensureCopilotBlock(opts.CopilotDir)
		if err != nil {
			return res, err
		}
		res.Copilot = path
	}
	return res, nil
}

// Uninstall removes everything Install wrote, preserving any user
// content in CLAUDE.md (and in the Copilot instructions file) outside
// the QASON markers. Running it on a machine where QASON was never
// installed is a no-op.
func Uninstall(opts Options) error {
	for _, a := range agentSpecs {
		_ = os.Remove(filepath.Join(opts.ClaudeDir, "agents", a.ID+".md"))
	}
	_ = os.RemoveAll(filepath.Join(opts.ClaudeDir, "skills", "qason"))

	if opts.CopilotDir != "" {
		_ = stripManagedBlock(copilotPath(opts.CopilotDir))
	}
	return stripManagedBlock(filepath.Join(opts.ClaudeDir, "CLAUDE.md"))
}

// renderAgentFile assembles the Claude Code sub-agent file: YAML
// frontmatter (name, description, model, color) followed by the prompt
// body.
func renderAgentFile(id, description, color string, body []byte) []byte {
	var b strings.Builder
	b.WriteString("---\n")
	fmt.Fprintf(&b, "name: %s\n", id)
	fmt.Fprintf(&b, "description: %s\n", description)
	b.WriteString("model: sonnet\n")
	fmt.Fprintf(&b, "color: %s\n", color)
	b.WriteString("---\n\n")
	b.Write(body)
	return []byte(b.String())
}

// installSkills copies every embedded skill into
// <ClaudeDir>/skills/qason/<id>/SKILL.md using the staging + atomic
// rename pattern: either the full new set lands, or the previous state
// survives untouched.
func installSkills(claudeDir string) (int, error) {
	finalDir := filepath.Join(claudeDir, "skills", "qason")
	stagingDir := finalDir + stagingSuffix

	if err := os.MkdirAll(filepath.Dir(finalDir), 0o755); err != nil {
		return 0, fmt.Errorf("create skills root: %w", err)
	}
	_ = os.RemoveAll(stagingDir) // stale staging from an interrupted run

	count := 0
	err := fs.WalkDir(assets.Skills, "skills", func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := assets.Skills.ReadFile(p)
		if err != nil {
			return err
		}
		rel := strings.TrimPrefix(p, "skills/")
		target := filepath.Join(stagingDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		if strings.HasSuffix(p, "/SKILL.md") {
			count++
		}
		return os.WriteFile(target, data, 0o644)
	})
	if err != nil {
		_ = os.RemoveAll(stagingDir)
		return 0, fmt.Errorf("stage skills: %w", err)
	}

	if err := os.RemoveAll(finalDir); err != nil {
		_ = os.RemoveAll(stagingDir)
		return 0, fmt.Errorf("clear previous skills: %w", err)
	}
	if err := os.Rename(stagingDir, finalDir); err != nil {
		return 0, fmt.Errorf("swap staging to final: %w", err)
	}
	return count, nil
}

// ensureOrchestratorBlock rewrites CLAUDE.md as: user content (with any
// previous QASON block stripped) + one fresh QASON block. Idempotent by
// construction — strip-then-append can never duplicate.
func ensureOrchestratorBlock(claudeDir string) error {
	return writeManagedBlock(filepath.Join(claudeDir, "CLAUDE.md"))
}

// copilotPath returns the instructions file GitHub Copilot reads for a
// project: <root>/.github/copilot-instructions.md.
func copilotPath(projectDir string) string {
	return filepath.Join(projectDir, ".github", "copilot-instructions.md")
}

// ensureCopilotBlock writes the same orchestrator into the file Copilot
// reads, and returns its path.
//
// Only this one file is needed. VS Code already discovers the agents and
// skills QASON installs: it scans ~/.claude/agents for Claude-format
// sub-agents and ~/.claude/skills for SKILL.md files. What it does NOT
// read is CLAUDE.md, so without this the three specialists are present
// but nothing routes work between them — the pipeline, which is the
// whole lesson, never starts.
func ensureCopilotBlock(projectDir string) (string, error) {
	path := copilotPath(projectDir)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", fmt.Errorf("create .github dir: %w", err)
	}
	if err := writeManagedBlock(path); err != nil {
		return "", err
	}
	return path, nil
}

// writeManagedBlock rewrites path as: existing content with any previous
// QASON block stripped, then one fresh block. The file need not exist.
func writeManagedBlock(path string) error {
	orch, err := assets.QAAgentPrompts.ReadFile("qaagents/orchestrator.md")
	if err != nil {
		return fmt.Errorf("embedded orchestrator: %w", err)
	}

	existing, _ := os.ReadFile(path) // missing file reads as empty

	content := stripBlock(string(existing))
	if content != "" {
		content += "\n\n"
	}
	content += beginMarker + "\n" + strings.TrimSpace(string(orch)) + "\n" + endMarker + "\n"
	return os.WriteFile(path, []byte(content), 0o644)
}

// stripManagedBlock removes the QASON block from path, leaving the
// user's own content. A missing file is not an error: nothing to strip.
func stripManagedBlock(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	stripped := stripBlock(string(data))
	if stripped != "" {
		stripped += "\n"
	}
	return os.WriteFile(path, []byte(stripped), 0o644)
}

// stripBlock removes the QASON-owned block (markers included) and
// trims trailing whitespace, leaving only user content.
func stripBlock(s string) string {
	start := strings.Index(s, beginMarker)
	if start >= 0 {
		end := strings.Index(s, endMarker)
		if end >= 0 {
			end += len(endMarker)
		} else {
			end = len(s) // torn block: strip to EOF rather than leak half a block
		}
		s = s[:start] + s[end:]
	}
	return strings.TrimRight(s, " \t\n")
}
