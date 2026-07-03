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
	}{
		{
			name:       "qa-analyst",
			fileName:   "qa-analyst.md",
			agentID:    "qa-analyst",
			bodyMarker: "QA Analyst",
		},
		{
			name:       "qa-test-designer",
			fileName:   "qa-test-designer.md",
			agentID:    "qa-test-designer",
			bodyMarker: "QA Test Designer",
		},
		{
			name:       "qa-automator",
			fileName:   "qa-automator.md",
			agentID:    "qa-automator",
			bodyMarker: "QA Automator",
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
