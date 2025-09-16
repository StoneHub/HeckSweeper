[CmdletBinding()]
param(
    [Parameter()] [switch] $Headless,
    [Parameter()] [int] $Seed,
    [Parameter()] [switch] $Compat
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$scriptRoot = Split-Path -Parent $PSCommandPath
$srcPath = Join-Path $scriptRoot 'src'

# Ensure UTF-8 regardless of locale
if ($PSVersionTable.PSVersion.Major -ge 6) {
    $encoding = [System.Text.UTF8Encoding]::new($false)
    if (-not [System.Console]::OutputEncoding.WebName.Equals($encoding.WebName)) {
        [System.Console]::OutputEncoding = $encoding
        [System.Console]::InputEncoding = $encoding
    }
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

$gameOptions = [ordered]@{
    Seed   = $Seed
    Compat = [bool]$Compat
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
