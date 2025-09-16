Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$script:DefaultWidth = 24
$script:DefaultHeight = 16
$script:DefaultSeedDensity = 0.18

function New-DungeonSweeperState {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Seed,
        [Parameter(Mandatory)] [bool] $Compat
    )

    $layout = New-DungeonLayout -Width $script:DefaultWidth -Height $script:DefaultHeight -Seed $Seed -Density $script:DefaultSeedDensity
    $board = New-DungeonBoard -Layout $layout
    $cursor = [pscustomobject]@{ X = [int]($layout.Width / 2); Y = [int]($layout.Height / 2) }

    [pscustomobject]@{
        Seed       = $Seed
        Compat     = $Compat
        Layout     = $layout
        Board      = $board
        Cursor     = $cursor
        Status     = 'Playing'
        LastAction = $null
    }
}

function Invoke-DungeonSweeperGame {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Seed,
        [Parameter(Mandatory)] [bool] $Compat
    )

    throw 'Interactive game loop not yet implemented.'
}

function Invoke-DungeonSweeperHeadless {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Seed,
        [Parameter(Mandatory)] [bool] $Compat
    )

    $state = New-DungeonSweeperState -Seed $Seed -Compat $Compat
    $board = $state.Board
    $layout = $state.Layout
    $log = [System.Collections.Generic.List[object]]::new()

    $total = $board.Width * $board.Height
    $safeZero = -1
    $safeThreat = -1
    for ($i = 0; $i -lt $total; $i++) {
        if ($board.Monsters[$i]) { continue }
        if ($board.Threats[$i] -eq 0 -and $safeZero -eq -1) { $safeZero = $i }
        if ($board.Threats[$i] -gt 0 -and $safeThreat -eq -1) { $safeThreat = $i }
        if ($safeZero -ne -1 -and $safeThreat -ne -1) { break }
    }

    if ($safeZero -eq -1) {
        $safeZero = (0..($total - 1) | Where-Object { -not $board.Monsters[$_] } | Select-Object -First 1)
    }

    if ($safeThreat -eq -1) {
        $safeThreat = (0..($total - 1) | Where-Object { -not $board.Monsters[$_] } | Select-Object -First 1)
    }

    $targets = [System.Collections.Generic.List[int]]::new()
    if ($safeZero -ge 0) { $targets.Add($safeZero) }
    if ($safeThreat -ge 0 -and $safeThreat -ne $safeZero) { $targets.Add($safeThreat) }

    foreach ($target in $targets) {
        $result = Reveal-DungeonCell -Board $board -Index $target
        $coord = ConvertTo-DungeonCoordinate -Width $board.Width -Height $board.Height -Index $target
        $log.Add([pscustomobject]@{
            Action          = 'Reveal'
            Index           = $target
            Coordinate      = $coord
            RevealedCount   = $result.RevealedIndexes.Count
            HitMonster      = $result.HitMonster
            RemainingSafe   = $result.RemainingSafe
            Completed       = $result.IsComplete
        })
        if ($result.HitMonster) { break }
    }

    if ($layout.MonsterIndexes.Count -gt 0 -and -not $board.IsExploded) {
        $monsterIndex = $layout.MonsterIndexes[0]
        $coord = ConvertTo-DungeonCoordinate -Width $board.Width -Height $board.Height -Index $monsterIndex

        $flagged = Toggle-DungeonFlag -Board $board -Index $monsterIndex
        $log.Add([pscustomobject]@{
            Action     = 'Flag'
            Index      = $monsterIndex
            Coordinate = $coord
            IsFlagged  = $flagged.IsFlagged
        })

        $result = Reveal-DungeonCell -Board $board -Index $monsterIndex
        $log.Add([pscustomobject]@{
            Action          = 'Reveal-Flagged'
            Index           = $monsterIndex
            Coordinate      = $coord
            RevealedCount   = $result.RevealedIndexes.Count
            HitMonster      = $result.HitMonster
            RemainingSafe   = $result.RemainingSafe
            Completed       = $result.IsComplete
        })

        $unflagged = Toggle-DungeonFlag -Board $board -Index $monsterIndex
        $log.Add([pscustomobject]@{
            Action     = 'Unflag'
            Index      = $monsterIndex
            Coordinate = $coord
            IsFlagged  = $unflagged.IsFlagged
        })

        $result = Reveal-DungeonCell -Board $board -Index $monsterIndex
        $log.Add([pscustomobject]@{
            Action          = 'Reveal-Monster'
            Index           = $monsterIndex
            Coordinate      = $coord
            RevealedCount   = $result.RevealedIndexes.Count
            HitMonster      = $result.HitMonster
            RemainingSafe   = $result.RemainingSafe
            Completed       = $result.IsComplete
        })
    }

    $status = if ($board.IsExploded) { 'Exploded' } elseif ($board.RemainingSafe -eq 0) { 'Cleared' } else { 'InProgress' }

    [pscustomobject]@{
        Seed          = $Seed
        Compat        = $Compat
        Dimensions    = [pscustomobject]@{ Width = $board.Width; Height = $board.Height }
        Monsters      = $layout.MonsterCount
        RemainingSafe = $board.RemainingSafe
        ExplodedIndex = $board.ExplodedIndex
        Status        = $status
        Steps         = $log.ToArray()
    }
}

Export-ModuleMember -Function New-DungeonSweeperState, Invoke-DungeonSweeperGame, Invoke-DungeonSweeperHeadless
