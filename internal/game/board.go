package game

import (
	"github.com/stonehub/hecksweeper/internal/constants"
)

// NewBoard creates a new game board with the specified dimensions
func NewBoard(width, height int) *Board {
	cells := make([]Cell, width*height)
	return &Board{
		Width:         width,
		Height:        height,
		Cells:         cells,
		Monsters:      make(map[int]bool),
		RemainingSafe: width * height,
	}
}

// GetIndex converts 2D coordinates to 1D array index
func (b *Board) GetIndex(x, y int) int {
	return y*b.Width + x
}

// GetPosition converts 1D index to 2D coordinates
func (b *Board) GetPosition(index int) Position {
	return Position{
		X: index % b.Width,
		Y: index / b.Width,
	}
}

// IsValidPosition checks if coordinates are within board bounds
func (b *Board) IsValidPosition(x, y int) bool {
	return x >= 0 && x < b.Width && y >= 0 && y < b.Height
}

// GetCell returns the cell at the given position
func (b *Board) GetCell(x, y int) *Cell {
	if !b.IsValidPosition(x, y) {
		return nil
	}
	return &b.Cells[b.GetIndex(x, y)]
}

// GetNeighborIndices returns indices of all valid neighbors (8-directional)
func (b *Board) GetNeighborIndices(index int) []int {
	pos := b.GetPosition(index)
	neighbors := []int{}

	directions := []Position{
		{-1, -1}, {0, -1}, {1, -1}, // Top row
		{-1, 0}, {1, 0}, // Middle row (skip center)
		{-1, 1}, {0, 1}, {1, 1}, // Bottom row
	}

	for _, dir := range directions {
		nx, ny := pos.X+dir.X, pos.Y+dir.Y
		if b.IsValidPosition(nx, ny) {
			neighbors = append(neighbors, b.GetIndex(nx, ny))
		}
	}

	return neighbors
}

// CalculateThreatCounts computes the number of adjacent monsters for each cell
func (b *Board) CalculateThreatCounts() {
	for i := range b.Cells {
		if b.Cells[i].HasMonster {
			continue
		}

		count := 0
		for _, neighborIdx := range b.GetNeighborIndices(i) {
			if b.Cells[neighborIdx].HasMonster {
				count++
			}
		}
		b.Cells[i].ThreatCount = count
	}
}

// ToggleFlag toggles the flag state of a cell
func (b *Board) ToggleFlag(x, y int) bool {
	cell := b.GetCell(x, y)
	if cell == nil || cell.Revealed {
		return false
	}

	cell.Flagged = !cell.Flagged
	return true
}

// RevealCell reveals a cell and performs flood-fill if it's safe (0 threats)
func (b *Board) RevealCell(x, y int) (hitMonster bool, cellsRevealed int) {
	cell := b.GetCell(x, y)
	if cell == nil || cell.Revealed || cell.Flagged {
		return false, 0
	}

	// Check if we hit a monster
	if cell.HasMonster {
		cell.Revealed = true
		return true, 1
	}

	// Flood-fill reveal using BFS
	queue := []int{b.GetIndex(x, y)}
	visited := make(map[int]bool)

	for len(queue) > 0 {
		idx := queue[0]
		queue = queue[1:]

		if visited[idx] {
			continue
		}
		visited[idx] = true

		currentCell := &b.Cells[idx]
		if currentCell.Revealed || currentCell.Flagged {
			continue
		}

		// Reveal the cell
		currentCell.Revealed = true
		cellsRevealed++
		b.RemainingSafe--

		// If this cell has no adjacent monsters, add neighbors to queue
		if currentCell.ThreatCount == 0 {
			for _, neighborIdx := range b.GetNeighborIndices(idx) {
				if !visited[neighborIdx] {
					queue = append(queue, neighborIdx)
				}
			}
		}
	}

	return false, cellsRevealed
}

// RevealAllMonsters reveals all monster cells (called on loss)
func (b *Board) RevealAllMonsters() {
	for i := range b.Cells {
		if b.Cells[i].HasMonster {
			b.Cells[i].Revealed = true
		}
	}
}

// IsWon checks if the game is won (all safe cells revealed)
func (b *Board) IsWon() bool {
	return b.RemainingSafe == 0
}

// NewGame creates a new game with the specified settings
func NewGame(width, height int, seed int64, useUnicode bool) *Game {
	board := GenerateBoard(width, height, constants.DefaultMonsterDensity, seed)

	return &Game{
		Board:      board,
		CursorPos:  Position{X: width / 2, Y: height / 2},
		State:      constants.GameStatePlaying,
		MoveCount:  0,
		FlagCount:  0,
		Seed:       seed,
		UseUnicode: useUnicode,
	}
}

// MoveCursor moves the cursor in the specified direction
func (g *Game) MoveCursor(dx, dy int) {
	newX := g.CursorPos.X + dx
	newY := g.CursorPos.Y + dy

	if g.Board.IsValidPosition(newX, newY) {
		g.CursorPos.X = newX
		g.CursorPos.Y = newY
		g.MoveCount++
	}
}

// ToggleFlag toggles flag at cursor position
func (g *Game) ToggleFlag() {
	if g.State != constants.GameStatePlaying {
		return
	}

	if g.Board.ToggleFlag(g.CursorPos.X, g.CursorPos.Y) {
		cell := g.Board.GetCell(g.CursorPos.X, g.CursorPos.Y)
		if cell.Flagged {
			g.FlagCount++
		} else {
			g.FlagCount--
		}
	}
}

// RevealAtCursor reveals the cell at cursor position
func (g *Game) RevealAtCursor() {
	if g.State != constants.GameStatePlaying {
		return
	}

	hitMonster, _ := g.Board.RevealCell(g.CursorPos.X, g.CursorPos.Y)

	if hitMonster {
		g.State = constants.GameStateLost
		g.Board.RevealAllMonsters()
	} else if g.Board.IsWon() {
		g.State = constants.GameStateWon
	}
}
