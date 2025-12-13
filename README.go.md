# HeckSweeper (Go Edition)

**A lightweight terminal roguelite inspired by Minesweeper**

Built with Go + [bubbletea](https://github.com/charmbracelet/bubbletea) for a clean, fast, cross-platform terminal experience.

![HeckSweeper Demo](demo.gif)

## Features

✨ **Pure Terminal** - Runs in any modern terminal with ANSI support
🎮 **Roguelite Gameplay** - Dungeon exploration meets Minesweeper mechanics
🎨 **Beautiful Styling** - Powered by lipgloss with color-coded threat levels
🔢 **Deterministic Seeds** - Reproducible games for speedrunning & challenges
⚡ **Tiny Binary** - ~4-5MB executable, no runtime dependencies
🌍 **Cross-Platform** - Linux, macOS, Windows (anywhere Go runs)
⌨️ **Multiple Control Schemes** - Arrow keys, WASD, or Vim (hjkl)

## Quick Start

### Prerequisites
- Go 1.21+ (for building)
- OR just download a pre-built binary from [Releases](#)

### Build & Run
```bash
# Clone the repo
git clone https://github.com/StoneHub/HeckSweeper.git
cd HeckSweeper

# Build and run
make run

# Or build manually
go build -o hecksweeper ./cmd/hecksweeper
./hecksweeper
```

## Controls

| Action | Keys |
|--------|------|
| **Move** | Arrow Keys / WASD / HJKL (Vim) |
| **Reveal** | Space / Enter |
| **Flag** | F |
| **Quit** | Q / Esc / Ctrl+C |
| **Restart** | R (after game over) |
| **New Game** | N (after game over) |

## Gameplay

- 🟢 **Goal**: Reveal all safe cells without hitting a monster
- ☠️ **Monsters**: Hidden throughout the dungeon (18% of cells by default)
- 🔢 **Numbers**: Show how many monsters are in adjacent cells (8-directional)
- ⚑ **Flags**: Mark cells you suspect contain monsters
- 💥 **Loss**: Revealing a monster ends the game
- 🎉 **Victory**: Reveal all safe cells to win

### Threat Level Colors
```
0 monsters → Green
1-2 monsters → Cyan/Light Green
3-4 monsters → Yellow/Orange
5-6 monsters → Dark Orange/Red-Orange
7-8 monsters → Red/Dark Red
```

## Command-Line Options

```bash
# Custom board size
./hecksweeper -width 32 -height 24

# Fixed seed for reproducible games
./hecksweeper -seed 42

# ASCII mode (for terminals without Unicode support)
./hecksweeper -unicode=false

# Combine options
./hecksweeper -width 16 -height 12 -seed 123 -unicode=true
```

## Makefile Targets

```bash
make build        # Build the binary
make release      # Build optimized binary (smaller, no debug symbols)
make run          # Build and run with default settings
make run-ascii    # Run in ASCII compatibility mode
make run-seed     # Run with fixed seed (42)
make run-large    # Run with 32×24 board
make run-small    # Run with 16×12 board
make test         # Run tests
make install      # Install to /usr/local/bin (requires sudo)
make clean        # Remove build artifacts
```

## Project Structure

```
.
├── cmd/
│   └── hecksweeper/
│       └── main.go          # Entry point
├── internal/
│   ├── game/
│   │   ├── board.go         # Board state and logic
│   │   ├── generator.go     # Procedural generation (seeded RNG)
│   │   └── types.go         # Game types
│   ├── ui/
│   │   ├── model.go         # Bubbletea model
│   │   ├── update.go        # Input handling
│   │   ├── view.go          # Rendering logic
│   │   └── styles.go        # Lipgloss styles
│   └── constants/
│       └── constants.go     # Game constants & glyphs
├── go.mod
├── go.sum
├── Makefile.go
└── README.go.md
```

## Architecture

### Clean Separation of Concerns

**Game Logic** (`internal/game/`)
- Pure game state (no UI dependencies)
- Deterministic board generation
- Flood-fill reveal algorithm (BFS)
- Win/loss detection

**UI Layer** (`internal/ui/`)
- Bubbletea Elm-inspired architecture (Model-Update-View)
- Keyboard input handling
- Terminal rendering with lipgloss
- No game logic (just presentation)

**Constants** (`internal/constants/`)
- Game parameters (board size, monster density)
- Glyph sets (Unicode and ASCII)
- Game state constants

### Key Algorithms

1. **Seeded Random Generation** - Uses `math/rand` with fixed seeds for reproducible dungeons
2. **Flood-Fill Reveal** - BFS queue-based expansion from zero-threat cells
3. **Neighbor Calculation** - 8-directional adjacency for threat counting
4. **Double-Buffering** - Bubbletea's efficient terminal rendering

## Development

### Running Tests
```bash
go test ./...

# With coverage
go test -cover ./...

# Generate coverage report
make test-coverage
```

### Code Formatting
```bash
go fmt ./...
# or
make fmt
```

### Updating Dependencies
```bash
go get -u ./...
go mod tidy
# or
make deps
```

## Design Philosophy

### Why Go + bubbletea?

After failing with PowerShell/PS2EXE (see [REDESIGN_PLAN.md](REDESIGN_PLAN.md)), we chose Go for:

1. **Zero Runtime Dependencies** - Single binary, no PS2EXE Unicode issues
2. **Cross-Platform** - Trivial compilation to Windows/Linux/macOS
3. **Small Binaries** - ~4-5MB (vs 150-200MB for Electron)
4. **Fast Startup** - <100ms to launch
5. **Clean Architecture** - Go's simplicity encourages good design
6. **Terminal Native** - bubbletea handles ANSI complexity elegantly

### Creative Touches

- **Dynamic Threat Colors** - Visual gradient from green (safe) to red (danger)
- **Dual Glyph Sets** - Unicode for modern terminals, ASCII fallback
- **Elm Architecture** - Predictable state management via bubbletea
- **Keyboard Shortcuts** - Multiple control schemes (arrows, WASD, Vim)
- **Seeded Generation** - Speedrunning and reproducible challenges
- **Clean Visuals** - Box-drawing borders, styled stats bar

## Contributing

Contributions welcome! Areas for enhancement:

- [ ] Add difficulty presets (easy/normal/hard)
- [ ] Implement daily challenge mode
- [ ] Add speedrun timer
- [ ] Create replay system
- [ ] Add sound effects (via beep/terminal bell)
- [ ] Procedural dungeon themes
- [ ] Power-ups (extra reveals, safe zones)
- [ ] Leaderboard (local file)

## License

MIT License - See LICENSE file

## Credits

- Built with [bubbletea](https://github.com/charmbracelet/bubbletea) by Charm
- Styled with [lipgloss](https://github.com/charmbracelet/lipgloss)
- Inspired by classic Minesweeper and roguelike games

## Related

- [REDESIGN_PLAN.md](REDESIGN_PLAN.md) - Full redesign rationale
- [Original PowerShell version](src/) - The failed first attempt (lessons learned!)

---

**Made with ❤️ and Go**
*A HeckSweeper production*
