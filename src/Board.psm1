Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function New-DungeonBoardState {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Width,
        [Parameter(Mandatory)] [int] $Height,
        [Parameter(Mandatory)] [object[]] $Cells
    )

    [pscustomobject]@{
        Width  = $Width
        Height = $Height
        Cells  = $Cells
    }
}

Export-ModuleMember -Function New-DungeonBoardState
