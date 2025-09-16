# Repository Guidelines

## Project Structure & Module Organization
PowerShell gameplay lives in `DungeonSweeper.ps1` with supporting modules in `src/Board.psm1`, `src/Gen.psm1`, `src/Render.psm1`, `src/Input.psm1`, and `src/Game.psm1`. Shared scripts reside under `scripts/`, including `bootstrap.ps1` for environment prep and `package.ps1` for Windows builds. Binary artifacts and packaging outputs land in `build/` (gitignored). Store static glyph references in `assets/`, and keep future design docs such as `dungeon_sweeper_ps_wsl_agent_dev_plan.md` at the repository root.

## Build, Test, and Development Commands
Run `make run` for the default WSL gameplay loop; it calls `pwsh` with ANSI rendering enabled. Use `make test` to execute the headless logic smoke test (`DungeonSweeper.ps1 -Seed 123 -Headless`). Package the Windows executable from WSL with `make package`, which shells into the host `powershell.exe` and invokes `scripts/package.ps1`. Execute `scripts/bootstrap.ps1` manually the first time if the Makefile is not available yet.

## Coding Style & Naming Conventions
Write PowerShell 7 code with four-space indentation and braces on the same line as keywords (`if (...) {`). Prefer verb-noun function names (`Get-BoardState`, `Invoke-RenderFrame`) and PascalCase public functions inside modules. Keep private helpers scoped with the `script:` prefix or local functions. Align ANSI sequences and ASCII fallbacks in separate rendering paths with concise comments explaining non-obvious buffer work.

## Testing Guidelines
Exercise deterministic behaviour with the `-Seed` parameter and capture console output when adding regressions tests. Extend the `make test` target with additional scripted scenarios rather than ad-hoc scripts. Place test assets under `tests/` if full fixtures are needed, mirroring the module layout. Document any new headless switches in the README and ensure they degrade gracefully in plain ASCII mode.

## Commit & Pull Request Guidelines
Write imperative, 72-character subject lines (`Add flood-fill reveal BFS`). Reference issues with `Fixes #123` when applicable and document gameplay changes or compatibility decisions in the body. PRs should link to relevant sections of the dev plan, include reproduction steps (`make run`, scenarios exercised), and attach terminal recordings or screenshots when UI changes are involved. Confirm that CI (or local `make test`) passes before requesting review.

## Environment & Packaging Notes
Verify `pwsh` availability in WSL and ensure UTF-8 output (`[Console]::OutputEncoding`). Packaging requires PS2EXE on the Windows host; the wrapper installs it on demand but note the download prompt for first-time contributors. Avoid committing generated `.exe` files—use `build/` for scratch artifacts.
