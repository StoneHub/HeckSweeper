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
$treatControlC = [System.Console]::TreatControlCAsInput
try {
    if ([System.Console]::CursorVisible) {
        $cursorWasVisible = $true
        [System.Console]::CursorVisible = $false
    } else {
        $cursorWasVisible = $false
    }
    [System.Console]::TreatControlCAsInput = $true

    $null = Invoke-DungeonSweeperGame @gameOptions
}
finally {
    [System.Console]::TreatControlCAsInput = $treatControlC
    if ($cursorWasVisible -and -not [System.Console]::CursorVisible) {
        try { [System.Console]::CursorVisible = $true } catch { }
    }
}
