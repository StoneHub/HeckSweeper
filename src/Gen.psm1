Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function New-DungeonLayout {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Width,
        [Parameter(Mandatory)] [int] $Height,
        [Parameter(Mandatory)] [int] $MonsterCount,
        [Parameter(Mandatory)] [System.Random] $Random
    )

    [pscustomobject]@{
        Width        = $Width
        Height       = $Height
        MonsterCount = $MonsterCount
        Random       = $Random
    }
}

Export-ModuleMember -Function New-DungeonLayout
