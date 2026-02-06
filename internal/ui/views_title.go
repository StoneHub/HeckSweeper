package ui

import (
	"fmt"
	"strings"
)

// renderTitle renders the main title/menu screen
func (m Model) renderTitle() string {
	var b strings.Builder

	title := `
  ╔══════════════════════════════════════╗
  ║                                      ║
  ║   ██╗  ██╗███████╗ ██████╗██╗  ██╗  ║
  ║   ██║  ██║██╔════╝██╔════╝██║ ██╔╝  ║
  ║   ███████║█████╗  ██║     █████╔╝   ║
  ║   ██╔══██║██╔══╝  ██║     ██╔═██╗   ║
  ║   ██║  ██║███████╗╚██████╗██║  ██╗  ║
  ║   ╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝  ║
  ║          S W E E P E R               ║
  ║                                      ║
  ║      A roguelite minesweeper         ║
  ║                                      ║
  ╚══════════════════════════════════════╝`

	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	menuItems := []string{
		"    [N] New Run",
		"    [D] Daily Challenge",
		"    [L] Leaderboard",
		"    [Q] Quit",
	}
	b.WriteString(menuStyle.Render(strings.Join(menuItems, "\n")))
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render(fmt.Sprintf("  Seed: %d", m.seed)))

	return boardContainerStyle.Render(b.String())
}
