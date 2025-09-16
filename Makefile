SHELL := /bin/bash
WSL_PWSH := pwsh
WIN_PWSH := /mnt/c/Windows/System32/WindowsPowerShell/v1.0/powershell.exe
BUILD_DIR := build
EXE := $(BUILD_DIR)/DungeonSweeper.exe
WSL_PATH := $(shell pwd)
WSL_DISTRO := $(shell printf "%s" "$$WSL_DISTRO_NAME")
ifeq ($(WSL_DISTRO),)
WSL_DISTRO := $(shell /mnt/c/Windows/System32/wsl.exe -l --quiet 2>/dev/null | head -n1)
endif

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
	@"$(WIN_PWSH)" -NoLogo -NoProfile -ExecutionPolicy Bypass -File scripts/run-window.ps1 -Distro "$(WSL_DISTRO)" -WorkingDir "$(WSL_PATH)"

# Deterministic smoke test for CI
test: setup-wsl
	@$(WSL_PWSH) -NoLogo -NoProfile -ExecutionPolicy Bypass -File ./scripts/test-headless.ps1

package: $(EXE)

$(EXE): DungeonSweeper.ps1 scripts/package.ps1 | $(BUILD_DIR) setup-win
	@"$(WIN_PWSH)" -NoLogo -NoProfile -ExecutionPolicy Bypass -File scripts/package.ps1 -Input ./DungeonSweeper.ps1 -Output $(EXE)
	@echo "Packed -> $(EXE)"

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)

clean:
	rm -rf $(BUILD_DIR)
