<<<<<<< HEAD
# Code-Playground
=======
# DungeonSweeper

DungeonSweeper is a PowerShell 7 terminal roguelite inspired by Minesweeper. Explore a fixed 24×16 dungeon from WSL with crisp ANSI rendering or an ASCII fallback.

## Requirements
- PowerShell 7 (`pwsh`) available in WSL for development and gameplay.
- Windows PowerShell (`powershell.exe`) on the host for packaging to `.exe`.
- Optional: PS2EXE module (installed automatically during packaging).

## Quick Start
```bash
make run     # launch the interactive game loop
make test    # run the deterministic headless smoke test (Seed=123)
make package # build build/DungeonSweeper.exe via PS2EXE
```
The bootstrap script (`scripts/bootstrap.ps1`) can be run manually if you are not using the Makefile yet.

## Controls
Use the arrow keys or WASD (also HJKL) to move the cursor. `Space`/`Enter` reveal the current room, `F` toggles a flag, and `Q` (or `Esc`/`Ctrl+C`) quits. Hidden rooms show `·` (or `.` in ASCII mode), flagged rooms use `⚑`/`F`, and revealed monsters display `☠`/`X`.

## Headless & Compatibility Modes
Enable compatibility rendering with `./DungeonSweeper.ps1 -Compat`. The command `./DungeonSweeper.ps1 -Headless -Seed <value>` emits a structured summary for testing and reproducible scenarios. ASCII mode uses plain borders and cursor markers so legacy terminals remain usable.

## Packaging Notes
`make package` calls `scripts/package.ps1` from the Windows host to produce `build/DungeonSweeper.exe`. The wrapper ensures PS2EXE is available and writes the executable without polluting the repository. Remove the generated binary before committing.
>>>>>>> 8b89520 (Add .gitignore, enhance input handling, and implement packaging script)
