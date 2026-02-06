# HeckSweeper

A roguelite minesweeper for the terminal. Clear floors, collect power-ups, climb the leaderboard. Play locally or over SSH.

## Play

```bash
# Local
go run ./cmd/hecksweeper

# Over SSH (connect to a hosted server)
ssh -p 2222 localhost
```

## How It Works

Each run is a series of minesweeper floors with escalating difficulty. Clear a floor, pick a power-up, descend deeper. Die and it's over.

- **Floors 1-6**: Board grows from 12x8 to 24x16, monster density from 12% to 22%
- **Floor 7+**: Max size, density keeps climbing (caps at 35%)
- **Power-ups**: Scout, Shield, X-Ray, Lucky Charm, and more -- 8 total across 4 rarity tiers
- **Daily Challenge**: Same seed for everyone, compete on the leaderboard
- **Scoring**: Cells revealed + speed bonus + efficiency bonus, scaled by floor depth

## Controls

| Key | Action |
|-----|--------|
| Arrow / WASD / HJKL | Move cursor |
| Space / Enter | Reveal cell |
| F | Toggle flag |
| Esc | Back / Quit |

## Build

Requires Go 1.24+.

```bash
make -f Makefile.new build       # Build local client
make -f Makefile.new server      # Build SSH server
make -f Makefile.new test        # Run tests (40+ tests)
make -f Makefile.new release     # Optimized binaries
make -f Makefile.new cross       # Cross-compile all platforms
```

## Host an SSH Server

```bash
# Build and run
make -f Makefile.new run-server

# Or directly
go run ./cmd/server --port 2222

# Docker
docker build -t hecksweeper .
docker run -p 2222:2222 -v hecksweeper-data:/data hecksweeper
```

The server auto-generates an SSH host key on first run. Scores persist to a JSON file.

Players are identified by SSH key fingerprint for leaderboards. No account needed -- just `ssh -p 2222 host` and play.

### Server Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--host` | 0.0.0.0 | Listen address |
| `--port` | 2222 | SSH port |
| `--key` | .ssh/hecksweeper_ed25519 | Host key path |
| `--data` | hecksweeper_data.json | Score data file |

## Architecture

```
cmd/hecksweeper/     Local client entry point
cmd/server/          SSH server entry point
internal/game/       Game logic (board, run, power-ups)
internal/ui/         Bubbletea views and input handling
internal/server/     Wish SSH server + handler
internal/daily/      Daily challenge seed generation
internal/storage/    Persistent score storage
internal/constants/  Shared constants and glyphs
```

Game logic is fully separated from UI. The SSH server wraps the same bubbletea model used locally -- each connection gets an isolated game instance.

## Tech Stack

- **Go** + **bubbletea** (TUI framework)
- **lipgloss** (terminal styling)
- **Wish** (SSH server, serves bubbletea apps over SSH)
- Zero CGO, fully static binaries
