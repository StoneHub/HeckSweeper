param(
    [int] $Seed = 123,
    [switch] $Compat
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$bootstrap = Join-Path $repoRoot 'scripts' 'bootstrap.ps1'
$entry = Join-Path $repoRoot 'DungeonSweeper.ps1'

if (Test-Path -LiteralPath $bootstrap) {
    try { & $bootstrap | Out-Null } catch { Write-Verbose $_.Exception.Message }
}

$result = & $entry -Headless -Seed $Seed -Compat:$Compat
$result | ConvertTo-Json -Depth 4
