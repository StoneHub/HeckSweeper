Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function ConvertTo-DungeonIndex {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Width,
        [Parameter(Mandatory)] [int] $Height,
        [Parameter(Mandatory)] [int] $X,
        [Parameter(Mandatory)] [int] $Y
    )

    if ($X -lt 0 -or $X -ge $Width -or $Y -lt 0 -or $Y -ge $Height) {
        return -1
    }

    return ($Y * $Width) + $X
}

function ConvertTo-DungeonCoordinate {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Width,
        [Parameter(Mandatory)] [int] $Height,
        [Parameter(Mandatory)] [int] $Index
    )

    $total = $Width * $Height
    if ($Index -lt 0 -or $Index -ge $total) {
        throw "Index $Index is outside the board bounds."
    }

    $x = $Index % $Width
    $y = [int]($Index / $Width)

    [pscustomobject]@{ X = $x; Y = $y }
}

function Get-DungeonNeighborIndexes {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Width,
        [Parameter(Mandatory)] [int] $Height,
        [Parameter(Mandatory)] [int] $Index
    )

    $total = $Width * $Height
    if ($Index -lt 0 -or $Index -ge $total) {
        throw "Index $Index is outside the board bounds."
    }

    $x = $Index % $Width
    $y = [int]($Index / $Width)

    $neighbors = [System.Collections.Generic.List[int]]::new()
    foreach ($dy in -1..1) {
        foreach ($dx in -1..1) {
            if ($dx -eq 0 -and $dy -eq 0) { continue }
            $nx = $x + $dx
            $ny = $y + $dy
            if ($nx -lt 0 -or $ny -lt 0 -or $nx -ge $Width -or $ny -ge $Height) {
                continue
            }
            $neighbors.Add(($ny * $Width) + $nx)
        }
    }

    return $neighbors.ToArray()
}

function New-DungeonBoard {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] $Layout
    )

    $Width = $Layout.Width
    $Height = $Layout.Height
    $monsterIndexes = $Layout.MonsterIndexes
    $total = $Width * $Height

    $threats = [int[]]::new($total)
    $monsters = [bool[]]::new($total)
    $revealed = [bool[]]::new($total)
    $flags = [bool[]]::new($total)

    foreach ($idx in $monsterIndexes) {
        if ($idx -lt 0 -or $idx -ge $total) {
            throw "Monster index $idx is outside the board bounds."
        }
        $monsters[$idx] = $true
    }

    foreach ($idx in $monsterIndexes) {
        foreach ($neighbor in (Get-DungeonNeighborIndexes -Width $Width -Height $Height -Index $idx)) {
            if (-not $monsters[$neighbor]) {
                $threats[$neighbor]++
            }
        }
    }

    [pscustomobject]@{
        Width         = $Width
        Height        = $Height
        Threats       = $threats
        Monsters      = $monsters
        Revealed      = $revealed
        Flags         = $flags
        TotalSafe     = $total - $monsterIndexes.Count
        RemainingSafe = $total - $monsterIndexes.Count
        IsExploded    = $false
        ExplodedIndex = $null
    }
}

function Reveal-DungeonCell {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] $Board,
        [Parameter(Mandatory)] [int] $Index
    )

    $total = $Board.Width * $Board.Height
    if ($Index -lt 0 -or $Index -ge $total) {
        throw "Index $Index is outside the board bounds."
    }

    if ($Board.IsExploded) {
        return [pscustomobject]@{
            RevealedIndexes = @()
            HitMonster      = $true
            RemainingSafe   = $Board.RemainingSafe
            IsComplete      = $false
        }
    }

    if ($Board.Revealed[$Index]) {
        return [pscustomobject]@{
            RevealedIndexes = @()
            HitMonster      = $false
            RemainingSafe   = $Board.RemainingSafe
            IsComplete      = ($Board.RemainingSafe -eq 0)
        }
    }

    if ($Board.Flags[$Index]) {
        return [pscustomobject]@{
            RevealedIndexes = @()
            HitMonster      = $false
            RemainingSafe   = $Board.RemainingSafe
            IsComplete      = ($Board.RemainingSafe -eq 0)
        }
    }

    if ($Board.Monsters[$Index]) {
        $Board.Revealed[$Index] = $true
        $Board.IsExploded = $true
        $Board.ExplodedIndex = $Index
        return [pscustomobject]@{
            RevealedIndexes = @($Index)
            HitMonster      = $true
            RemainingSafe   = $Board.RemainingSafe
            IsComplete      = $false
        }
    }

    $queue = [System.Collections.Generic.Queue[int]]::new()
    $seen = [bool[]]::new($total)
    $revealedIndexes = [System.Collections.Generic.List[int]]::new()

    $queue.Enqueue($Index)
    $seen[$Index] = $true

    while ($queue.Count -gt 0) {
        $current = $queue.Dequeue()
        if ($Board.Revealed[$current]) { continue }
        if ($Board.Flags[$current]) { continue }

        $Board.Revealed[$current] = $true
        $revealedIndexes.Add($current)
        $Board.RemainingSafe--

        $currentThreat = $Board.Threats[$current]
        if ($currentThreat -eq 0) {
            foreach ($neighbor in (Get-DungeonNeighborIndexes -Width $Board.Width -Height $Board.Height -Index $current)) {
                if (-not $seen[$neighbor] -and -not $Board.Monsters[$neighbor]) {
                    $queue.Enqueue($neighbor)
                    $seen[$neighbor] = $true
                }
            }
        }
    }

    $isComplete = ($Board.RemainingSafe -le 0)

    [pscustomobject]@{
        RevealedIndexes = $revealedIndexes.ToArray()
        HitMonster      = $false
        RemainingSafe   = $Board.RemainingSafe
        IsComplete      = $isComplete
    }
}

function Toggle-DungeonFlag {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] $Board,
        [Parameter(Mandatory)] [int] $Index
    )

    $total = $Board.Width * $Board.Height
    if ($Index -lt 0 -or $Index -ge $total) {
        throw "Index $Index is outside the board bounds."
    }

    if ($Board.Revealed[$Index]) {
        return [pscustomobject]@{
            Index     = $Index
            IsFlagged = $Board.Flags[$Index]
        }
    }

    $Board.Flags[$Index] = -not $Board.Flags[$Index]

    [pscustomobject]@{
        Index     = $Index
        IsFlagged = $Board.Flags[$Index]
    }
}

Export-ModuleMember -Function `
    ConvertTo-DungeonIndex, `
    ConvertTo-DungeonCoordinate, `
    Get-DungeonNeighborIndexes, `
    New-DungeonBoard, `
    Reveal-DungeonCell, `
    Toggle-DungeonFlag
