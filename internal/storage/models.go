package storage

import "time"

// Player represents a known player (identified by SSH key fingerprint)
type Player struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
}

// Score represents a recorded score entry
type Score struct {
	PlayerID   string    `json:"player_id"`
	PlayerName string    `json:"player_name"`
	Date       string    `json:"date"`
	Score      int       `json:"score"`
	Floors     int       `json:"floors"`
	DurationMs int64     `json:"duration_ms"`
	Seed       int64     `json:"seed"`
	IsDaily    bool      `json:"is_daily"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// LeaderboardEntry is a score with rank information
type LeaderboardEntry struct {
	Rank       int    `json:"rank"`
	PlayerName string `json:"player_name"`
	Score      int    `json:"score"`
	Floors     int    `json:"floors"`
	IsYou      bool   `json:"is_you"`
}
