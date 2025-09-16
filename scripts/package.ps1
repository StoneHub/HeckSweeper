param(
  [Parameter(Mandatory)] [string] $Input,
  [Parameter(Mandatory)] [string] $Output
)

$ErrorActionPreference = 'Stop'

function Ensure-PS2EXE {
  if (-not (Get-Command Invoke-ps2exe -ErrorAction SilentlyContinue)) {
    Write-Host 'Installing PS2EXE module for packaging...'
    Install-Module PS2EXE -Scope CurrentUser -Force -AllowClobber
  }
}

Ensure-PS2EXE
$destination = Split-Path -Parent (Resolve-Path -Path $Output -ErrorAction SilentlyContinue)
if (-not $destination) {
  $destination = Split-Path -Parent $Output
}
if (-not (Test-Path -LiteralPath $destination)) {
  $null = New-Item -ItemType Directory -Force -Path $destination
}

Invoke-ps2exe -inputFile $Input -outputFile $Output -noConsole -title 'DungeonSweeper' -verbose
Write-Host "OK: $Output"
