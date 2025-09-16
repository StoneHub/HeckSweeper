Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

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

    Write-Verbose "Running headless diagnostics (stub) with Seed=$Seed Compat=$Compat"
    [pscustomobject]@{
        Seed   = $Seed
        Compat = $Compat
        Status = 'NotImplemented'
    }
}

Export-ModuleMember -Function Invoke-DungeonSweeperGame, Invoke-DungeonSweeperHeadless
