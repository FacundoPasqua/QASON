// Command qason installs the QASON educational QA agents into Claude
// Code. Three sub-commands, zero dependencies — the whole CLI fits in
// one file on purpose: it is course material first, tooling second.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/FacundoPasqua/qason/internal/installer"
)

const version = "0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	cmd := os.Args[1]
	flags := flag.NewFlagSet(cmd, flag.ExitOnError)
	claudeDir := flags.String("claude-dir", defaultClaudeDir(), "Claude Code config directory")
	_ = flags.Parse(os.Args[2:])

	switch cmd {
	case "install":
		res, err := installer.Install(installer.Options{ClaudeDir: *claudeDir})
		if err != nil {
			fmt.Fprintf(os.Stderr, "qason: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("QASON installed into %s\n", *claudeDir)
		fmt.Printf("  ✓ %d QA agents (qa-analyst, qa-test-designer, qa-automator)\n", len(res.Agents))
		fmt.Printf("  ✓ %d skills\n", res.Skills)
		fmt.Println("  ✓ orchestrator registered in CLAUDE.md")
		fmt.Println()
		fmt.Println("Open Claude Code and try: \"Analyze this ticket and create tests for it\"")
	case "uninstall":
		if err := installer.Uninstall(installer.Options{ClaudeDir: *claudeDir}); err != nil {
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

func defaultClaudeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".claude"
	}
	return filepath.Join(home, ".claude")
}

func usage() {
	fmt.Println(`qason — educational QA agents for Claude Code (the QATES teaching edition)

Usage:
  qason install     install the 3 QA agents, skills, and orchestrator
  qason uninstall   remove everything qason installed
  qason version     print the version

Flags:
  --claude-dir      Claude Code config directory (default: ~/.claude)`)
}
