package game

import (
	"testing"
	"time"
)

func TestGetFloorConfig(t *testing.T) {
	// Floor 1 should be the smallest/easiest
	cfg1 := GetFloorConfig(1)
	if cfg1.Width != 12 || cfg1.Height != 8 {
		t.Errorf("Floor 1: expected 12x8, got %dx%d", cfg1.Width, cfg1.Height)
	}
	if cfg1.MonsterDensity != 0.12 {
		t.Errorf("Floor 1: expected density 0.12, got %f", cfg1.MonsterDensity)
	}

	// Floor 6 should be the last predefined
	cfg6 := GetFloorConfig(6)
	if cfg6.Width != 24 || cfg6.Height != 16 {
		t.Errorf("Floor 6: expected 24x16, got %dx%d", cfg6.Width, cfg6.Height)
	}
	if cfg6.MonsterDensity != 0.22 {
		t.Errorf("Floor 6: expected density 0.22, got %f", cfg6.MonsterDensity)
	}

	// Floor 7+ should scale beyond predefined
	cfg7 := GetFloorConfig(7)
	if cfg7.Width != 24 || cfg7.Height != 16 {
		t.Errorf("Floor 7: expected 24x16, got %dx%d", cfg7.Width, cfg7.Height)
	}
	if cfg7.MonsterDensity != 0.23 {
		t.Errorf("Floor 7: expected density 0.23, got %f", cfg7.MonsterDensity)
	}

	// Floor 0 and negative should default to floor 1
	cfgZero := GetFloorConfig(0)
	if cfgZero.Width != cfg1.Width {
		t.Errorf("Floor 0 should default to floor 1 config")
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
