param(
    [int] $Seed = 123,
    [switch] $ShowMonsters,
    [switch] $Compat,
    [switch] $IncludeThreats = $true
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$scriptRoot = Split-Path -Parent $PSScriptRoot
$srcPath = Join-Path $scriptRoot 'src'

$modules = 'Gen','Board'
foreach ($name in $modules) {
    $modulePath = Join-Path $srcPath ("{0}.psm1" -f $name)
    Import-Module -Name $modulePath -Force
}

$layout = New-DungeonLayout -Width 24 -Height 16 -Seed $Seed
$board = New-DungeonBoard -Layout $layout

$rows = @()
for ($y = 0; $y -lt $board.Height; $y++) {
    $cells = @()
    for ($x = 0; $x -lt $board.Width; $x++) {
        $index = ConvertTo-DungeonIndex -Width $board.Width -Height $board.Height -X $x -Y $y
        if ($layout.MonsterIndexes -contains $index) {
            $cells += if ($ShowMonsters) { 'M' } else { '·' }
            continue
        }

        $threat = $board.Threats[$index]
        if ($threat -gt 0 -and $IncludeThreats) {
            $cells += [string]$threat
        } else {
            $cells += ' '
        }
    }
    $rows += ,$cells
}

Write-Host "Seed: $Seed"
Write-Host "Monsters: $($layout.MonsterCount)"
Write-Host "Threat grid (M=monster)"

for ($y = 0; $y -lt $board.Height; $y++) {
    $line = ''
    for ($x = 0; $x -lt $board.Width; $x++) {
        $value = $rows[$y][$x]
        if (($value -eq ' ') -and -not $ShowMonsters) {
            $value = '.'
        }
        $line += $value.PadLeft(2)
    }
    Write-Host $line
}
