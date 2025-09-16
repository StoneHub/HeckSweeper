SHELL := /bin/bash
WSL_PWSH := pwsh
WIN_PWSH := /mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe
BUILD_DIR := build
EXE := $(BUILD_DIR)/DungeonSweeper.exe
WSL_REPO := $(shell pwd)
WIN_REPO := $(shell wslpath -w $(WSL_REPO))
WIN_INPUT := $(shell wslpath -w $(WSL_REPO)/DungeonSweeper.ps1)
WIN_OUTPUT := $(shell wslpath -w $(WSL_REPO)/$(EXE))
WIN_PACKAGE_SCRIPT := $(shell wslpath -w $(WSL_REPO)/scripts/package.ps1)
WIN_RUN_WINDOW_SCRIPT := $(shell wslpath -w $(WSL_REPO)/scripts/run-window.ps1)

.PHONY: all run run-window test package clean setup-wsl setup-win

all: run

setup-wsl:
	@if ! command -v $(WSL_PWSH) >/dev/null 2>&1; then \
	  echo "PowerShell 7 (pwsh) is required in WSL"; \
	  exit 1; \
	fi
	@$(WSL_PWSH) -NoLogo -NoProfile -ExecutionPolicy Bypass -File ./scripts/bootstrap.ps1
	@echo "WSL setup ok"

setup-win:
	@if [ ! -x "$(WIN_PWSH)" ]; then \
	  echo "powershell.exe not found at $(WIN_PWSH)"; \
	  exit 1; \
	fi
	@echo "Windows host PowerShell found"

run: setup-wsl
	@$(WSL_PWSH) -NoLogo -ExecutionPolicy Bypass -File ./DungeonSweeper.ps1

run-window: setup-wsl setup-win
	@"$(WIN_PWSH)" -NoLogo -NoProfile -ExecutionPolicy Bypass -File "$(WIN_RUN_WINDOW_SCRIPT)"

# Deterministic smoke test for CI
test: setup-wsl
	@$(WSL_PWSH) -NoLogo -NoProfile -ExecutionPolicy Bypass -File ./scripts/test-headless.ps1

package: $(EXE)

$(EXE): DungeonSweeper.ps1 scripts/package.ps1 | $(BUILD_DIR) setup-win
	@"$(WIN_PWSH)" -NoLogo -NoProfile -ExecutionPolicy Bypass -File "$(WIN_PACKAGE_SCRIPT)" -InFile "$(WIN_INPUT)" -OutFile "$(WIN_OUTPUT)"
	@echo "Packed -> $(EXE)"

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

clean:
	rm -rf $(BUILD_DIR)
