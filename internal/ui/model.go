package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stonehub/hecksweeper/internal/constants"
	"github.com/stonehub/hecksweeper/internal/game"
)

// Model represents the UI state for bubbletea
type Model struct {
	game       *game.Game
	glyphs     constants.GlyphSet
	width      int
	height     int
	seed       int64
	useUnicode bool
	quitting   bool
}

// NewModel creates a new UI model
func NewModel(width, height int, seed int64, useUnicode bool) Model {
	glyphs := constants.ASCIIGlyphs
	if useUnicode {
		glyphs = constants.UnicodeGlyphs
	}

	return Model{
		game:       game.NewGame(width, height, seed, useUnicode),
		glyphs:     glyphs,
		width:      width,
		height:     height,
		seed:       seed,
		useUnicode: useUnicode,
		quitting:   false,
	}
}

// Init initializes the model (required by bubbletea)
func (m Model) Init() tea.Cmd {
	return nil
}
