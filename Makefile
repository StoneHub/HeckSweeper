SHELL := /bin/bash
WSL_PWSH := pwsh
WIN_PWSH := /mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe
BUILD_DIR := build
EXE := $(BUILD_DIR)/DungeonSweeper.exe

.PHONY: all run test package clean setup-wsl setup-win

all: run

setup-wsl:
	@if ! command -v $(WSL_PWSH) >/dev/null 2>&1; then \
	  echo "PowerShell 7 (pwsh) is required in WSL"; \
	  exit 1; \
	fi
	@$(WSL_PWSH) -NoLogo -NoProfile -Command "\
	  $$ErrorActionPreference='Stop'; \
	  [Console]::OutputEncoding=[Text.UTF8Encoding]::new($$false) \
	"
	@echo "WSL setup ok"

setup-win:
	@if [ ! -x "$(WIN_PWSH)" ]; then \
	  echo "powershell.exe not found at $(WIN_PWSH)"; \
	  exit 1; \
	fi
	@echo "Windows host PowerShell found"

run: setup-wsl
	@$(WSL_PWSH) -NoLogo -ExecutionPolicy Bypass -File ./DungeonSweeper.ps1

# Deterministic smoke test for CI
test: setup-wsl
	@$(WSL_PWSH) -NoLogo -NoProfile -ExecutionPolicy Bypass -Command "\
	  $$ErrorActionPreference='Stop'; \
	  $$result = & ./DungeonSweeper.ps1 -Headless -Seed 123; \
	  $$result | ConvertTo-Json -Depth 4 \
	"

package: $(EXE)

$(EXE): DungeonSweeper.ps1 scripts/package.ps1 | $(BUILD_DIR) setup-win
	@"$(WIN_PWSH)" -NoLogo -NoProfile -ExecutionPolicy Bypass -File scripts/package.ps1 -Input ./DungeonSweeper.ps1 -Output $(EXE)
	@echo "Packed -> $(EXE)"

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

clean:
	rm -rf $(BUILD_DIR)
