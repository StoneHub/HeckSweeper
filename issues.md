# Packaging & Window Mode Issues

## PowerShell Packaging Constraints
- PS2EXE embeds PowerShell 5.1, so console APIs we call in WSL (UTF-8 encoding, cursor visibility, `TreatControlCAsInput`) either throw or silently do nothing. We wrapped those calls in `try { } catch { }` and defaulted to ASCII mode when UTF-8 wasn’t available, but it still limits fidelity.
- The executable runs the script from `build\` while our modules live under `src\`. When `$PSCommandPath` comes back null in the packaged app, we must derive the project root from `$PSScriptRoot`, the process path, or the current directory. Failing to do so breaks module imports.

## Unicode Glyph Parsing Failures
- `src/Render.psm1` contains glyphs like `☠`, `⚑`, and box-drawing characters. In WSL/PowerShell 7 the file loads correctly, but PS2EXE reads the UTF-8 (no BOM) file using the local ANSI code page, mangling the characters (`'âˆ¬`). The EXE then crashes before rendering.
- Fix options: re-save the PowerShell modules that use Unicode as UTF-8 **with BOM**, switch to ASCII glyphs only, or replace the literals with expressions (`[char]0x2620`). Without one of these changes the packaged `DungeonSweeper.exe` will always throw the “Unexpected token 'â'” parse error.

## Windows Terminal Launch Complications
- `wt.exe new-window …` only works if Windows Terminal (Store version) is installed and defaults to remembering prior tabs. Even with correct quoting for `pwsh -Command "& { … }"` it frequently opened extra tabs or failed outright on systems without `wt.exe`.
- Maintaining a reliable “new window” workflow adds platform-specific quirks, so we reverted to `make run` inside WSL as the supported entry point.

## Practical Path Forward
- Running the game via `make run` in WSL remains the most stable option.
- If packaging is revisited, convert Unicode-rich modules to UTF-8 with BOM (or ASCII) before invoking PS2EXE, and keep the compatibility fallbacks for console APIs.
- Leave the Windows Terminal helper optional; rely on the WSL loop for development and play.
