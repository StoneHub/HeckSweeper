# DungeonSweeper_PS – WSL Agent Dev Plan

Goal: ship an MVP playable in terminal from WSL and package a Windows `.exe` via a Makefile that the agent generates.

## 0) Constraints
- Code: pure PowerShell 7 (cross‑platform). No native DLLs. No third‑party libs.
- Render: ANSI VT + Unicode. ASCII fallback.
- WSL daily dev; final Windows packaging via `PS2EXE` using host `powershell.exe` from WSL.

## 1) Repo layout
```
/ (root)
  DungeonSweeper.ps1            # entrypoint
  /src                          # modules split after MVP
    Board.psm1
    Gen.psm1
    Render.psm1
    Input.psm1
    Game.psm1
  /scripts
    bootstrap.ps1               # one‑time env checks
    package.ps1                 # PS2EXE wrapper (Windows host)
  /assets
    glyph_compat.md
  Makefile
  README.md
```

## 2) MVP definition
- Fixed map size that fits 80×25.
- Cursor move, reveal, flag, win/lose.
- Monsters as mines with radius=1 for MVP; counts shown.
- Fog‑of‑war, flood‑reveal of zero‑threat.
- Double‑buffer redraw from (0,0). No `Clear()` calls. Hide/show cursor around loop.
- Seeded RNG for reproducibility.

## 3) Milestones
**M1 – Bootstrap (WSL)**
- Ensure `pwsh` in WSL and UTF‑8 defaults.
- Verify Windows host interop path for `powershell.exe`.

**M2 – Play loop**
- Input polling with `Console.KeyAvailable` + `ReadKey()`.
- Back‑buffer builder, HUD line, status line.

**M3 – Board + logic**
- Grid model, monster placement, threat counts.
- Flood‑reveal BFS, flagging and win/lose checks.

**M4 – Packaging**
- PS2EXE packaging via Windows host from WSL.
- Emit `/build/DungeonSweeper.exe`.

**M5 – Compat pass**
- ASCII mode, 16‑color fallback.
- Legacy conhost sanity test.

## 4) Acceptance tests
- Start → first frame under 100ms.
- Reveal zero → expands deterministically.
- Win detection when non‑monster cells revealed.
- Packaging produces an `.exe` < 10 MB that launches to the game screen.

## 5) Makefile (WSL‑first, Windows packaging from WSL)
```make
# Makefile
SHELL := /bin/bash
WSL_PWSH := pwsh
WIN_PWSH := /mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe
BUILD_DIR := build
EXE := $(BUILD_DIR)/DungeonSweeper.exe

.PHONY: all run test package clean setup-wsl setup-win

all: run

setup-wsl:
	@$(WSL_PWSH) -NoLogo -NoProfile -Command "\
	  $ErrorActionPreference='Stop'; \
	  if(-not (Get-Command pwsh -ErrorAction SilentlyContinue)) { throw 'PowerShell 7 not found in WSL'; } \
	  [Console]::OutputEncoding=[Text.UTF8Encoding]::new($false) \
	"
	@echo "WSL setup ok"

run: setup-wsl
	@$(WSL_PWSH) -NoLogo -ExecutionPolicy Bypass -File ./DungeonSweeper.ps1

# headless logic smoke test (expand later)
test:
	@$(WSL_PWSH) -NoLogo -ExecutionPolicy Bypass -Command "\
	  Write-Host 'Seed=123'; \
	  . ./DungeonSweeper.ps1 -Seed 123 -Headless \
	"

package: $(EXE)

$(EXE): DungeonSweeper.ps1 scripts/package.ps1 | $(BUILD_DIR)
	@"$(WIN_PWSH)" -NoLogo -NoProfile -ExecutionPolicy Bypass -File scripts/package.ps1 -Input ./DungeonSweeper.ps1 -Output $(EXE)
	@echo "Packed -> $(EXE)"

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

clean:
	rm -rf $(BUILD_DIR)
```

## 6) scripts/package.ps1 (runs on Windows host from WSL)
```powershell
param(
  [Parameter(Mandatory)] [string]$Input,
  [Parameter(Mandatory)] [string]$Output
)
$ErrorActionPreference='Stop'
function Ensure-PS2EXE {
  if (-not (Get-Command Invoke-ps2exe -ErrorAction SilentlyContinue)) {
    Install-Module PS2EXE -Scope CurrentUser -Force -AllowClobber
  }
}
Ensure-PS2EXE
$null = New-Item -ItemType Directory -Force -Path (Split-Path $Output)
Invoke-ps2exe -inputFile $Input -outputFile $Output -noConsole -title 'DungeonSweeper' -verbose
Write-Host "OK: $Output"
```

## 7) scripts/bootstrap.ps1 (WSL session)
```powershell
$ErrorActionPreference='Stop'
$PSStyle.OutputRendering='Ansi'
[Console]::OutputEncoding=[Text.UTF8Encoding]::new($false) | Out-Null
$wt = $env:WT_SESSION
if (-not $wt) { Write-Host 'Note: not in Windows Terminal; enabling 16-color fallback at runtime' }
```

## 8) DungeonSweeper.ps1 flags to implement now
- `-Seed [int]` reproducible runs.
- `-Compat` ASCII + 16‑color mode.
- `-Headless` logic self‑test that exits.

## 9) Dev tasks for the agent
1. Add flags and bootstrap. Ensure UTF‑8 and cursor hide/show.
2. Implement board model and threat counts.
3. Implement flood‑reveal and win/lose.
4. Write renderer with double buffer (string builder per row → single Write()).
5. Add key handling loop (arrows/WASD, Space=Reveal, F=Flag, Q=Quit).
6. Add ASCII fallback rendering.
7. Create `scripts/package.ps1` and Makefile. Verify `make package` emits `.exe` from WSL.
8. Write README with usage: `make run`, `make package`.

## 10) Risks + mitigations
- **PS2EXE availability:** install on first package run.
- **Path issues:** Use relative paths in repo. Let Makefile create `build/`.
- **Console flicker:** never clear; full frame write from (0,0).
- **Glyph gaps:** default to ASCII in `-Compat`.

## 11) Done criteria for MVP PR
- `make run` plays a round.
- `make package` creates a Windows `.exe` that starts the game full‑screen‑like in Windows Terminal.
- README shows controls and flags.
```

