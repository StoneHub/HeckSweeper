Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function New-DungeonRenderer {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [switch] $Compat
    )

    [pscustomobject]@{
        Compat = [bool]$Compat
    }
}

function Write-DungeonFrame {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] $Renderer,
        [Parameter(Mandatory)] $Board,
        [Parameter(Mandatory)] $Cursor
    )
}

Export-ModuleMember -Function New-DungeonRenderer, Write-DungeonFrame
