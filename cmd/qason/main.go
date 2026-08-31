// Command qason installs the QASON educational QA agents into Claude
// Code. Three sub-commands, zero dependencies — the whole CLI fits in
// one file on purpose: it is course material first, tooling second.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/FacundoPasqua/qason/internal/installer"
	"github.com/FacundoPasqua/qason/internal/tui"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		// No sub-command: launch the interactive wizard.
		if err := tui.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "qason: %v\n", err)
			os.Exit(1)
		}
		return
	}

	cmd := os.Args[1]
	flags := flag.NewFlagSet(cmd, flag.ExitOnError)
	claudeDir := flags.String("claude-dir", installer.DefaultClaudeDir(), "Claude Code config directory")
	copilotDir := flags.String("copilot", "", "project root that also gets .github/copilot-instructions.md")
	_ = flags.Parse(os.Args[2:])

	opts := installer.Options{ClaudeDir: *claudeDir, CopilotDir: *copilotDir}

	switch cmd {
	case "install":
		res, err := installer.Install(opts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "qason: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("QASON installed into %s\n", *claudeDir)
		fmt.Printf("  ✓ %d QA agents (qa-analyst, qa-test-designer, qa-automator)\n", len(res.Agents))
		fmt.Printf("  ✓ %d skills\n", res.Skills)
		fmt.Println("  ✓ orchestrator registered in CLAUDE.md")
		if res.Copilot != "" {
			fmt.Printf("  ✓ orchestrator registered in %s\n", res.Copilot)
		}
		fmt.Println()
		fmt.Println("Open Claude Code and try: \"Analyze this ticket and create tests for it\"")
		if res.Copilot != "" {
			fmt.Println("Copilot reads the agents and skills from ~/.claude on its own.")
		}
	case "uninstall":
		if err := installer.Uninstall(opts); err != nil {
			fmt.Fprintf(os.Stderr, "qason: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("QASON uninstalled. Your own CLAUDE.md content was preserved.")
	case "version":
		fmt.Println("qason " + version)
	default:
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Println(`qason — educational QA agents for Claude Code (the QATES teaching edition)

Usage:
  qason             launch the interactive install wizard
  qason install     install the 3 QA agents, skills, and orchestrator
  qason uninstall   remove everything qason installed
  qason version     print the version

Flags:
  --claude-dir      Claude Code config directory (default: ~/.claude)
  --copilot DIR     also register the orchestrator for GitHub Copilot, in
                    DIR/.github/copilot-instructions.md (VS Code already
                    reads the agents and skills from ~/.claude)`)
}
