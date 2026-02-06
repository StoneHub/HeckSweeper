# HeckSweeper Revival Plan: Roguelite + SSH Multiplayer

> **Goal**: Transform HeckSweeper from a vanilla terminal minesweeper into a roguelite
> dungeon crawler with progression, daily challenges, and SSH-accessible multiplayer —
> all served over `ssh play@hecksweeper.com`.

---

## Table of Contents

1. [Architecture Overview](#1-architecture-overview)
2. [Phase 1 — Roguelite Floor System](#2-phase-1--roguelite-floor-system)
3. [Phase 2 — Power-Up System](#3-phase-2--power-up-system)
4. [Phase 3 — Wish SSH Server](#4-phase-3--wish-ssh-server)
5. [Phase 4 — Daily Challenges & Leaderboards](#5-phase-4--daily-challenges--leaderboards)
6. [Phase 5 — Polish & Distribution](#6-phase-5--polish--distribution)
7. [File Map](#7-file-map)
8. [Data Models](#8-data-models)
9. [Technical Decisions](#9-technical-decisions)

---

## 1. Architecture Overview

### Current State (what we have)

```
cmd/hecksweeper/main.go      → CLI entry point, flag parsing
internal/game/types.go       → Cell, Board, Position, Game structs
internal/game/board.go       → Board logic, BFS flood-fill, cursor, reveal
internal/game/generator.go   → Seeded board generation with exclusion zones
internal/ui/model.go         → Bubbletea Model (wraps Game)
internal/ui/update.go        → Key handling (movement, reveal, flag)
internal/ui/view.go          → Board rendering, status bar, game-over screen
internal/ui/styles.go        → Lipgloss color palette and style definitions
internal/constants/           → Glyphs, dimensions, game states
```

The game logic (`internal/game/`) is cleanly separated from the UI (`internal/ui/`).
This is the key property that makes the revival feasible — we can layer roguelite
mechanics on top of the game layer without touching the renderer, and we can wrap
the entire UI in a Wish SSH handler without touching the game layer.

### Target State (what we're building)

```
cmd/
  hecksweeper/main.go        → Local play entry point (unchanged for local)
  server/main.go             → Wish SSH server entry point (NEW)

internal/
  game/
    types.go                 → Extended with Run, Floor, PowerUp types
    board.go                 → Unchanged (single-floor logic stays here)
    generator.go             → Extended with floor-scaled difficulty
    run.go                   → NEW: Multi-floor run manager
    powerups.go              → NEW: Power-up definitions and effects
  ui/
    model.go                 → Extended with run state, floor transitions
    update.go                → Extended with power-up selection, floor nav
    view.go                  → Extended with floor header, power-up picker
    views_powerup.go         → NEW: Power-up selection screen renderer
    views_title.go           → NEW: Title screen / main menu renderer
    views_summary.go         → NEW: Run summary / death screen renderer
    styles.go                → Extended with new screen styles
  server/
    server.go                → NEW: Wish SSH server setup
    handler.go               → NEW: Bubbletea middleware handler
    auth.go                  → NEW: Public key + anonymous auth
  daily/
    daily.go                 → NEW: Daily seed generation (date-based)
    leaderboard.go           → NEW: Score tracking and ranking
  storage/
    store.go                 → NEW: Storage interface
    sqlite.go                → NEW: SQLite implementation
    models.go                → NEW: DB models (scores, runs, players)
```

### Dependency Additions

```go
// go.mod additions
require (
    github.com/charmbracelet/wish   // SSH server framework
    github.com/charmbracelet/ssh    // SSH primitives
    github.com/gliderlabs/ssh       // Underlying SSH (wish dep)
    modernc.org/sqlite              // Pure-Go SQLite (no CGO)
)
```

We use `modernc.org/sqlite` instead of `mattn/go-sqlite3` to keep the binary
fully static with zero CGO dependency. This matters for cross-compilation and
single-binary deployment.

---

## 2. Phase 1 — Roguelite Floor System

**The core idea**: Each "run" is a series of minesweeper floors. Clear a floor →
pick a power-up → descend to a harder floor. Die → run ends, see your stats.

### 2.1 Run State Machine

```
[Title Screen] → [Floor N] → [Floor Cleared] → [Power-Up Pick] → [Floor N+1]
                      ↓                                               ↓
                 [Death Screen] ←──────────────────────────────────────┘
                      ↓
               [Run Summary]
                      ↓
               [Title Screen]
```

Game states to add to `constants.go`:

```go
GameStateTitle     = "title"      // Main menu
GameStateFloorIntro = "floor_intro" // "Descending to Floor 3..."
GameStatePlaying   = "playing"    // Existing — the minesweeper board
GameStateCleared   = "cleared"    // Floor cleared, show stats
GameStatePowerUp   = "powerup"   // Choosing a power-up
GameStateDead      = "dead"       // Hit a monster, run over
GameStateSummary   = "summary"    // Final run stats
```

### 2.2 Floor Scaling

Each floor increases difficulty. The parameters that scale:

| Floor | Board Size | Monster Density | New Mechanic        |
|-------|-----------|-----------------|---------------------|
| 1     | 12×8      | 12%             | Tutorial-safe       |
| 2     | 16×10     | 14%             | —                   |
| 3     | 18×12     | 16%             | —                   |
| 4     | 20×14     | 18%             | Armored monsters*   |
| 5     | 22×14     | 20%             | —                   |
| 6     | 24×16     | 22%             | Fog of war*         |
| 7+    | 24×16     | 22% + 1%/floor  | Stacking difficulty |

*Armored monsters and fog of war are stretch goals for later phases. The initial
implementation just scales size and density.

Floor parameters live in `internal/game/run.go`:

```go
type FloorConfig struct {
    Width          int
    Height         int
    MonsterDensity float64
    FloorNumber    int
}

func GetFloorConfig(floor int) FloorConfig {
    // Returns scaled parameters for the given floor number
}
```

### 2.3 Run Manager

The `Run` struct manages the progression through floors:

```go
type Run struct {
    Seed        int64
    FloorNum    int
    Score       int        // Cumulative across floors
    PowerUps    []PowerUp  // Collected power-ups
    CurrentGame *Game      // Active floor's game state
    FloorStats  []FloorResult
    StartedAt   time.Time
}

type FloorResult struct {
    Floor      int
    Moves      int
    Cleared    bool
    TimeSpent  time.Duration
}
```

### 2.4 Scoring

Points are earned per floor:

```
base_score       = cells_revealed * 10
speed_bonus      = max(0, (time_limit - time_spent) * 2)
efficiency_bonus = max(0, (optimal_moves - actual_moves) * 5)
floor_multiplier = floor_number * 1.5
floor_score      = (base_score + speed_bonus + efficiency_bonus) * floor_multiplier
```

The score is displayed in the stats bar and accumulated across the run.

### 2.5 Changes to Existing Files

**`internal/game/types.go`** — Add `Run`, `FloorConfig`, `FloorResult` structs.
The existing `Game` struct stays as-is (it represents one floor).

**`internal/ui/model.go`** — The `Model` struct gains a `run *game.Run` field.
The `game` field becomes `run.CurrentGame`. The `NewModel` constructor changes
to initialize a `Run` instead of a bare `Game`.

**`internal/ui/update.go`** — Key handling becomes state-dependent:
- `GameStateTitle`: Enter to start run, D for daily challenge, Q to quit
- `GameStatePlaying`: Existing controls (unchanged)
- `GameStateCleared`: Enter to continue to power-up selection
- `GameStatePowerUp`: 1/2/3 to pick a power-up
- `GameStateDead`/`GameStateSummary`: R to restart, Q to quit

**`internal/ui/view.go`** — The `View()` method dispatches to sub-renderers
based on state. The board renderer stays unchanged. New screens get their own
view files to keep things organized.

---

## 3. Phase 2 — Power-Up System

### 3.1 Power-Up Design

After clearing each floor, the player chooses 1 of 3 randomly offered power-ups.
Power-ups persist for the entire run.

| Power-Up         | Effect                                      | Rarity  |
|-----------------|----------------------------------------------|---------|
| **Scout**        | Reveal a random 3×3 safe area on each floor  | Common  |
| **Shield**       | Survive one monster hit (consumed on use)     | Rare    |
| **Dowsing Rod**  | First reveal on each floor is always safe     | Common  |
| **X-Ray**        | Permanently reveal all edge cells on each floor | Uncommon |
| **Lucky Charm**  | Reduce monster density by 2% for all future floors | Uncommon |
| **Swift Boots**  | +50% speed bonus multiplier                  | Common  |
| **Flag Master**  | Auto-flag cells when all neighbors are revealed | Common |
| **Second Wind**  | On death, revive once per run with board reset | Legendary |

### 3.2 Power-Up Data Model

```go
type PowerUp struct {
    ID          string
    Name        string
    Description string
    Rarity      Rarity  // Common, Uncommon, Rare, Legendary
    Apply       func(run *Run, game *Game)  // Called at floor start
    OnReveal    func(run *Run, game *Game, x, y int) bool  // Hook into reveal
    OnDeath     func(run *Run, game *Game) bool  // Return true to prevent death
}
```

Power-ups hook into the game at specific points:
- **Floor start**: Scout, X-Ray, Dowsing Rod activate here
- **On reveal**: Flag Master triggers after each reveal
- **On death**: Shield, Second Wind intercept the death

### 3.3 Power-Up Selection Screen

After clearing a floor, the player sees:

```
╔═══════════════════════════════════════╗
║         FLOOR 3 CLEARED!             ║
║   Score: 1,250  |  Time: 0:42       ║
╚═══════════════════════════════════════╝

  Choose a power-up:

  [1] Scout (Common)
      Reveal a random 3×3 safe area on each floor

  [2] Shield (Rare)
      Survive one monster hit (consumed on use)

  [3] Lucky Charm (Uncommon)
      Reduce monster density by 2% for future floors

  Press 1, 2, or 3 to choose
```

The three offered power-ups are randomly selected (weighted by rarity) from the
pool, using the run's RNG seed so daily challenge players get the same choices.

### 3.4 Implementation

New file `internal/game/powerups.go`:
- Power-up registry (all power-ups defined as data)
- Random selection with rarity weighting
- Application functions for each power-up type

Changes to `internal/game/board.go`:
- `RevealCell` gains a hook point for power-up effects (Shield intercept)
- The function signature stays the same; hooks are checked via the `Run` struct

Changes to `internal/game/run.go`:
- `Run.ApplyFloorStartPowerUps()` — called when a new floor begins
- `Run.HandleDeath()` — checks for death-prevention power-ups before confirming death

---

## 4. Phase 3 — Wish SSH Server

### 4.1 Overview

Wish wraps our existing bubbletea `Model` in an SSH server. Each SSH connection
gets its own `tea.Program` with isolated state. The game works identically to
local play — the only difference is the transport layer.

### 4.2 Server Entry Point

New file `cmd/server/main.go`:

```go
func main() {
    // Parse server flags
    host := flag.String("host", "0.0.0.0", "Listen address")
    port := flag.Int("port", 2222, "SSH port")
    keyPath := flag.String("key", ".ssh/server_ed25519", "Host key path")

    // Create Wish server
    s, err := wish.NewServer(
        wish.WithAddress(net.JoinHostPort(*host, strconv.Itoa(*port))),
        wish.WithHostKeyPath(*keyPath),
        wish.WithMiddleware(
            bubbletea.Middleware(newHandler),  // Our game handler
            activeterm.Middleware(),            // Ensure active terminal
            logging.Middleware(),               // Connection logging
        ),
    )

    // Start with graceful shutdown
    // ...
}
```

### 4.3 Handler

New file `internal/server/handler.go`:

```go
func newHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
    pty, _, _ := sess.Pty()
    // Create a new game model sized to the client's terminal
    width := min(pty.Window.Width/2, 24)   // 2 chars per cell
    height := min(pty.Window.Height-8, 16) // Leave room for chrome
    seed := time.Now().UnixNano()

    model := ui.NewModel(width, height, seed, true)
    return model, []tea.ProgramOption{tea.WithAltScreen()}
}
```

Each connection gets a fresh game. Terminal dimensions from the SSH pty are used
to auto-size the board — no configuration needed from the player.

### 4.4 Authentication

For a game, we want low friction. The auth strategy:

1. **Anonymous access** (default): Anyone can `ssh -p 2222 host` and play.
   No username/password required.
2. **Optional public key**: If a player connects with a key, we can associate
   their scores with a persistent identity for leaderboards. The key fingerprint
   becomes their player ID.
3. **No passwords**: We never store or check passwords.

```go
func authHandler(ctx ssh.Context, key ssh.PublicKey) bool {
    // Always allow connection
    // If key is provided, store fingerprint for leaderboard identity
    if key != nil {
        fingerprint := gossh.FingerprintSHA256(key)
        ctx.SetValue("player_id", fingerprint)
    }
    return true
}
```

Players without keys are "anonymous" and can play but won't appear on persistent
leaderboards. This keeps the barrier to entry at zero while rewarding identity.

### 4.5 Session Management

Considerations for a public SSH game server:

- **Idle timeout**: Disconnect after 10 minutes of inactivity
- **Max concurrent sessions**: Configurable (default 100)
- **Rate limiting**: Max 5 connections per IP per minute
- **Resource limits**: Each session uses ~2-5 MB RAM (bubbletea + game state)
- **Graceful shutdown**: On SIGTERM, finish active games before closing

### 4.6 Deployment

The server is a single binary. Deployment options:

```bash
# Direct
./hecksweeper-server --port 2222 --key ./host_key

# systemd service
[Unit]
Description=HeckSweeper SSH Game Server
After=network.target

[Service]
ExecStart=/usr/local/bin/hecksweeper-server --port 2222
Restart=always
User=hecksweeper
NoNewPrivileges=true
ProtectSystem=strict

[Install]
WantedBy=multi-user.target
```

```bash
# Docker
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o server ./cmd/server

FROM alpine:latest
COPY --from=build /app/server /usr/local/bin/
EXPOSE 2222
CMD ["server"]
```

---

## 5. Phase 4 — Daily Challenges & Leaderboards

### 5.1 Daily Seed

The daily challenge uses a deterministic seed derived from the date:

```go
func DailySeed() int64 {
    now := time.Now().UTC()
    dateStr := now.Format("2006-01-02")
    h := fnv.New64a()
    h.Write([]byte("hecksweeper-daily-" + dateStr))
    return int64(h.Sum64())
}
```

All players on the same day get the same seed → same board layouts → same
power-up offerings → fair competition. The `"hecksweeper-daily-"` prefix is a
salt so the seeds aren't predictable from the date alone.

### 5.2 Leaderboard Storage

Using SQLite (pure-Go via `modernc.org/sqlite`) for persistence:

```sql
CREATE TABLE players (
    id          TEXT PRIMARY KEY,  -- SSH key fingerprint or "anon_<hash>"
    display_name TEXT,
    first_seen  DATETIME,
    last_seen   DATETIME
);

CREATE TABLE daily_scores (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id   TEXT REFERENCES players(id),
    date        TEXT,              -- "2026-02-06"
    score       INTEGER,
    floors      INTEGER,
    duration_ms INTEGER,
    submitted   DATETIME,
    UNIQUE(player_id, date)        -- One score per player per day
);

CREATE TABLE run_scores (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id   TEXT REFERENCES players(id),
    seed        INTEGER,
    score       INTEGER,
    floors      INTEGER,
    duration_ms INTEGER,
    submitted   DATETIME
);
```

### 5.3 Leaderboard Display

Accessible from the title screen (press `L`):

```
╔═══════════════════════════════════════╗
║      DAILY LEADERBOARD  2026-02-06   ║
╠═══════════════════════════════════════╣
║  #   Player          Score   Floors  ║
║  1.  ssh:a3f2...     4,850   7       ║
║  2.  ssh:9bc1...     3,200   5       ║
║  3.  ssh:f78d...     2,910   5       ║
║  4.  YOU →           2,100   4       ║
║  5.  ssh:12ab...     1,800   3       ║
╠═══════════════════════════════════════╣
║  Your best: 2,100 (Floor 4)         ║
║  Players today: 23                    ║
╚═══════════════════════════════════════╝

  [D] Daily  [A] All-time  [B] Back
```

Player identities are shown as truncated SSH key fingerprints. Players can
optionally set a display name via a command-line flag or environment variable
when connecting: `ssh -p 2222 host -o SetEnv=PLAYER_NAME=xXDungeonLordXx`

### 5.4 Anti-Cheat (Lightweight)

For a fun terminal game, we don't need serious anti-cheat. Simple measures:

- **Server-side validation**: The game runs on the server. Clients only send
  keystrokes. There's no client-side state to tamper with.
- **One daily attempt**: First completed run score is recorded. You can replay
  for fun but the leaderboard score doesn't update (unless it's higher — TBD).
- **Rate limiting**: Can't spam runs to brute-force optimal plays.

The SSH-served model inherently prevents most cheating since all game logic
executes server-side. This is a massive advantage over a client-side binary.

---

## 6. Phase 5 — Polish & Distribution

### 6.1 Title Screen

```
    ╔══════════════════════════════════════╗
    ║                                      ║
    ║   ██╗  ██╗███████╗ ██████╗██╗  ██╗  ║
    ║   ██║  ██║██╔════╝██╔════╝██║ ██╔╝  ║
    ║   ███████║█████╗  ██║     █████╔╝   ║
    ║   ██╔══██║██╔══╝  ██║     ██╔═██╗   ║
    ║   ██║  ██║███████╗╚██████╗██║  ██╗  ║
    ║   ╚═╝  ╚═╝╚══════╝ ╚═════╝╚═╝  ╚═╝  ║
    ║         S W E E P E R                ║
    ║                                      ║
    ║     A roguelite minesweeper          ║
    ║                                      ║
    ║   [N] New Run                        ║
    ║   [D] Daily Challenge                ║
    ║   [L] Leaderboard                    ║
    ║   [Q] Quit                           ║
    ║                                      ║
    ╚══════════════════════════════════════╝
```

### 6.2 Floor Transition Animation

Between floors, a brief transition screen:

```
    Descending to Floor 3...

    ░░░░░░░░░░░░░░░████████░░░░░░░░░░░░░░

    Board: 18×12  |  Monsters: 16%
    Power-ups: Scout, Shield
```

The progress bar animates using a bubbletea ticker. This gives a brief pause
between floors that builds tension and communicates what's coming.

### 6.3 Distribution

**Local play** (unchanged):
```bash
go install github.com/stonehub/hecksweeper/cmd/hecksweeper@latest
```

**SSH play**:
```bash
ssh -p 2222 play.hecksweeper.com
```

**Self-host the server**:
```bash
go install github.com/stonehub/hecksweeper/cmd/server@latest
hecksweeper-server --port 2222
```

### 6.4 Makefile Updates

Extend `Makefile.new` with server targets:

```makefile
server:
	go build -o hecksweeper-server ./cmd/server

run-server: server
	./hecksweeper-server --port 2222

docker:
	docker build -t hecksweeper .

deploy:
	# fly.io deployment
	flyctl deploy
```

---

## 7. File Map

Complete list of files to create or modify, by phase:

### Phase 1 — Roguelite Floors
| Action | File | Description |
|--------|------|-------------|
| CREATE | `internal/game/run.go` | Run manager, floor config, scoring |
| MODIFY | `internal/game/types.go` | Add Run, FloorConfig, FloorResult types |
| MODIFY | `internal/constants/constants.go` | Add new game states |
| MODIFY | `internal/ui/model.go` | Replace bare Game with Run |
| MODIFY | `internal/ui/update.go` | State-dependent key routing |
| MODIFY | `internal/ui/view.go` | Dispatch to sub-renderers |
| CREATE | `internal/ui/views_title.go` | Title screen renderer |
| CREATE | `internal/ui/views_summary.go` | Death/run summary renderer |
| CREATE | `internal/game/run_test.go` | Tests for floor scaling and scoring |

### Phase 2 — Power-Ups
| Action | File | Description |
|--------|------|-------------|
| CREATE | `internal/game/powerups.go` | Power-up registry and effects |
| CREATE | `internal/ui/views_powerup.go` | Power-up selection screen |
| MODIFY | `internal/game/board.go` | Hook points for power-up effects |
| MODIFY | `internal/game/run.go` | Power-up application at floor start |
| CREATE | `internal/game/powerups_test.go` | Tests for power-up mechanics |

### Phase 3 — Wish SSH
| Action | File | Description |
|--------|------|-------------|
| CREATE | `cmd/server/main.go` | SSH server entry point |
| CREATE | `internal/server/server.go` | Wish server setup |
| CREATE | `internal/server/handler.go` | Bubbletea SSH handler |
| CREATE | `internal/server/auth.go` | Auth (anonymous + pubkey) |
| MODIFY | `go.mod` | Add wish, ssh dependencies |
| MODIFY | `Makefile.new` | Add server build targets |

### Phase 4 — Daily & Leaderboards
| Action | File | Description |
|--------|------|-------------|
| CREATE | `internal/daily/daily.go` | Daily seed generation |
| CREATE | `internal/daily/leaderboard.go` | Score ranking logic |
| CREATE | `internal/storage/store.go` | Storage interface |
| CREATE | `internal/storage/sqlite.go` | SQLite implementation |
| CREATE | `internal/storage/models.go` | DB models |
| CREATE | `internal/ui/views_leaderboard.go` | Leaderboard renderer |
| MODIFY | `internal/ui/update.go` | Leaderboard navigation keys |
| CREATE | `internal/daily/daily_test.go` | Daily seed tests |

### Phase 5 — Polish
| Action | File | Description |
|--------|------|-------------|
| CREATE | `Dockerfile` | Container build |
| MODIFY | `Makefile.new` | Docker + deploy targets |
| MODIFY | `README.md` | Updated with SSH instructions |
| MODIFY | `internal/ui/styles.go` | New screen styles |

---

## 8. Data Models

### Complete Type Definitions (Phase 1+2)

```go
// internal/game/types.go additions

type Rarity int
const (
    Common    Rarity = iota
    Uncommon
    Rare
    Legendary
)

type FloorConfig struct {
    Width          int
    Height         int
    MonsterDensity float64
    FloorNumber    int
}

type FloorResult struct {
    Floor         int
    Moves         int
    CellsRevealed int
    TotalSafe     int
    Cleared       bool
    TimeSpent     time.Duration
    Score         int
}

type PowerUpID string

type ActivePowerUp struct {
    ID       PowerUpID
    Name     string
    Uses     int   // -1 = unlimited, 0 = expired
}

type Run struct {
    Seed         int64
    FloorNum     int
    TotalScore   int
    PowerUps     []ActivePowerUp
    CurrentGame  *Game
    FloorHistory []FloorResult
    StartedAt    time.Time
    IsDaily      bool
}
```

### Storage Models (Phase 4)

```go
// internal/storage/models.go

type Player struct {
    ID          string    // SSH fingerprint
    DisplayName string
    FirstSeen   time.Time
    LastSeen    time.Time
}

type DailyScore struct {
    PlayerID   string
    Date       string    // "2006-01-02"
    Score      int
    Floors     int
    DurationMs int64
}

type RunScore struct {
    PlayerID   string
    Seed       int64
    Score      int
    Floors     int
    DurationMs int64
}
```

---

## 9. Technical Decisions

### Why SQLite over flat files?

- Atomic writes (no corruption on crash)
- Query support for leaderboards (ORDER BY, LIMIT, date filtering)
- `modernc.org/sqlite` is pure Go — no CGO, fully static binary
- Single file (`hecksweeper.db`) — trivial backup and migration

### Why not a web frontend?

- SSH is the distribution channel. Players connect with a tool they already have.
- No HTTPS certificates, no CDN, no JavaScript bundle, no browser compatibility.
- Terminal rendering is already solved by bubbletea + lipgloss.
- The server-side execution model gives us anti-cheat for free.

### Why anonymous-first auth?

- Zero friction: `ssh -p 2222 host` and you're playing.
- Optional identity via SSH keys for those who want leaderboard persistence.
- No account creation, no passwords, no email verification.

### Why not WebSocket/browser instead of SSH?

- SSH gives us encrypted transport, terminal emulation, and window resize
  handling for free.
- The charm ecosystem (wish + bubbletea) is purpose-built for this.
- "ssh to play a game" is a compelling novelty that drives word-of-mouth.

### Board sizing for SSH sessions

The server reads the client's terminal dimensions from the SSH pty and
auto-scales the board. Minimum viable terminal: 80×24 (standard). The floor
config's width/height are treated as maximums — if the client terminal is
smaller, we clamp down. This ensures the game is playable on any reasonable
terminal.

---

## Implementation Order

The phases are designed to be independently valuable:

1. **Phase 1 (Roguelite Floors)**: Makes the game worth playing. This is the
   minimum viable revival — without progression, it's still just minesweeper.

2. **Phase 2 (Power-Ups)**: Adds replayability and strategic depth. Each run
   feels different because of the power-up choices.

3. **Phase 3 (Wish SSH)**: Changes the distribution model. Anyone can play
   without installing anything. This is the viral hook.

4. **Phase 4 (Daily + Leaderboards)**: Adds retention. Players come back every
   day for the daily challenge and to compete on the leaderboard.

5. **Phase 5 (Polish)**: Makes it deployable and presentable.

Each phase produces a working, shippable game. You can stop after any phase
and have something complete.
