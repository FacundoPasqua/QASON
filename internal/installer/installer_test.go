package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInstall_WritesAgentFiles verifies that Install writes the 3 agent
// files with correct YAML frontmatter and rebranded QASON prompt bodies.
func TestInstall_WritesAgentFiles(t *testing.T) {
	tests := []struct {
		name       string
		fileName   string
		agentID    string
		bodyMarker string
		color      string
	}{
		{
			name:       "qa-analyst",
			fileName:   "qa-analyst.md",
			agentID:    "qa-analyst",
			bodyMarker: "QA Analyst",
			color:      "cyan",
		},
		{
			name:       "qa-test-designer",
			fileName:   "qa-test-designer.md",
			agentID:    "qa-test-designer",
			bodyMarker: "QA Test Designer",
			color:      "green",
		},
		{
			name:       "qa-automator",
			fileName:   "qa-automator.md",
			agentID:    "qa-automator",
			bodyMarker: "QA Automator",
			color:      "purple",
		},
	}

	claudeDir := t.TempDir()
	result, err := Install(Options{ClaudeDir: claudeDir})
	if err != nil {
		t.Fatalf("Install() returned error: %v", err)
	}

	if len(result.Agents) != 3 {
		t.Fatalf("Result.Agents = %d entries, want 3 (got %v)", len(result.Agents), result.Agents)
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(claudeDir, "agents", tc.fileName)
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("expected agent file %s to exist, got error: %v", path, err)
			}
			content := string(data)

			if !strings.HasPrefix(content, "---\n") {
				t.Errorf("agent file %s should start with '---\\n', got prefix: %q", path, content[:min(20, len(content))])
			}
			if !strings.Contains(content, "name: "+tc.agentID) {
				t.Errorf("agent file %s should contain 'name: %s', content:\n%s", path, tc.agentID, content)
			}
			if !strings.Contains(content, "color: "+tc.color) {
				t.Errorf("agent file %s should contain %q — the colors are what make the pipeline legible in the task list", path, "color: "+tc.color)
			}
			if !strings.Contains(content, "model: sonnet") {
				t.Errorf("agent file %s should contain 'model: sonnet', content:\n%s", path, content)
			}
			if !strings.Contains(content, tc.bodyMarker) {
				t.Errorf("agent file %s should contain body marker %q, content:\n%s", path, tc.bodyMarker, content)
			}
		})
	}
}

// TestInstall_WritesSkillsAtomically verifies skills are installed under
// skills/qason/<id>/SKILL.md and that the staging dir is cleaned up.
func TestInstall_WritesSkillsAtomically(t *testing.T) {
	claudeDir := t.TempDir()
	result, err := Install(Options{ClaudeDir: claudeDir})
	if err != nil {
		t.Fatalf("Install() returned error: %v", err)
	}

	prdAnalyzerPath := filepath.Join(claudeDir, "skills", "qason", "prd-analyzer", "SKILL.md")
	if _, err := os.Stat(prdAnalyzerPath); err != nil {
		t.Errorf("expected %s to exist, got error: %v", prdAnalyzerPath, err)
	}

	unitTestGenPath := filepath.Join(claudeDir, "skills", "qason", "unit-test-gen", "SKILL.md")
	if _, err := os.Stat(unitTestGenPath); err != nil {
		t.Errorf("expected %s to exist, got error: %v", unitTestGenPath, err)
	}

	skillsRoot := filepath.Join(claudeDir, "skills", "qason")
	entries, err := os.ReadDir(skillsRoot)
	if err != nil {
		t.Fatalf("expected skills root %s to exist, got error: %v", skillsRoot, err)
	}
	dirCount := 0
	for _, e := range entries {
		if e.IsDir() {
			dirCount++
		}
	}

	if dirCount != result.Skills {
		t.Errorf("dir count under skills/qason = %d, want equal to Result.Skills = %d", dirCount, result.Skills)
	}
	if dirCount < 30 {
		t.Errorf("dir count under skills/qason = %d, want >= 30", dirCount)
	}

	stagingPath := filepath.Join(claudeDir, "skills", "qason.qason-staging")
	if _, err := os.Stat(stagingPath); !os.IsNotExist(err) {
		t.Errorf("staging dir %s should not exist after Install returns, stat err = %v", stagingPath, err)
	}
}

// TestInstall_OrchestratorBlockIdempotent verifies the CLAUDE.md orchestrator
// block is inserted once, preserves existing user content, and running
// Install twice does not duplicate the block.
func TestInstall_OrchestratorBlockIdempotent(t *testing.T) {
	claudeDir := t.TempDir()
	claudeMDPath := filepath.Join(claudeDir, "CLAUDE.md")

	if err := os.WriteFile(claudeMDPath, []byte("# My personal notes\n"), 0o644); err != nil {
		t.Fatalf("failed to pre-create CLAUDE.md: %v", err)
	}

	if _, err := Install(Options{ClaudeDir: claudeDir}); err != nil {
		t.Fatalf("first Install() returned error: %v", err)
	}
	if _, err := Install(Options{ClaudeDir: claudeDir}); err != nil {
		t.Fatalf("second Install() returned error: %v", err)
	}

	data, err := os.ReadFile(claudeMDPath)
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "# My personal notes") {
		t.Errorf("CLAUDE.md should preserve user content '# My personal notes', got:\n%s", content)
	}
	if !strings.Contains(content, "QASON QA Orchestrator") {
		t.Errorf("CLAUDE.md should contain 'QASON QA Orchestrator', got:\n%s", content)
	}

	beginCount := strings.Count(content, "<!-- QASON:BEGIN -->")
	if beginCount != 1 {
		t.Errorf("CLAUDE.md should contain exactly one '<!-- QASON:BEGIN -->' marker, found %d, content:\n%s", beginCount, content)
	}
}

// TestUninstall_RemovesEverythingPreservingUserContent verifies Uninstall
// removes agent files and the skills/qason dir, and strips the QASON block
// from CLAUDE.md while preserving user content.
func TestUninstall_RemovesEverythingPreservingUserContent(t *testing.T) {
	claudeDir := t.TempDir()
	claudeMDPath := filepath.Join(claudeDir, "CLAUDE.md")

	if err := os.WriteFile(claudeMDPath, []byte("# My personal notes\n"), 0o644); err != nil {
		t.Fatalf("failed to pre-create CLAUDE.md: %v", err)
	}

	if _, err := Install(Options{ClaudeDir: claudeDir}); err != nil {
		t.Fatalf("Install() returned error: %v", err)
	}

	if err := Uninstall(Options{ClaudeDir: claudeDir}); err != nil {
		t.Fatalf("Uninstall() returned error: %v", err)
	}

	agentFiles := []string{"qa-analyst.md", "qa-test-designer.md", "qa-automator.md"}
	for _, f := range agentFiles {
		path := filepath.Join(claudeDir, "agents", f)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("agent file %s should be removed after Uninstall, stat err = %v", path, err)
		}
	}

	skillsQasonPath := filepath.Join(claudeDir, "skills", "qason")
	if _, err := os.Stat(skillsQasonPath); !os.IsNotExist(err) {
		t.Errorf("skills/qason dir %s should be removed after Uninstall, stat err = %v", skillsQasonPath, err)
	}

	data, err := os.ReadFile(claudeMDPath)
	if err != nil {
		t.Fatalf("failed to read CLAUDE.md after Uninstall: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "# My personal notes") {
		t.Errorf("CLAUDE.md should still preserve user content '# My personal notes' after Uninstall, got:\n%s", content)
	}
	if strings.Contains(content, "<!-- QASON:BEGIN -->") || strings.Contains(content, "<!-- QASON:END -->") {
		t.Errorf("CLAUDE.md should not contain QASON markers after Uninstall, got:\n%s", content)
	}
}

// TestUninstall_NoopWhenNotInstalled verifies Uninstall is a safe no-op
// when QASON was never installed in the given ClaudeDir.
func TestUninstall_NoopWhenNotInstalled(t *testing.T) {
	claudeDir := t.TempDir()

	if err := Uninstall(Options{ClaudeDir: claudeDir}); err != nil {
		t.Errorf("Uninstall() on a never-installed dir should return nil, got: %v", err)
	}
}

// TestInstall_CopilotOptIn verifies the Copilot instructions file is written
// only when CopilotDir is set. VS Code finds the agents and skills in
// ~/.claude by itself; this file is the one piece it cannot discover, so
// writing it unasked would edit a project the user never pointed us at.
func TestInstall_CopilotOptIn(t *testing.T) {
	claudeDir := t.TempDir()

	res, err := Install(Options{ClaudeDir: claudeDir})
	if err != nil {
		t.Fatalf("Install() without CopilotDir returned error: %v", err)
	}
	if res.Copilot != "" {
		t.Errorf("Result.Copilot should be empty when CopilotDir is unset, got %q", res.Copilot)
	}

	projectDir := t.TempDir()
	res, err = Install(Options{ClaudeDir: claudeDir, CopilotDir: projectDir})
	if err != nil {
		t.Fatalf("Install() with CopilotDir returned error: %v", err)
	}
	want := filepath.Join(projectDir, ".github", "copilot-instructions.md")
	if res.Copilot != want {
		t.Errorf("Result.Copilot = %q, want %q", res.Copilot, want)
	}
	data, err := os.ReadFile(want)
	if err != nil {
		t.Fatalf("expected Copilot instructions at %s, got error: %v", want, err)
	}
	if !strings.Contains(string(data), "QASON QA Orchestrator") {
		t.Errorf("Copilot instructions should contain the orchestrator, got:\n%s", data)
	}
}

// TestInstall_CopilotIdempotentAndReversible verifies the Copilot file gets
// the same managed-block contract as CLAUDE.md: user content survives, the
// block never duplicates, and Uninstall strips the block but leaves the
// user's own instructions behind.
func TestInstall_CopilotIdempotentAndReversible(t *testing.T) {
	claudeDir := t.TempDir()
	projectDir := t.TempDir()
	path := filepath.Join(projectDir, ".github", "copilot-instructions.md")

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create .github dir: %v", err)
	}
	if err := os.WriteFile(path, []byte("# Team Copilot rules\n"), 0o644); err != nil {
		t.Fatalf("failed to pre-create Copilot instructions: %v", err)
	}

	opts := Options{ClaudeDir: claudeDir, CopilotDir: projectDir}
	for i := 0; i < 2; i++ {
		if _, err := Install(opts); err != nil {
			t.Fatalf("Install() #%d returned error: %v", i+1, err)
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read Copilot instructions: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "# Team Copilot rules") {
		t.Errorf("Copilot instructions should preserve user content, got:\n%s", content)
	}
	if n := strings.Count(content, "<!-- QASON:BEGIN -->"); n != 1 {
		t.Errorf("Copilot instructions should contain exactly one begin marker, found %d", n)
	}

	if err := Uninstall(opts); err != nil {
		t.Fatalf("Uninstall() returned error: %v", err)
	}
	data, err = os.ReadFile(path)
	if err != nil {
		t.Fatalf("Uninstall() should leave the Copilot file in place, got error: %v", err)
	}
	content = string(data)
	if strings.Contains(content, "<!-- QASON:BEGIN -->") {
		t.Errorf("Uninstall() should strip the QASON block, got:\n%s", content)
	}
	if !strings.Contains(content, "# Team Copilot rules") {
		t.Errorf("Uninstall() must not destroy the user's own Copilot rules, got:\n%s", content)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
