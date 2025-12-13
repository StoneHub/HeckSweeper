# HeckSweeper Redesign Plan

## Executive Summary

The original PowerShell-based implementation has proven unviable for production due to fundamental platform limitations (PS2EXE encoding issues, console API fragility, cross-platform inconsistencies). This plan outlines a redesign using modern technologies that preserve the "terminal aesthetic" while enabling robust, cross-platform distribution.

---

## Lessons Learned from Failed Approach

### Critical Failures

1. **PowerShell/PS2EXE Packaging is Fundamentally Broken**
   - UTF-8 encoding issues cause parse failures
   - PS2EXE embeds PowerShell 5.1, creating version mismatches
   - Console APIs fail silently or throw in packaged builds
   - No viable workaround without abandoning Unicode entirely

2. **Platform Fragility**
   - WSL + Windows Terminal + PowerShell creates 3-layer complexity
   - Different behavior across environments (conhost vs. Windows Terminal)
   - Path resolution requires brittle hacks
   - Windows Terminal launch is unreliable

3. **Limited Extensibility**
   - Hard-coded 24×16 board prevents scaling
   - Pure console rendering limits UI possibilities
   - No asset system or configuration infrastructure
   - Testing requires headless mode hacks

4. **Poor Developer Experience**
   - Slow iteration cycles (packaging takes ~10-15 seconds)
   - Sparse error messages from PowerShell modules
   - No debugging tools or hot reload
   - Documentation/code drift

### What Actually Worked

✅ **Core Game Logic** - The board generation, flood-fill reveal, and win/loss detection algorithms are solid
✅ **Modular Architecture** - Clean separation (Board, Gen, Render, Input, Game) is maintainable
✅ **Deterministic Seeding** - Reproducible gameplay is valuable for testing
✅ **Double-Buffered Rendering** - StringBuilder approach minimizes flicker
✅ **Fallback Strategies** - ASCII mode as backup shows good defensive design

---

## Design Goals for Redesign

### Must Have
1. **Terminal Aesthetic** - Monospace fonts, ASCII/Unicode art, retro feel
2. **Cross-Platform** - Works on Windows, macOS, Linux without platform-specific hacks
3. **Single-File Distribution** - One executable, no runtime dependencies
4. **Reliable Rendering** - No encoding issues, consistent behavior
5. **Extensible Architecture** - Easy to add new mechanics, board sizes, game modes

### Nice to Have
6. **Modern UI Options** - Terminal emulation in a GUI window (not limited to actual terminal)
7. **Rich Graphics** - Sprite support, animations, particle effects (while maintaining retro style)
8. **Sound/Music** - Optional audio (chiptune aesthetic)
9. **Multiple Themes** - Color schemes, glyph sets
10. **Persistent State** - Saves, leaderboards, achievements

---

## Technology Stack Options

### Option A: Rust + crossterm (Pure Terminal)
**Stack**: Rust + crossterm/ratatui + serde

**Pros:**
- True cross-platform terminal rendering (no encoding issues)
- Single binary compilation with zero runtime dependencies
- Excellent performance
- Rich ecosystem (crossterm handles all ANSI complexity)
- ratatui provides TUI widgets (borders, tables, layouts)
- Strong type safety prevents bugs

**Cons:**
- Steeper learning curve if unfamiliar with Rust
- Terminal-only (can't add GUI features later without major rewrite)
- Text-based rendering limits visual flair

**Best For:** Pure terminal game that ships as a robust CLI tool

---

### Option B: Python + Rich/Textual (Hybrid Terminal)
**Stack**: Python 3.10+ + Rich/Textual + PyInstaller

**Pros:**
- Rapid development (familiar syntax, great for prototyping)
- Rich library provides beautiful terminal formatting out-of-box
- Textual enables reactive TUI with widgets (buttons, progress bars, etc.)
- PyInstaller creates single-file executables
- Easy to add web interface later (Flask/FastAPI)

**Cons:**
- Larger binary size (~15-30MB with PyInstaller)
- Slower runtime performance than compiled languages
- PyInstaller can have platform-specific quirks
- Python dependency management can be complex

**Best For:** Rapid prototyping with option to pivot to web/GUI later

---

### Option C: Go + bubbletea (Modern TUI)
**Stack**: Go + bubbletea (Elm-inspired TUI framework)

**Pros:**
- Cross-compilation to all platforms trivially
- Single binary with no runtime (small ~5-10MB)
- bubbletea provides elegant reactive model (Elm architecture)
- Fast compilation and runtime
- Great standard library (testing, HTTP, file I/O)
- Simple concurrency model for animations/timers

**Cons:**
- Less mature TUI ecosystem than Rust
- Go's syntax can feel verbose for game logic
- Terminal-only (no easy GUI upgrade path)

**Best For:** Clean, maintainable TUI with excellent distribution story

---

### Option D: JavaScript/TypeScript + Electron (Terminal Emulator)
**Stack**: TypeScript + xterm.js + Electron

**Pros:**
- Not limited to terminal - renders terminal emulator in GUI window
- Full control over rendering (can add sprites, animations, shaders)
- Huge ecosystem (npm packages for everything)
- Hot reload for fast iteration
- Can add web multiplayer, Discord integration, etc.
- Beautiful typography (custom fonts, ligatures)

**Cons:**
- Large binary size (~150-200MB with Electron)
- Higher memory usage (~100MB+ at runtime)
- Slower startup time
- "Fake terminal" feel may not be authentic

**Best For:** Terminal aesthetic with modern features (cloud saves, multiplayer, rich UI)

---

### Option E: C# + Spectre.Console (Cross-Platform .NET)
**Stack**: C# 12 + Spectre.Console + .NET 8 AOT

**Pros:**
- Spectre.Console provides gorgeous terminal UI (tables, trees, spinners)
- .NET 8 AOT compilation creates small, fast native binaries
- Excellent debugging in Visual Studio/Rider
- Strong typing and LINQ for clean game logic
- Could add Blazor WebAssembly version later

**Cons:**
- Requires .NET SDK for development (though binaries are standalone)
- AOT can have compatibility issues with reflection-heavy code
- Terminal-focused (but could use AvaloniaUI for GUI later)

**Best For:** Developers familiar with C# wanting professional tooling

---

## Recommended Approach: **Option D (TypeScript + Electron + xterm.js)**

### Rationale

Given the requirement "runs like a terminal but not necessarily limited to using that," **Electron + xterm.js** provides the perfect balance:

1. **Authentic Terminal Feel** - xterm.js is a real terminal emulator (used by VS Code, Hyper)
2. **No Platform Limitations** - No encoding issues, no console API fragility
3. **Extensibility** - Can add sprites, animations, sound without rewriting
4. **Modern DX** - TypeScript, hot reload, great debugging
5. **Distribution** - electron-builder creates installers for Windows/macOS/Linux
6. **Future-Proof** - Easy to add web version, multiplayer, mods, etc.

### Trade-Off Justification

Yes, the binary will be ~150MB vs. ~5MB for Rust/Go. But:
- Storage is cheap (game size is negligible on modern systems)
- Development speed matters more for a solo/small team project
- The ability to iterate quickly and add features beats binary size
- Users expect ~100-200MB for indie games

---

## Proposed Architecture (TypeScript + Electron)

### Tech Stack Details

```
Frontend: TypeScript 5 + xterm.js + xterm-addon-fit
Backend: Electron 28+ (main/renderer process split)
Build: Vite (fast bundling) + electron-builder (packaging)
Testing: Vitest (unit) + Playwright (E2E)
State: Zustand or Redux Toolkit (game state management)
Audio: Howler.js (sound effects)
Config: JSON5 (human-friendly config files)
```

### Module Structure

```
src/
├── main/                    # Electron main process
│   ├── index.ts            # App initialization, window management
│   ├── config.ts           # Load user config, defaults
│   └── ipc.ts              # IPC handlers (save/load game)
│
├── renderer/                # Renderer process (game UI)
│   ├── terminal/
│   │   ├── Terminal.ts     # xterm.js wrapper, ANSI rendering
│   │   ├── InputHandler.ts # Keyboard event processing
│   │   └── Renderer.ts     # Frame rendering to xterm buffer
│   │
│   ├── game/
│   │   ├── Board.ts        # Board state (cells, flags, reveal)
│   │   ├── Generator.ts    # Map generation (seeded RNG)
│   │   ├── GameLoop.ts     # Main loop (tick, input, render)
│   │   └── Actions.ts      # Game actions (move, reveal, flag)
│   │
│   ├── state/
│   │   ├── store.ts        # Zustand store (game state)
│   │   └── types.ts        # TypeScript types/interfaces
│   │
│   ├── audio/
│   │   └── SoundManager.ts # Sound effects (reveal, flag, win/loss)
│   │
│   └── ui/
│       ├── HUD.ts          # Status bar (flags, moves, timer)
│       ├── Menu.ts         # Main menu, settings
│       └── Themes.ts       # Color schemes, glyph sets
│
├── shared/
│   ├── constants.ts        # Game constants (board sizes, difficulty)
│   └── utils.ts            # Shared utilities (RNG, indexing)
│
└── assets/
    ├── sounds/             # .wav/.ogg files
    ├── fonts/              # Monospace fonts (Fira Code, Hack)
    └── themes/             # JSON theme definitions
```

### Key Design Patterns

1. **Entity-Component-System (Light)**
   - Cells are simple data structures
   - Systems process cells (RevealSystem, FlagSystem, RenderSystem)
   - Avoids deep inheritance trees

2. **Command Pattern for Actions**
   - Each user input becomes a Command object
   - Commands are undoable (for replay/undo features)
   - Easy to serialize for save/load

3. **Double Buffering**
   - Maintain current and next frame buffers
   - Write to xterm only on delta (minimize redraws)
   - Use xterm.write() batch API for performance

4. **Immutable State Updates**
   - Use immer (via Zustand) for immutable state
   - Makes debugging easier (state history)
   - Enables time-travel debugging

---

## Implementation Phases

### Phase 1: Foundation (Week 1-2)
**Goal:** Minimal playable game in Electron terminal

**Tasks:**
- [ ] Set up Electron + Vite + TypeScript project
- [ ] Integrate xterm.js with fit addon
- [ ] Implement basic ANSI rendering (colors, box-drawing)
- [ ] Port Board.ts (cell state, reveal, flag logic)
- [ ] Port Generator.ts (seeded RNG, monster placement)
- [ ] Implement keyboard input handling (arrow keys, WASD)
- [ ] Create simple game loop (render → input → update)
- [ ] Win/loss detection
- [ ] Basic smoke test (automated input sequence)

**Success Criteria:**
- 24×16 board renders with Unicode glyphs
- Cursor movement works
- Reveal/flag mechanics function correctly
- Win when all safe cells revealed
- Lose when monster revealed

---

### Phase 2: Polish & Features (Week 3-4)
**Goal:** Production-ready game with quality-of-life improvements

**Tasks:**
- [ ] Add HUD (remaining flags, moves counter, timer)
- [ ] Implement main menu (New Game, Settings, Quit)
- [ ] Add settings UI (difficulty, board size, theme)
- [ ] Create 3 difficulty presets (Easy 16×12, Normal 24×16, Hard 32×24)
- [ ] Sound effects (reveal, flag, win, lose)
- [ ] Smooth cursor animations
- [ ] Flood-fill reveal animation (wave effect)
- [ ] Victory/defeat screens with stats
- [ ] Save/load game state (JSON files)
- [ ] High score leaderboard (local file)

**Success Criteria:**
- Game feels polished and responsive
- Settings persist between sessions
- Audio enhances experience (with mute option)
- UI/UX is intuitive

---

### Phase 3: Distribution (Week 5)
**Goal:** Packaged installers for all platforms

**Tasks:**
- [ ] Configure electron-builder for Windows (NSIS installer)
- [ ] Configure electron-builder for macOS (DMG, code signing)
- [ ] Configure electron-builder for Linux (AppImage, deb, rpm)
- [ ] Add auto-updater support (electron-updater)
- [ ] Create landing page (GitHub Pages or simple static site)
- [ ] Write user manual / gameplay guide
- [ ] Set up GitHub Releases workflow
- [ ] Test installers on clean VMs

**Success Criteria:**
- One-click install on Windows/macOS/Linux
- Installers are <200MB
- Auto-update works
- No runtime dependencies required

---

### Phase 4: Extended Features (Post-MVP)
**Goal:** Make the game stand out

**Potential Features:**
- [ ] Daily challenge mode (global seed shared daily)
- [ ] Procedural dungeon themes (catacombs, mines, temples)
- [ ] Power-ups (reveal radius, extra flags, undo move)
- [ ] Campaign mode (progressive difficulty)
- [ ] Speedrun timer with splits
- [ ] Replay system (record/playback games)
- [ ] Steam/itch.io integration
- [ ] Custom map editor
- [ ] Mod support (load custom themes/rulesets)
- [ ] Web version (deploy to web, same TypeScript codebase)

---

## Migration Strategy

### Code Reuse from PowerShell Version

1. **Board Logic (Board.psm1)** → `Board.ts`
   - Port cell state structure (Revealed, HasMonster, IsFlagged)
   - Port neighbor indexing logic
   - Port flood-fill reveal (BFS queue)
   - **Estimated Effort:** 4-6 hours

2. **Generator (Gen.psm1)** → `Generator.ts`
   - Port seeded RNG initialization
   - Port monster placement algorithm
   - **Estimated Effort:** 2-3 hours

3. **Rendering (Render.psm1)** → `Renderer.ts`
   - Replace StringBuilder with xterm.js buffer writes
   - Port glyph mapping (Unicode/ASCII modes)
   - **Estimated Effort:** 4-6 hours

4. **Input (Input.psm1)** → `InputHandler.ts`
   - Replace Console.ReadKey with xterm.onKey event
   - Port key mapping (arrows, WASD, HJKL)
   - **Estimated Effort:** 2-3 hours

5. **Game Loop (Game.psm1)** → `GameLoop.ts`
   - Replace PowerShell loop with requestAnimationFrame
   - Port game state transitions
   - **Estimated Effort:** 3-4 hours

**Total Migration Effort:** ~15-22 hours (2-3 days)

---

## Testing Strategy

### Unit Tests (Vitest)
- Board state mutations (reveal, flag, unflag)
- Generator determinism (same seed = same board)
- Win/loss detection
- Neighbor calculation edge cases
- Input command parsing

**Target Coverage:** 80%+ for game logic

### Integration Tests
- Full game flow (new game → moves → win/loss)
- Save/load roundtrip
- Settings persistence

### E2E Tests (Playwright)
- Automated playthrough (scripted inputs)
- UI interaction (menu navigation, settings changes)
- Cross-platform smoke tests (Windows/macOS/Linux runners)

### Manual Testing Checklist
- [ ] All key bindings work
- [ ] Audio plays correctly (and mute works)
- [ ] Resizing window adapts terminal
- [ ] High DPI displays render sharply
- [ ] Performance is smooth (60 FPS)
- [ ] No memory leaks over long sessions

---

## Success Metrics

### Technical
- ✅ Single-file executable for each platform
- ✅ <200MB installer size
- ✅ <100MB memory usage at runtime
- ✅ Startup time <2 seconds
- ✅ 60 FPS rendering (smooth animations)
- ✅ Zero encoding/Unicode issues
- ✅ 80%+ test coverage

### User Experience
- ✅ Intuitive controls (no tutorial needed)
- ✅ Polished visuals (retro but professional)
- ✅ Satisfying audio feedback
- ✅ Clear win/loss conditions
- ✅ Settings are discoverable and work

### Distribution
- ✅ Installable on Windows 10/11, macOS 12+, Ubuntu 20.04+
- ✅ No runtime dependencies
- ✅ Auto-update works
- ✅ GitHub releases with binaries
- ✅ 100+ GitHub stars (aspirational)

---

## Risk Mitigation

### Risk: Electron bundle size too large
**Mitigation:** Use electron-builder compression, exclude unnecessary Node modules, consider Tauri (Rust alternative to Electron with ~10MB binaries) if size becomes critical

### Risk: xterm.js performance issues on large boards
**Mitigation:** Profile early, optimize rendering (delta updates only), cap board size at 50×50, consider canvas-based rendering for larger boards

### Risk: Cross-platform audio issues
**Mitigation:** Test early on all platforms, use common audio formats (.ogg), provide mute option, make audio optional

### Risk: Development takes longer than estimated
**Mitigation:** Release Phase 1 as "Early Access" to get feedback, prioritize core gameplay over features, cut Phase 4 features if needed

---

## Alternative: Lightweight Option (Go + bubbletea)

If Electron's size/complexity is a concern, **Go + bubbletea** is the best alternative:

### Quick Comparison

| Aspect | Electron + xterm.js | Go + bubbletea |
|--------|---------------------|----------------|
| Binary Size | ~150-200MB | ~5-10MB |
| Memory Usage | ~100MB | ~10-20MB |
| Startup Time | 1-2s | <100ms |
| Development Speed | Fast (TypeScript, hot reload) | Moderate (Go compilation) |
| UI Capabilities | Full GUI potential | Terminal only |
| Learning Curve | Low (if familiar with JS/TS) | Moderate (Go + Elm patterns) |
| Distribution | electron-builder | go build (trivial) |
| Future Extensibility | High (web, GUI, plugins) | Low (terminal-locked) |

### When to Choose Go
- Binary size is critical (embedded systems, bandwidth-constrained)
- Pure terminal aesthetic is non-negotiable
- Team is already Go-proficient
- Want minimal runtime footprint

### When to Choose Electron
- Want room to grow (GUI, web, advanced features)
- Team is JavaScript/TypeScript-focused
- Development speed matters more than binary size
- Want modern debugging/tooling

---

## Next Steps

1. **Decision Point:** Choose tech stack (recommend Electron, but Go is valid alternative)
2. **Repository Setup:**
   - Create new repo `hecksweeper-v2` or `hecksweeper-electron`
   - Initialize with chosen template (electron-vite-template or bubbletea starter)
   - Set up CI/CD (GitHub Actions for build/test/release)
3. **Phase 1 Kickoff:**
   - Implement basic Electron + xterm.js integration
   - Port core game logic (Board, Generator)
   - Get minimal playable demo running
4. **Feedback Loop:**
   - Share early build for playtesting
   - Iterate on feel (cursor speed, reveal animation, audio)
   - Refine based on feedback

---

## Appendix: Quick Start Templates

### Electron + TypeScript + Vite
```bash
npm create @quick-start/electron
# Choose: TypeScript + Vite + React (or Vanilla)
cd hecksweeper-v2
npm install xterm xterm-addon-fit
npm run dev
```

### Go + bubbletea
```bash
go mod init github.com/yourusername/hecksweeper-v2
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss  # for styling
# Create main.go with bubbletea boilerplate
go run main.go
```

### Rust + ratatui
```bash
cargo new hecksweeper-v2
cd hecksweeper-v2
cargo add ratatui crossterm
# Edit main.rs with ratatui setup
cargo run
```

---

**End of Redesign Plan**

*Last Updated: 2025-12-13*
