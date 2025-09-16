param(
    [string] $Distro,
    [string] $WorkingDir
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

$bootstrapCommand = "& ./scripts/bootstrap.ps1; ./DungeonSweeper.ps1"
$wslArgs = @()
$wslArgs += '-d'
$wslArgs += $Distro
$wslArgs += '--cd'
$wslArgs += $WorkingDir
$wslArgs += 'pwsh'
$wslArgs += '-NoLogo'
$wslArgs += '-NoProfile'
$wslArgs += '-ExecutionPolicy'
$wslArgs += 'Bypass'
$wslArgs += '-Command'
$wslArgs += $bootstrapCommand

$wtArgs = @('new-window', 'wsl.exe') + $wslArgs
Start-Process -FilePath $wt.Path -ArgumentList $wtArgs | Out-Null
Write-Host 'Launched DungeonSweeper in a new Windows Terminal window.'
