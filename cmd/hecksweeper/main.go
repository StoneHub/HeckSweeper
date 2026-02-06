package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stonehub/hecksweeper/internal/ui"
)

func main() {
	seed := flag.Int64("seed", 0, "Random seed (0 for random)")
	useUnicode := flag.Bool("unicode", true, "Use Unicode glyphs (false for ASCII)")
	flag.Parse()

	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}

	model := ui.NewModel(*seed, *useUnicode)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running game: %v\n", err)
		os.Exit(1)
	}
}
