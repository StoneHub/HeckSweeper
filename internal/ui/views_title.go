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

	menu := strings.Join([]string{
		"    [N] New Run",
		"    [Q] Quit",
	}, "\n")
	b.WriteString(menuStyle.Render(menu))
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render(fmt.Sprintf("  Seed: %d", m.seed)))

	return boardContainerStyle.Render(b.String())
}
