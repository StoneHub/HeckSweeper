Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function Read-DungeonInput {
    [CmdletBinding()]
    param(
        [Parameter()] [switch] $NonBlocking
    )

    if ($NonBlocking) {
        if (-not [System.Console]::KeyAvailable) {
            return $null
        }
    }

    return [System.Console]::ReadKey($true)
}

Export-ModuleMember -Function Read-DungeonInput
