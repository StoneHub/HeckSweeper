package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func tempStorePath(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return filepath.Join(dir, "test_data.json")
}

func TestNewStore(t *testing.T) {
	path := tempStorePath(t)
	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("NewStore failed: %v", err)
	}
	if store == nil {
		t.Fatal("Store should not be nil")
	}
}

func TestRecordPlayer(t *testing.T) {
	path := tempStorePath(t)
	store, _ := NewStore(path)

	store.RecordPlayer("player1", "Alice")

	// Verify data persisted
	store2, _ := NewStore(path)
	if _, ok := store2.data.Players["player1"]; !ok {
		t.Error("Player should be persisted")
	}
	if store2.data.Players["player1"].DisplayName != "Alice" {
		t.Error("Player name should be Alice")
	}
}

func TestSubmitScore(t *testing.T) {
	path := tempStorePath(t)
	store, _ := NewStore(path)

	score := Score{
		PlayerID:   "player1",
		PlayerName: "Alice",
		Date:       "2026-02-06",
		Score:      1000,
		Floors:     3,
		IsDaily:    false,
	}
	store.SubmitScore(score)

	if len(store.data.Scores) != 1 {
		t.Errorf("Expected 1 score, got %d", len(store.data.Scores))
	}
}

func TestDailyScoreBestOnly(t *testing.T) {
	path := tempStorePath(t)
	store, _ := NewStore(path)

	// Submit first daily score
	store.SubmitScore(Score{
		PlayerID: "player1",
		Date:     "2026-02-06",
		Score:    500,
		IsDaily:  true,
	})

	// Submit better daily score (same player, same date)
	store.SubmitScore(Score{
		PlayerID: "player1",
		Date:     "2026-02-06",
		Score:    800,
		IsDaily:  true,
	})

	// Should only keep the better score
	entries := store.GetDailyLeaderboard("2026-02-06", 10, "")
	if len(entries) != 1 {
		t.Errorf("Expected 1 entry (best only), got %d", len(entries))
	}
	if entries[0].Score != 800 {
		t.Errorf("Expected score 800 (best), got %d", entries[0].Score)
	}

	// Submit worse score — should not replace
	store.SubmitScore(Score{
		PlayerID: "player1",
		Date:     "2026-02-06",
		Score:    300,
		IsDaily:  true,
	})
	entries = store.GetDailyLeaderboard("2026-02-06", 10, "")
	if entries[0].Score != 800 {
		t.Errorf("Worse score should not replace best: got %d", entries[0].Score)
	}
}

func TestLeaderboardRanking(t *testing.T) {
	path := tempStorePath(t)
	store, _ := NewStore(path)

	store.SubmitScore(Score{PlayerID: "p1", PlayerName: "Alice", Date: "2026-02-06", Score: 500, Floors: 2, IsDaily: true})
	store.SubmitScore(Score{PlayerID: "p2", PlayerName: "Bob", Date: "2026-02-06", Score: 1000, Floors: 4, IsDaily: true})
	store.SubmitScore(Score{PlayerID: "p3", PlayerName: "Charlie", Date: "2026-02-06", Score: 750, Floors: 3, IsDaily: true})

	entries := store.GetDailyLeaderboard("2026-02-06", 10, "p3")

	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Should be sorted by score descending
	if entries[0].PlayerName != "Bob" || entries[0].Score != 1000 || entries[0].Rank != 1 {
		t.Errorf("Rank 1 should be Bob with 1000, got %s with %d", entries[0].PlayerName, entries[0].Score)
	}
	if entries[1].PlayerName != "Charlie" || entries[1].Rank != 2 {
		t.Errorf("Rank 2 should be Charlie, got %s", entries[1].PlayerName)
	}
	if entries[2].PlayerName != "Alice" || entries[2].Rank != 3 {
		t.Errorf("Rank 3 should be Alice, got %s", entries[2].PlayerName)
	}

	// Charlie should be marked as "you"
	if !entries[1].IsYou {
		t.Error("Charlie should be marked as IsYou")
	}
	if entries[0].IsYou || entries[2].IsYou {
		t.Error("Others should not be IsYou")
	}
}

func TestLeaderboardLimit(t *testing.T) {
	path := tempStorePath(t)
	store, _ := NewStore(path)

	for i := 0; i < 20; i++ {
		store.SubmitScore(Score{
			PlayerID: "p" + string(rune('a'+i)),
			Date:     "2026-02-06",
			Score:    i * 100,
			IsDaily:  true,
		})
	}

	entries := store.GetDailyLeaderboard("2026-02-06", 5, "")
	if len(entries) != 5 {
		t.Errorf("Expected 5 entries (limited), got %d", len(entries))
	}
}

func TestGetPlayerCount(t *testing.T) {
	path := tempStorePath(t)
	store, _ := NewStore(path)

	store.SubmitScore(Score{PlayerID: "p1", Date: "2026-02-06", IsDaily: true})
	store.SubmitScore(Score{PlayerID: "p2", Date: "2026-02-06", IsDaily: true})
	store.SubmitScore(Score{PlayerID: "p3", Date: "2026-02-07", IsDaily: true})

	count := store.GetPlayerCount("2026-02-06")
	if count != 2 {
		t.Errorf("Expected 2 players for 2026-02-06, got %d", count)
	}
}

func TestStorePersistence(t *testing.T) {
	path := tempStorePath(t)

	// Create store, add data
	store1, _ := NewStore(path)
	store1.SubmitScore(Score{PlayerID: "p1", Score: 500})

	// Reload from disk
	store2, _ := NewStore(path)
	if len(store2.data.Scores) != 1 {
		t.Error("Data should persist across store instances")
	}
	if store2.data.Scores[0].Score != 500 {
		t.Error("Score should persist correctly")
	}
}

func TestStoreEmptyFile(t *testing.T) {
	path := tempStorePath(t)

	// Create empty file
	os.WriteFile(path, []byte{}, 0644)

	store, err := NewStore(path)
	if err != nil {
		t.Fatalf("Should handle empty file: %v", err)
	}
	if store == nil {
		t.Fatal("Store should not be nil")
	}
}
