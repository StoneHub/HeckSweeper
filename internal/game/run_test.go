package game

import (
	"testing"
	"time"
)

func TestGetFloorConfig(t *testing.T) {
	// Floor 1 should be the tutorial floor (smallest/easiest)
	cfg1 := GetFloorConfig(1)
	if cfg1.Width != 8 || cfg1.Height != 6 {
		t.Errorf("Floor 1: expected 8x6, got %dx%d", cfg1.Width, cfg1.Height)
	}
	if cfg1.MonsterDensity != 0.10 {
		t.Errorf("Floor 1: expected density 0.10, got %f", cfg1.MonsterDensity)
	}

	// Floor 7 should be the last predefined
	cfg7 := GetFloorConfig(7)
	if cfg7.Width != 24 || cfg7.Height != 16 {
		t.Errorf("Floor 7: expected 24x16, got %dx%d", cfg7.Width, cfg7.Height)
	}
	if cfg7.MonsterDensity != 0.22 {
		t.Errorf("Floor 7: expected density 0.22, got %f", cfg7.MonsterDensity)
	}

	// Floor 8+ should scale beyond predefined
	cfg8 := GetFloorConfig(8)
	if cfg8.Width != 24 || cfg8.Height != 16 {
		t.Errorf("Floor 8: expected 24x16, got %dx%d", cfg8.Width, cfg8.Height)
	}
	if cfg8.MonsterDensity != 0.23 {
		t.Errorf("Floor 8: expected density 0.23, got %f", cfg8.MonsterDensity)
	}

	// Floor 0 and negative should default to floor 1
	cfgZero := GetFloorConfig(0)
	if cfgZero.Width != cfg1.Width {
		t.Errorf("Floor 0 should default to floor 1 config")
	}
}

func TestSaveAndRestoreSnapshot(t *testing.T) {
	run := NewRun(42, true, false)
	run.StartNextFloor()

	g := run.CurrentGame
	originalRemaining := g.Board.RemainingSafe
	originalCursor := g.CursorPos
	originalMoveCount := g.MoveCount

	// Save snapshot
	run.SaveSnapshot()

	// Modify game state
	g.CursorPos = Position{X: 0, Y: 0}
	g.Board.RemainingSafe = 0
	g.MoveCount = 99

	// Restore
	run.RestoreSnapshot()

	if g.Board.RemainingSafe != originalRemaining {
		t.Errorf("RemainingSafe should be restored: got %d, want %d", g.Board.RemainingSafe, originalRemaining)
	}
	if g.CursorPos != originalCursor {
		t.Errorf("CursorPos should be restored")
	}
	if g.MoveCount != originalMoveCount {
		t.Errorf("MoveCount should be restored to %d, got %d", originalMoveCount, g.MoveCount)
	}
	if run.LastSnapshot != nil {
		t.Error("LastSnapshot should be nil after restore")
	}
}

func TestGetUnrevealedMonsterPositions(t *testing.T) {
	run := NewRun(42, true, false)
	run.StartNextFloor()

	g := run.CurrentGame
	totalMonsters := len(g.Board.Monsters)

	// Before any reveals, all monsters should be unrevealed
	positions := g.Board.GetUnrevealedMonsterPositions(0, 0)
	if len(positions) != totalMonsters {
		t.Errorf("Expected %d unrevealed monsters, got %d", totalMonsters, len(positions))
	}

	// Reveal all monsters
	g.Board.RevealAllMonsters()
	positions = g.Board.GetUnrevealedMonsterPositions(0, 0)
	if len(positions) != 0 {
		t.Errorf("Expected 0 unrevealed monsters after reveal, got %d", len(positions))
	}
}

func TestGetFloorConfigDensityCap(t *testing.T) {
	// Very high floor should cap density at 0.35
	cfg := GetFloorConfig(100)
	if cfg.MonsterDensity > 0.35 {
		t.Errorf("Density should cap at 0.35, got %f", cfg.MonsterDensity)
	}
}

func TestNewRun(t *testing.T) {
	run := NewRun(42, true, false)

	if run.Seed != 42 {
		t.Errorf("Expected seed 42, got %d", run.Seed)
	}
	if run.FloorNum != 0 {
		t.Errorf("Expected floor 0, got %d", run.FloorNum)
	}
	if run.TotalScore != 0 {
		t.Errorf("Expected score 0, got %d", run.TotalScore)
	}
	if run.IsDaily {
		t.Error("Expected non-daily run")
	}
	if !run.UseUnicode {
		t.Error("Expected unicode enabled")
	}
}

func TestStartNextFloor(t *testing.T) {
	run := NewRun(42, true, false)

	run.StartNextFloor()
	if run.FloorNum != 1 {
		t.Errorf("Expected floor 1, got %d", run.FloorNum)
	}
	if run.CurrentGame == nil {
		t.Fatal("CurrentGame should not be nil after StartNextFloor")
	}

	// Board should match floor 1 config
	cfg := GetFloorConfig(1)
	if run.CurrentGame.Board.Width != cfg.Width {
		t.Errorf("Expected width %d, got %d", cfg.Width, run.CurrentGame.Board.Width)
	}
	if run.CurrentGame.Board.Height != cfg.Height {
		t.Errorf("Expected height %d, got %d", cfg.Height, run.CurrentGame.Board.Height)
	}

	// Advance to floor 2
	run.StartNextFloor()
	if run.FloorNum != 2 {
		t.Errorf("Expected floor 2, got %d", run.FloorNum)
	}
	cfg2 := GetFloorConfig(2)
	if run.CurrentGame.Board.Width != cfg2.Width {
		t.Errorf("Expected width %d, got %d", cfg2.Width, run.CurrentGame.Board.Width)
	}
}

func TestFloorSeedDeterminism(t *testing.T) {
	// Same seed should produce identical floors
	run1 := NewRun(42, true, false)
	run1.StartNextFloor()

	run2 := NewRun(42, true, false)
	run2.StartNextFloor()

	for i := range run1.CurrentGame.Board.Cells {
		if run1.CurrentGame.Board.Cells[i].HasMonster != run2.CurrentGame.Board.Cells[i].HasMonster {
			t.Errorf("Runs with same seed should produce identical floors at cell %d", i)
		}
	}

	// Different seeds should produce different floors
	run3 := NewRun(99, true, false)
	run3.StartNextFloor()

	different := false
	for i := range run1.CurrentGame.Board.Cells {
		if run1.CurrentGame.Board.Cells[i].HasMonster != run3.CurrentGame.Board.Cells[i].HasMonster {
			different = true
			break
		}
	}
	if !different {
		t.Error("Different seeds should produce different floors")
	}
}

func TestCalculateFloorScore(t *testing.T) {
	// Score should be positive for revealed cells
	score := CalculateFloorScore(50, 60, 30*time.Second, 1)
	if score <= 0 {
		t.Errorf("Score should be positive, got %d", score)
	}

	// Higher floor should give higher score (same stats)
	score1 := CalculateFloorScore(50, 60, 30*time.Second, 1)
	score3 := CalculateFloorScore(50, 60, 30*time.Second, 3)
	if score3 <= score1 {
		t.Errorf("Floor 3 score (%d) should be > floor 1 score (%d)", score3, score1)
	}

	// Faster time should give higher score
	scoreFast := CalculateFloorScore(50, 60, 10*time.Second, 1)
	scoreSlow := CalculateFloorScore(50, 60, 100*time.Second, 1)
	if scoreFast <= scoreSlow {
		t.Errorf("Fast score (%d) should be > slow score (%d)", scoreFast, scoreSlow)
	}

	// More revealed cells should give higher score
	scoreMore := CalculateFloorScore(80, 60, 30*time.Second, 1)
	scoreLess := CalculateFloorScore(20, 60, 30*time.Second, 1)
	if scoreMore <= scoreLess {
		t.Errorf("More cells score (%d) should be > fewer cells score (%d)", scoreMore, scoreLess)
	}
}

func TestTotalFloorsCleared(t *testing.T) {
	run := NewRun(42, true, false)

	run.FloorHistory = []FloorResult{
		{Floor: 1, Cleared: true, Score: 100},
		{Floor: 2, Cleared: true, Score: 200},
		{Floor: 3, Cleared: false, Score: 50},
	}

	if run.TotalFloorsCleared() != 2 {
		t.Errorf("Expected 2 floors cleared, got %d", run.TotalFloorsCleared())
	}
}

func TestRunIsActive(t *testing.T) {
	run := NewRun(42, true, false)

	// Before starting, run should be active
	if !run.IsActive() {
		t.Error("New run should be active")
	}

	// After starting floor, should be active
	run.StartNextFloor()
	if !run.IsActive() {
		t.Error("Run with playing game should be active")
	}

	// After dying, should not be active
	run.CurrentGame.State = "lost"
	if run.IsActive() {
		t.Error("Run with lost game should not be active")
	}
}
