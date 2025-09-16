param(
    [string] $Distro,
    [string] $WorkingDir,
    [string] $Title = 'DungeonSweeper'
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

if (-not $Distro) {
    if ($env:WSL_DISTRO_NAME) {
        $Distro = $env:WSL_DISTRO_NAME
    } else {
        throw 'WSL distribution name not provided. Pass -Distro or ensure WSL_DISTRO_NAME is set.'
    }
}

if (-not $WorkingDir) {
    throw 'Working directory must be provided.'
}

$wt = Get-Command wt.exe -ErrorAction SilentlyContinue
if (-not $wt) {
    throw 'Windows Terminal (wt.exe) not found. Install Windows Terminal or adjust scripts/run-window.ps1.'
}

$escapedCommand = '"& { ./scripts/bootstrap.ps1; ./DungeonSweeper.ps1 }"'

$wtArgs = @(
    'new-window',
    '--title', $Title,
    '--',
    'wsl.exe',
    '-d', $Distro,
    '--cd', $WorkingDir,
    'pwsh',
    '-NoLogo',
    '-NoProfile',
    '-ExecutionPolicy', 'Bypass',
    '-Command', $escapedCommand
)

Start-Process -FilePath $wt.Path -ArgumentList $wtArgs | Out-Null
Write-Host "Launched DungeonSweeper in a new Windows Terminal window (distro=$Distro, cwd=$WorkingDir)."
