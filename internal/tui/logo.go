package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// logoArt is the QASON wordmark, same block-character family as the
// QATES logo — students who later meet the full product should feel
// at home.
const logoArt = ` ██████╗  █████╗ ███████╗ ██████╗ ███╗   ██╗
██╔═══██╗██╔══██╗██╔════╝██╔═══██╗████╗  ██║
██║   ██║███████║███████╗██║   ██║██╔██╗ ██║
██║▄▄ ██║██╔══██║╚════██║██║   ██║██║╚██╗██║
╚██████╔╝██║  ██║███████║╚██████╔╝██║ ╚████║
 ╚══▀▀═╝ ╚═╝  ╚═╝╚══════╝ ╚═════╝ ╚═╝  ╚═══╝`

const logoTagline = "QASON · Educational QA agents for Claude Code"

// logoGradient shades the wordmark from light to deep orange, one
// color per line — same per-line gradient trick as QATES, in the
// QASON brand palette.
var logoGradient = []lipgloss.Color{"214", "208", "208", "202", "202", "166"}

var (
	brandStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
	taglineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Italic(true)
)

// logoMinWidth is the narrowest terminal that renders the wordmark
// without wrapping.
var logoMinWidth = lipgloss.Width(logoArt) + 4

// renderLogo draws the orange QASON wordmark and the tagline. Narrow
// terminals (or unknown width, before the first WindowSizeMsg) fall
// back to a one-line brand so nothing wraps.
func renderLogo(width int) string {
	if width > 0 && width < logoMinWidth {
		return brandStyle.Render("QASON") + taglineStyle.Render(" · "+logoTagline)
	}

	var b strings.Builder
	logoWidth := lipgloss.Width(logoArt)
	for i, line := range strings.Split(logoArt, "\n") {
		color := logoGradient[i%len(logoGradient)]
		b.WriteString(lipgloss.NewStyle().Foreground(color).Render(line))
		b.WriteString("\n")
	}
	b.WriteString(strings.Repeat(" ", centerPad(logoWidth, lipgloss.Width(logoTagline))))
	b.WriteString(taglineStyle.Render(logoTagline))
	return b.String()
}

func centerPad(total, inner int) int {
	pad := (total - inner) / 2
	if pad < 0 {
		return 0
	}
	return pad
}
