// Package tui implements the QASON install wizard: a three-screen
// Bubble Tea app (menu → working → done). Deliberately small — like
// the installer, it doubles as course material for reading.
package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/FacundoPasqua/qason/internal/installer"
)

// screen tracks which view is active.
type screen int

const (
	screenMenu    screen = iota // menu: Install / Uninstall / Quit
	screenWorking               // action running (spinner)
	screenDone                  // result or error
)

// actionDoneMsg reports the finished action.
type actionDoneMsg struct {
	agents    int   // agents installed (install only)
	skills    int   // skills installed (install only)
	uninstall bool  // true when the action was uninstall
	err       error // non-nil on failure
}

// tickMsg drives the spinner while an action runs.
type tickMsg struct{}

func tick() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

var menuItems = []string{"Install", "Uninstall", "Quit"}

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
	itemStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	okStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("70"))
	errStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
)

// Model is the Bubble Tea model for the wizard.
type Model struct {
	screen screen
	cursor int
	width  int
	height int
	frame  int
	action string // "install" | "uninstall"
	done   actionDoneMsg

	installFn   func() (agents, skills int, err error)
	uninstallFn func() error
}

// New returns the production model, wired to the real installer
// against the default Claude Code directory.
func New() Model {
	dir := installer.DefaultClaudeDir()
	return NewWithActions(
		func() (int, int, error) {
			res, err := installer.Install(installer.Options{ClaudeDir: dir})
			return len(res.Agents), res.Skills, err
		},
		func() error {
			return installer.Uninstall(installer.Options{ClaudeDir: dir})
		},
	)
}

// NewWithActions returns a model with injectable actions — the seam
// tests use so no test ever touches a real ~/.claude.
func NewWithActions(install func() (int, int, error), uninstall func() error) Model {
	return Model{installFn: install, uninstallFn: uninstall}
}

// Run starts the interactive wizard and blocks until it exits.
func Run() error {
	_, err := tea.NewProgram(New()).Run()
	return err
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case actionDoneMsg:
		m.done = msg
		m.screen = screenDone
		return m, nil
	case tickMsg:
		if m.screen == screenWorking {
			m.frame++
			return m, tick()
		}
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(key tea.KeyMsg) (tea.Model, tea.Cmd) {
	if key.Type == tea.KeyCtrlC {
		return m, tea.Quit
	}
	switch m.screen {
	case screenMenu:
		switch key.String() {
		case "q", "esc":
			return m, tea.Quit
		case "j", "down":
			if m.cursor < len(menuItems)-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			switch m.cursor {
			case 0:
				m.screen = screenWorking
				m.action = "install"
				return m, tea.Batch(m.runInstall, tick())
			case 1:
				m.screen = screenWorking
				m.action = "uninstall"
				return m, tea.Batch(m.runUninstall, tick())
			default:
				return m, tea.Quit
			}
		}
	case screenDone:
		switch key.String() {
		case "enter", "q", "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

// runInstall and runUninstall are tea.Cmds: they run the injected
// action (blocking, off the UI goroutine) and deliver the outcome.
func (m Model) runInstall() tea.Msg {
	agents, skills, err := m.installFn()
	return actionDoneMsg{agents: agents, skills: skills, err: err}
}

func (m Model) runUninstall() tea.Msg {
	return actionDoneMsg{uninstall: true, err: m.uninstallFn()}
}

func (m Model) View() string {
	var b strings.Builder
	b.WriteString("\n")
	b.WriteString(renderLogo(m.width))
	b.WriteString("\n\n")

	switch m.screen {
	case screenMenu:
		for i, item := range menuItems {
			if i == m.cursor {
				b.WriteString("  ❯ " + selectedStyle.Render(item) + "\n")
			} else {
				b.WriteString("    " + itemStyle.Render(item) + "\n")
			}
		}
		b.WriteString("\n" + helpStyle.Render("  ↑/↓ move · enter select · q quit"))
	case screenWorking:
		verb := "Installing QASON"
		if m.action == "uninstall" {
			verb = "Uninstalling QASON"
		}
		b.WriteString(fmt.Sprintf("  %s %s...", spinnerFrames[m.frame%len(spinnerFrames)], verb))
	case screenDone:
		switch {
		case m.done.err != nil:
			b.WriteString("  " + errStyle.Render("✗ "+m.done.err.Error()))
		case m.done.uninstall:
			b.WriteString("  " + okStyle.Render("✓ QASON uninstalled.") + " Your own CLAUDE.md content was preserved.")
		default:
			b.WriteString("  " + okStyle.Render("✓ Installed:") +
				fmt.Sprintf(" %d QA agents, %d skills, orchestrator registered in CLAUDE.md\n", m.done.agents, m.done.skills))
			b.WriteString("\n  Open Claude Code and try: \"Analyze this ticket and create tests for it\"")
		}
		b.WriteString("\n\n" + helpStyle.Render("  enter to exit"))
	}
	b.WriteString("\n")
	return b.String()
}
