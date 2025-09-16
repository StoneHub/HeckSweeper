Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function New-DungeonLayout {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Width,
        [Parameter(Mandatory)] [int] $Height,
        [Parameter(Mandatory)] [int] $Seed,
        [Parameter()] [int] $MonsterCount,
        [Parameter()] [double] $Density = 0.18
    )

    if ($Width -le 1 -or $Height -le 1) {
        throw 'Board dimensions must be greater than 1x1.'
    }

    $totalCells = $Width * $Height
    if (-not $PSBoundParameters.ContainsKey('MonsterCount')) {
        $rawCount = [Math]::Round($totalCells * $Density)
        $MonsterCount = [Math]::Max([Math]::Min([int]$rawCount, $totalCells - 1), 1)
    }

    if ($MonsterCount -ge $totalCells) {
        throw 'Monster count must leave at least one safe cell.'
    }

    $random = [System.Random]::new($Seed)
    $chosen = [System.Collections.Generic.HashSet[int]]::new()
    while ($chosen.Count -lt $MonsterCount) {
        $index = $random.Next(0, $totalCells)
        $null = $chosen.Add($index)
    }

    $monsterPositions = $chosen.ToArray()
    [Array]::Sort($monsterPositions)

    [pscustomobject]@{
        Seed           = $Seed
        Width          = $Width
        Height         = $Height
        MonsterCount   = $MonsterCount
        MonsterIndexes = $monsterPositions
    }
}

Export-ModuleMember -Function New-DungeonLayout
