package game

import (
	"math/rand"
	"testing"
)

func TestPickPowerUpChoices(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	choices := PickPowerUpChoices(rng, 3, nil)
	if len(choices) != 3 {
		t.Errorf("Expected 3 choices, got %d", len(choices))
	}

	// All choices should be unique
	seen := map[PowerUpID]bool{}
	for _, c := range choices {
		if seen[c.ID] {
			t.Errorf("Duplicate power-up: %s", c.ID)
		}
		seen[c.ID] = true
	}
}

func TestPickPowerUpChoicesDeterministic(t *testing.T) {
	rng1 := rand.New(rand.NewSource(42))
	rng2 := rand.New(rand.NewSource(42))

	choices1 := PickPowerUpChoices(rng1, 3, nil)
	choices2 := PickPowerUpChoices(rng2, 3, nil)

	for i := range choices1 {
		if choices1[i].ID != choices2[i].ID {
			t.Errorf("Same seed should produce same choices: %s vs %s", choices1[i].ID, choices2[i].ID)
		}
	}
}

func TestPickPowerUpChoicesExcludes(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	// Exclude some power-ups
	exclude := []PowerUpID{PowerUpScout, PowerUpDowsing, PowerUpXRay, PowerUpLuckyCharm, PowerUpSwiftBoots, PowerUpFlagMaster}
	choices := PickPowerUpChoices(rng, 3, exclude)

	// Should only get shield and second wind (the re-acquirable ones) or remaining
	for _, c := range choices {
		for _, excl := range exclude {
			if c.ID == excl {
				t.Errorf("Should not offer excluded power-up: %s", c.ID)
			}
		}
	}
}

func TestAddPowerUp(t *testing.T) {
	run := NewRun(42, true, false)

	run.AddPowerUp(PowerUpDef{ID: PowerUpScout, Name: "Scout", Rarity: Common})
	if len(run.PowerUps) != 1 {
		t.Errorf("Expected 1 power-up, got %d", len(run.PowerUps))
	}
	if run.PowerUps[0].ID != PowerUpScout {
		t.Errorf("Expected Scout, got %s", run.PowerUps[0].ID)
	}
	if run.PowerUps[0].Uses != -1 {
		t.Errorf("Scout should have unlimited uses, got %d", run.PowerUps[0].Uses)
	}

	// Shield should have 1 use
	run.AddPowerUp(PowerUpDef{ID: PowerUpShield, Name: "Shield", Rarity: Rare})
	if run.PowerUps[1].Uses != 1 {
		t.Errorf("Shield should have 1 use, got %d", run.PowerUps[1].Uses)
	}
}

func TestHasPowerUp(t *testing.T) {
	run := NewRun(42, true, false)
	run.AddPowerUp(PowerUpDef{ID: PowerUpScout, Name: "Scout", Rarity: Common})

	if !run.HasPowerUp(PowerUpScout) {
		t.Error("Should have Scout power-up")
	}
	if run.HasPowerUp(PowerUpShield) {
		t.Error("Should not have Shield power-up")
	}

	// Consumed power-up should not count
	run.PowerUps[0].Consumed = true
	if run.HasPowerUp(PowerUpScout) {
		t.Error("Consumed power-up should not count")
	}
}

func TestPowerUpNames(t *testing.T) {
	run := NewRun(42, true, false)
	run.AddPowerUp(PowerUpDef{ID: PowerUpScout, Name: "Scout", Rarity: Common})
	run.AddPowerUp(PowerUpDef{ID: PowerUpShield, Name: "Shield", Rarity: Rare})

	names := run.PowerUpNames()
	if len(names) != 2 {
		t.Errorf("Expected 2 names, got %d", len(names))
	}

	// Consume one
	run.PowerUps[1].Consumed = true
	names = run.PowerUpNames()
	if len(names) != 1 {
		t.Errorf("Expected 1 name after consumption, got %d", len(names))
	}
}

func TestHandleDeathPowerUpsShield(t *testing.T) {
	run := NewRun(42, true, false)
	run.StartNextFloor()
	run.AddPowerUp(PowerUpDef{ID: PowerUpShield, Name: "Shield", Rarity: Rare})

	// Simulate death
	g := run.CurrentGame
	g.State = "lost"

	prevented := HandleDeathPowerUps(run)
	if !prevented {
		t.Error("Shield should prevent death")
	}
	if g.State != "playing" {
		t.Errorf("Game state should be playing, got %s", g.State)
	}
	if !run.PowerUps[0].Consumed {
		t.Error("Shield should be consumed after use")
	}

	// Second death should not be prevented
	g.State = "lost"
	prevented = HandleDeathPowerUps(run)
	if prevented {
		t.Error("Consumed shield should not prevent death")
	}
}

func TestApplyScout(t *testing.T) {
	run := NewRun(42, true, false)
	run.AddPowerUp(PowerUpDef{ID: PowerUpScout, Name: "Scout", Rarity: Common})
	run.StartNextFloor()

	g := run.CurrentGame
	revealedCount := 0
	for _, c := range g.Board.Cells {
		if c.Revealed {
			revealedCount++
		}
	}

	// Scout should have revealed some cells at floor start
	if revealedCount == 0 {
		t.Error("Scout should reveal cells at floor start")
	}
}

func TestApplyXRay(t *testing.T) {
	run := NewRun(42, true, false)
	run.AddPowerUp(PowerUpDef{ID: PowerUpXRay, Name: "X-Ray", Rarity: Uncommon})
	run.StartNextFloor()

	g := run.CurrentGame
	b := g.Board

	// Check that edge cells without monsters are revealed
	for x := 0; x < b.Width; x++ {
		cell := b.GetCell(x, 0)
		if !cell.HasMonster && !cell.Revealed {
			t.Errorf("X-Ray should reveal safe edge cell at (%d, 0)", x)
		}
	}
}

func TestLuckyCharmReducesDensity(t *testing.T) {
	// Run without Lucky Charm
	run1 := NewRun(42, true, false)
	run1.StartNextFloor()
	monsters1 := len(run1.CurrentGame.Board.Monsters)

	// Run with Lucky Charm (need to start from same floor for comparison)
	run2 := NewRun(42, true, false)
	run2.AddPowerUp(PowerUpDef{ID: PowerUpLuckyCharm, Name: "Lucky Charm", Rarity: Uncommon})
	run2.StartNextFloor()
	monsters2 := len(run2.CurrentGame.Board.Monsters)

	// Lucky Charm should result in fewer monsters
	if monsters2 >= monsters1 {
		t.Errorf("Lucky Charm should reduce monsters: %d >= %d", monsters2, monsters1)
	}
}

func TestHandleDeathPowerUpsTimeWarp(t *testing.T) {
	run := NewRun(42, true, false)
	run.StartNextFloor()
	run.AddPowerUp(PowerUpDef{ID: PowerUpTimeWarp, Name: "Time Warp", Rarity: Rare})

	// Save snapshot before "reveal"
	run.SaveSnapshot()

	// Simulate death
	g := run.CurrentGame
	g.State = "lost"

	prevented := HandleDeathPowerUps(run)
	if !prevented {
		t.Error("Time Warp should prevent death")
	}
	if g.State != "playing" {
		t.Errorf("Game state should be playing after Time Warp, got %s", g.State)
	}

	// Time Warp should NOT be consumed permanently (it recharges)
	if run.PowerUps[0].Consumed {
		t.Error("Time Warp should not be consumed permanently")
	}

	// But it should be used for this floor
	if run.UndoUsedFloor != run.FloorNum {
		t.Error("UndoUsedFloor should be set to current floor")
	}

	// Second death on same floor should not be prevented by Time Warp
	run.SaveSnapshot()
	g.State = "lost"
	prevented = HandleDeathPowerUps(run)
	if prevented {
		t.Error("Time Warp should not prevent death twice on same floor")
	}
}

func TestTimeWarpPrioritizedOverShield(t *testing.T) {
	run := NewRun(42, true, false)
	run.StartNextFloor()
	run.AddPowerUp(PowerUpDef{ID: PowerUpShield, Name: "Shield", Rarity: Rare})
	run.AddPowerUp(PowerUpDef{ID: PowerUpTimeWarp, Name: "Time Warp", Rarity: Rare})

	// Save snapshot
	run.SaveSnapshot()

	// Simulate death
	run.CurrentGame.State = "lost"
	HandleDeathPowerUps(run)

	// Time Warp should be used first (renewable), Shield should be untouched
	if run.PowerUps[0].Consumed {
		t.Error("Shield should not be consumed when Time Warp is available")
	}
	if run.UndoUsedFloor != run.FloorNum {
		t.Error("Time Warp should have been used")
	}
}

func TestRarityName(t *testing.T) {
	tests := []struct {
		r    Rarity
		want string
	}{
		{Common, "Common"},
		{Uncommon, "Uncommon"},
		{Rare, "Rare"},
		{Legendary, "Legendary"},
		{Rarity(99), "Unknown"},
	}
	for _, tt := range tests {
		got := RarityName(tt.r)
		if got != tt.want {
			t.Errorf("RarityName(%d) = %s, want %s", tt.r, got, tt.want)
		}
	}
}
