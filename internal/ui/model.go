package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stonehub/hecksweeper/internal/constants"
	"github.com/stonehub/hecksweeper/internal/game"
)

// Model represents the UI state for bubbletea
type Model struct {
	run        *game.Run
	runState   string // Current run state (title, playing, cleared, dead, summary)
	lastResult game.FloorResult
	glyphs     constants.GlyphSet
	seed       int64
	useUnicode bool
	quitting   bool
}

// NewModel creates a new UI model starting at the title screen
func NewModel(seed int64, useUnicode bool) Model {
	glyphs := constants.ASCIIGlyphs
	if useUnicode {
		glyphs = constants.UnicodeGlyphs
	}

	return Model{
		runState:   constants.RunStateTitle,
		glyphs:     glyphs,
		seed:       seed,
		useUnicode: useUnicode,
		quitting:   false,
	}
}

// Init initializes the model (required by bubbletea)
func (m Model) Init() tea.Cmd {
	return nil
}

// startNewRun begins a new roguelite run and starts floor 1
func (m *Model) startNewRun(isDaily bool) {
	m.run = game.NewRun(m.seed, m.useUnicode, isDaily)
	m.run.StartNextFloor()
	m.runState = constants.RunStatePlaying
}
