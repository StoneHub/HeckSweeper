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
        LastAction = 'Explore the dungeon. Press Space to reveal.'
    }
}

function Get-DungeonAction {
    param([Parameter(Mandatory)] [System.ConsoleKeyInfo] $KeyInfo)

    if ($KeyInfo.Key -eq [System.ConsoleKey]::C -and ($KeyInfo.Modifiers -band [System.ConsoleModifiers]::Control)) {
        return [pscustomobject]@{ Type = 'Quit' }
    }

    switch ($KeyInfo.Key) {
        ([System.ConsoleKey]::LeftArrow) { return [pscustomobject]@{ Type = 'Move'; Dx = -1; Dy = 0 } }
        ([System.ConsoleKey]::RightArrow) { return [pscustomobject]@{ Type = 'Move'; Dx = 1; Dy = 0 } }
        ([System.ConsoleKey]::UpArrow) { return [pscustomobject]@{ Type = 'Move'; Dx = 0; Dy = -1 } }
        ([System.ConsoleKey]::DownArrow) { return [pscustomobject]@{ Type = 'Move'; Dx = 0; Dy = 1 } }
        ([System.ConsoleKey]::A) { return [pscustomobject]@{ Type = 'Move'; Dx = -1; Dy = 0 } }
        ([System.ConsoleKey]::D) { return [pscustomobject]@{ Type = 'Move'; Dx = 1; Dy = 0 } }
        ([System.ConsoleKey]::W) { return [pscustomobject]@{ Type = 'Move'; Dx = 0; Dy = -1 } }
        ([System.ConsoleKey]::S) { return [pscustomobject]@{ Type = 'Move'; Dx = 0; Dy = 1 } }
        ([System.ConsoleKey]::H) { return [pscustomobject]@{ Type = 'Move'; Dx = -1; Dy = 0 } }
        ([System.ConsoleKey]::L) { return [pscustomobject]@{ Type = 'Move'; Dx = 1; Dy = 0 } }
        ([System.ConsoleKey]::K) { return [pscustomobject]@{ Type = 'Move'; Dx = 0; Dy = -1 } }
        ([System.ConsoleKey]::J) { return [pscustomobject]@{ Type = 'Move'; Dx = 0; Dy = 1 } }
        ([System.ConsoleKey]::Spacebar) { return [pscustomobject]@{ Type = 'Reveal' } }
        ([System.ConsoleKey]::Enter) { return [pscustomobject]@{ Type = 'Reveal' } }
        ([System.ConsoleKey]::F) { return [pscustomobject]@{ Type = 'Flag' } }
        ([System.ConsoleKey]::Q) { return [pscustomobject]@{ Type = 'Quit' } }
        ([System.ConsoleKey]::Escape) { return [pscustomobject]@{ Type = 'Quit' } }
        Default { return $null }
    }
}

function Move-DungeonCursor {
    param(
        [Parameter(Mandatory)] $State,
        [Parameter(Mandatory)] [int] $Dx,
        [Parameter(Mandatory)] [int] $Dy
    )

    $board = $State.Board
    $newX = [Math]::Max([Math]::Min($State.Cursor.X + $Dx, $board.Width - 1), 0)
    $newY = [Math]::Max([Math]::Min($State.Cursor.Y + $Dy, $board.Height - 1), 0)

    if ($newX -eq $State.Cursor.X -and $newY -eq $State.Cursor.Y) {
        $State.LastAction = 'Reached the edge of the map.'
        return
    }

    $State.Cursor.X = $newX
    $State.Cursor.Y = $newY
    $State.LastAction = "Cursor moved to ({0},{1})." -f ($newX + 1), ($newY + 1)
}

function Reveal-AllMonsters {
    param([Parameter(Mandatory)] $Board)

    for ($i = 0; $i -lt $Board.Monsters.Length; $i++) {
        if ($Board.Monsters[$i]) {
            $Board.Revealed[$i] = $true
        }
    }
}

function Resolve-DungeonReveal {
    param([Parameter(Mandatory)] $State)

    $board = $State.Board
    $index = ConvertTo-DungeonIndex -Width $board.Width -Height $board.Height -X $State.Cursor.X -Y $State.Cursor.Y

    if ($board.Flags[$index]) {
        $State.LastAction = 'Remove the flag before revealing.'
        return
    }

    if ($board.Revealed[$index]) {
        $State.LastAction = 'Already revealed.'
        return
    }

    $result = Reveal-DungeonCell -Board $board -Index $index
    if ($result.HitMonster) {
        $State.Status = 'Lost'
        $State.LastAction = 'A monster awakens!'
        Reveal-AllMonsters -Board $board
        return
    }

    if ($result.RevealedIndexes.Count -le 0) {
        $State.LastAction = 'No new rooms cleared.'
    } elseif ($result.RevealedIndexes.Count -eq 1) {
        $State.LastAction = 'Cleared 1 room.'
    } else {
        $State.LastAction = "Cleared {0} rooms." -f $result.RevealedIndexes.Count
    }

    if ($result.IsComplete) {
        $State.Status = 'Won'
        $State.LastAction = 'All safe rooms cleared!'
    }
}

function Resolve-DungeonFlag {
    param([Parameter(Mandatory)] $State)

    $board = $State.Board
    $index = ConvertTo-DungeonIndex -Width $board.Width -Height $board.Height -X $State.Cursor.X -Y $State.Cursor.Y

    if ($board.Revealed[$index]) {
        $State.LastAction = 'Cannot flag a revealed room.'
        return
    }

    $result = Toggle-DungeonFlag -Board $board -Index $index
    if ($result.IsFlagged) {
        $State.LastAction = 'Flag placed.'
    } else {
        $State.LastAction = 'Flag removed.'
    }
}

function Invoke-DungeonSweeperGame {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Seed,
        [Parameter(Mandatory)] [bool] $Compat
    )

    $state = New-DungeonSweeperState -Seed $Seed -Compat $Compat
    $renderer = New-DungeonRenderer -Width $state.Board.Width -Height $state.Board.Height -Compat $Compat

    while ($true) {
        Write-DungeonFrame -Renderer $renderer -State $state

        if ($state.Status -ne 'Playing') {
            break
        }

        $keyInfo = Read-DungeonInput
        if ($null -eq $keyInfo) { continue }

        $action = Get-DungeonAction -KeyInfo $keyInfo
        if ($null -eq $action) {
            continue
        }

        switch ($action.Type) {
            'Move'   { Move-DungeonCursor -State $state -Dx $action.Dx -Dy $action.Dy }
            'Reveal' { Resolve-DungeonReveal -State $state }
            'Flag'   { Resolve-DungeonFlag -State $state }
            'Quit'   { $state.Status = 'Quit'; $state.LastAction = 'Retreat called. Thanks for playing.' }
        }
    }

    Write-DungeonFrame -Renderer $renderer -State $state
    return $state
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
