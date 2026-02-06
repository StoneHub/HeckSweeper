package storage

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
	"time"
)

// Store provides thread-safe persistent storage for scores and players
type Store struct {
	mu       sync.RWMutex
	filePath string
	data     storeData
}

type storeData struct {
	Players map[string]Player `json:"players"`
	Scores  []Score           `json:"scores"`
}

// NewStore creates or loads a store from disk
func NewStore(filePath string) (*Store, error) {
	s := &Store{
		filePath: filePath,
		data: storeData{
			Players: make(map[string]Player),
			Scores:  []Score{},
		},
	}

	// Try loading existing data
	if _, err := os.Stat(filePath); err == nil {
		raw, err := os.ReadFile(filePath)
		if err != nil {
			return nil, err
		}
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &s.data); err != nil {
				return nil, err
			}
		}
	}

	if s.data.Players == nil {
		s.data.Players = make(map[string]Player)
	}

	return s, nil
}

// save persists the store to disk (caller must hold lock)
func (s *Store) save() error {
	raw, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, raw, 0644)
}

// RecordPlayer creates or updates a player record
func (s *Store) RecordPlayer(id, displayName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	p, exists := s.data.Players[id]
	if !exists {
		p = Player{
			ID:        id,
			FirstSeen: now,
		}
	}
	if displayName != "" {
		p.DisplayName = displayName
	}
	p.LastSeen = now
	s.data.Players[id] = p
	_ = s.save()
}

// SubmitScore records a score. For daily challenges, only the best score per player per day is kept.
func (s *Store) SubmitScore(score Score) {
	s.mu.Lock()
	defer s.mu.Unlock()

	score.SubmittedAt = time.Now()

	if score.IsDaily {
		// Check for existing daily score for this player+date
		for i, existing := range s.data.Scores {
			if existing.IsDaily && existing.PlayerID == score.PlayerID && existing.Date == score.Date {
				// Keep the higher score
				if score.Score > existing.Score {
					s.data.Scores[i] = score
					_ = s.save()
				}
				return
			}
		}
	}

	s.data.Scores = append(s.data.Scores, score)
	_ = s.save()
}

// GetDailyLeaderboard returns the top N scores for a specific date
func (s *Store) GetDailyLeaderboard(date string, limit int, currentPlayerID string) []LeaderboardEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Filter to daily scores for this date
	var dailyScores []Score
	for _, sc := range s.data.Scores {
		if sc.IsDaily && sc.Date == date {
			dailyScores = append(dailyScores, sc)
		}
	}

	return s.buildLeaderboard(dailyScores, limit, currentPlayerID)
}

// GetAllTimeLeaderboard returns the top N scores overall
func (s *Store) GetAllTimeLeaderboard(limit int, currentPlayerID string) []LeaderboardEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.buildLeaderboard(s.data.Scores, limit, currentPlayerID)
}

// buildLeaderboard sorts and formats scores into leaderboard entries
func (s *Store) buildLeaderboard(scores []Score, limit int, currentPlayerID string) []LeaderboardEntry {
	// Sort by score descending
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	if limit > 0 && len(scores) > limit {
		scores = scores[:limit]
	}

	entries := make([]LeaderboardEntry, len(scores))
	for i, sc := range scores {
		name := sc.PlayerName
		if name == "" {
			// Use truncated player ID
			if len(sc.PlayerID) > 12 {
				name = sc.PlayerID[:12] + "..."
			} else {
				name = sc.PlayerID
			}
		}
		if name == "" {
			name = "Anonymous"
		}

		entries[i] = LeaderboardEntry{
			Rank:       i + 1,
			PlayerName: name,
			Score:      sc.Score,
			Floors:     sc.Floors,
			IsYou:      sc.PlayerID == currentPlayerID && currentPlayerID != "",
		}
	}

	return entries
}

// GetPlayerCount returns the total number of players who have submitted scores for a date
func (s *Store) GetPlayerCount(date string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	players := make(map[string]bool)
	for _, sc := range s.data.Scores {
		if sc.IsDaily && sc.Date == date {
			players[sc.PlayerID] = true
		}
	}
	return len(players)
}
