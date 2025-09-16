[CmdletBinding()]
param(
    [Parameter()] [switch] $Headless,
    [Parameter()] [int] $Seed,
    [Parameter()] [switch] $Compat
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$scriptRootCandidates = @()
if ($PSScriptRoot) {
    $scriptRootCandidates += $PSScriptRoot
}
if ($PSCommandPath) {
    $scriptRootCandidates += (Split-Path -Parent $PSCommandPath)
}

$processPath = $null
try {
    $processPath = Split-Path -Parent ([System.Diagnostics.Process]::GetCurrentProcess().MainModule.FileName)
} catch { }
if ($processPath) {
    $scriptRootCandidates += $processPath
}

$scriptRoot = $null
foreach ($candidate in ($scriptRootCandidates + (Get-Location).Path | Select-Object -Unique)) {
    if (-not $candidate) { continue }
    $candidateFull = [System.IO.Path]::GetFullPath($candidate)
    if (Test-Path (Join-Path $candidateFull 'src')) {
        $scriptRoot = $candidateFull
        break
    }
    $parent = Split-Path -Parent $candidateFull
    if ($parent -and (Test-Path (Join-Path $parent 'src'))) {
        $scriptRoot = $parent
        break
    }
}

if (-not $scriptRoot) {
    $scriptRoot = [System.IO.Path]::GetFullPath((Get-Location).Path)
}

$srcPath = Join-Path $scriptRoot 'src'

# Ensure UTF-8 regardless of host
$utf8Encoding = [System.Text.UTF8Encoding]::new($false)
$encodingSet = $false
try {
    [System.Console]::OutputEncoding = $utf8Encoding
    [System.Console]::InputEncoding = $utf8Encoding
    $encodingSet = $true
} catch {
    $encodingSet = $false
}

$modules = 'Board','Gen','Render','Input','Game'
foreach ($module in $modules) {
    $modulePath = Join-Path $srcPath ("{0}.psm1" -f $module)
    if (-not (Test-Path -LiteralPath $modulePath)) {
        throw "Required module '$modulePath' was not found."
    }
    Import-Module -Name $modulePath -Force
}

if (-not $PSBoundParameters.ContainsKey('Seed')) {
    $Seed = [System.Random]::new().Next()
}

$useCompat = [bool]$Compat
if (-not $encodingSet) {
    $useCompat = $true
}

$gameOptions = [ordered]@{
    Seed   = $Seed
    Compat = $useCompat
}

if ($Headless) {
    Invoke-DungeonSweeperHeadless @gameOptions
    return
}

$cursorWasVisible = $null
$treatControlC = $null
try { $treatControlC = [System.Console]::TreatControlCAsInput } catch { }
try {
    try {
        $cursorWasVisible = [System.Console]::CursorVisible
    } catch {
        $cursorWasVisible = $null
    }

    if ($cursorWasVisible) {
        try { [System.Console]::CursorVisible = $false } catch { $cursorWasVisible = $null }
    }

    try {
        [System.Console]::TreatControlCAsInput = $true
    } catch { }

    $null = Invoke-DungeonSweeperGame @gameOptions
}
finally {
    try { [System.Console]::TreatControlCAsInput = $treatControlC } catch { }
    if ($cursorWasVisible) {
        try { [System.Console]::CursorVisible = $true } catch { }
    }
}
