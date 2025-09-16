Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

if ($PSStyle -and $PSStyle.OutputRendering) {
    $PSStyle.OutputRendering = 'Ansi'
}

$utf8 = [System.Text.UTF8Encoding]::new($false)
[System.Console]::OutputEncoding = $utf8
[System.Console]::InputEncoding = $utf8

$wtSession = $env:WT_SESSION
if (-not $wtSession) {
    Write-Host 'Note: Terminal does not report WT_SESSION; renderer should fall back to 16-color mode if needed.'
}
