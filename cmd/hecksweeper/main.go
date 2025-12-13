package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stonehub/hecksweeper/internal/constants"
	"github.com/stonehub/hecksweeper/internal/ui"
)

func main() {
	// Parse command-line flags
	width := flag.Int("width", constants.DefaultWidth, "Board width")
	height := flag.Int("height", constants.DefaultHeight, "Board height")
	seed := flag.Int64("seed", 0, "Random seed (0 for random)")
	useUnicode := flag.Bool("unicode", true, "Use Unicode glyphs (false for ASCII)")
	flag.Parse()

	// Use current time as seed if not specified
	if *seed == 0 {
		*seed = time.Now().UnixNano()
	}

	// Create the bubbletea model
	model := ui.NewModel(*width, *height, *seed, *useUnicode)

	// Run the program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running game: %v\n", err)
		os.Exit(1)
	}
}
