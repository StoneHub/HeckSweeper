package game

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard(24, 16)

	if board.Width != 24 {
		t.Errorf("Expected width 24, got %d", board.Width)
	}

	if board.Height != 16 {
		t.Errorf("Expected height 16, got %d", board.Height)
	}

	expectedCells := 24 * 16
	if len(board.Cells) != expectedCells {
		t.Errorf("Expected %d cells, got %d", expectedCells, len(board.Cells))
	}

	if board.RemainingSafe != expectedCells {
		t.Errorf("Expected RemainingSafe to be %d, got %d", expectedCells, board.RemainingSafe)
	}
}

func TestGetIndexAndPosition(t *testing.T) {
	board := NewBoard(10, 10)

	tests := []struct {
		x, y  int
		index int
	}{
		{0, 0, 0},
		{9, 0, 9},
		{0, 1, 10},
		{5, 5, 55},
		{9, 9, 99},
	}

	for _, tt := range tests {
		index := board.GetIndex(tt.x, tt.y)
		if index != tt.index {
			t.Errorf("GetIndex(%d, %d) = %d, want %d", tt.x, tt.y, index, tt.index)
		}

		pos := board.GetPosition(tt.index)
		if pos.X != tt.x || pos.Y != tt.y {
			t.Errorf("GetPosition(%d) = (%d, %d), want (%d, %d)", tt.index, pos.X, pos.Y, tt.x, tt.y)
		}
	}
}

func TestIsValidPosition(t *testing.T) {
	board := NewBoard(10, 10)

	tests := []struct {
		x, y  int
		valid bool
	}{
		{0, 0, true},
		{9, 9, true},
		{5, 5, true},
		{-1, 0, false},
		{0, -1, false},
		{10, 0, false},
		{0, 10, false},
		{10, 10, false},
	}

	for _, tt := range tests {
		valid := board.IsValidPosition(tt.x, tt.y)
		if valid != tt.valid {
			t.Errorf("IsValidPosition(%d, %d) = %v, want %v", tt.x, tt.y, valid, tt.valid)
		}
	}
}

func TestGetNeighborIndices(t *testing.T) {
	board := NewBoard(5, 5)

	// Center cell should have 8 neighbors
	centerIdx := board.GetIndex(2, 2)
	neighbors := board.GetNeighborIndices(centerIdx)
	if len(neighbors) != 8 {
		t.Errorf("Center cell should have 8 neighbors, got %d", len(neighbors))
	}

	// Corner cell should have 3 neighbors
	cornerIdx := board.GetIndex(0, 0)
	neighbors = board.GetNeighborIndices(cornerIdx)
	if len(neighbors) != 3 {
		t.Errorf("Corner cell should have 3 neighbors, got %d", len(neighbors))
	}

	// Edge cell should have 5 neighbors
	edgeIdx := board.GetIndex(2, 0)
	neighbors = board.GetNeighborIndices(edgeIdx)
	if len(neighbors) != 5 {
		t.Errorf("Edge cell should have 5 neighbors, got %d", len(neighbors))
	}
}

func TestToggleFlag(t *testing.T) {
	board := NewBoard(5, 5)

	// Should be able to flag an unrevealed cell
	if !board.ToggleFlag(0, 0) {
		t.Error("Should be able to flag unrevealed cell")
	}

	cell := board.GetCell(0, 0)
	if !cell.Flagged {
		t.Error("Cell should be flagged")
	}

	// Should be able to unflag
	if !board.ToggleFlag(0, 0) {
		t.Error("Should be able to unflag cell")
	}

	if cell.Flagged {
		t.Error("Cell should not be flagged")
	}

	// Should not be able to flag revealed cell
	cell.Revealed = true
	if board.ToggleFlag(0, 0) {
		t.Error("Should not be able to flag revealed cell")
	}
}

func TestGenerateBoardDeterminism(t *testing.T) {
	seed := int64(42)

	board1 := GenerateBoard(10, 10, 0.2, seed)
	board2 := GenerateBoard(10, 10, 0.2, seed)

	// Should generate identical boards with same seed
	for i := range board1.Cells {
		if board1.Cells[i].HasMonster != board2.Cells[i].HasMonster {
			t.Errorf("Boards with same seed should be identical at index %d", i)
		}
	}
}

func TestRevealCell(t *testing.T) {
	// Create a simple 3x3 board with one monster
	board := NewBoard(3, 3)
	board.Cells[0].HasMonster = true
	board.Monsters[0] = true
	board.RemainingSafe = 8
	board.CalculateThreatCounts()

	// Reveal a safe cell
	hitMonster, revealed := board.RevealCell(2, 2)
	if hitMonster {
		t.Error("Should not hit monster")
	}
	if revealed == 0 {
		t.Error("Should reveal at least one cell")
	}

	// Reveal the monster
	hitMonster, _ = board.RevealCell(0, 0)
	if !hitMonster {
		t.Error("Should hit monster")
	}
}

func TestIsWon(t *testing.T) {
	board := NewBoard(3, 3)

	// Initially not won (cells not revealed)
	if board.IsWon() {
		t.Error("Game should not be won initially")
	}

	// Reveal all safe cells (RemainingSafe = 0)
	board.RemainingSafe = 0
	if !board.IsWon() {
		t.Error("Game should be won when RemainingSafe = 0")
	}
}
