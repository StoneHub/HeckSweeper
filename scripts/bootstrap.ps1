Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

if ($PSStyle -and $PSStyle.OutputRendering) {
    $PSStyle.OutputRendering = 'Ansi'
}

$utf8 = [System.Text.UTF8Encoding]::new($false)
try { [System.Console]::OutputEncoding = $utf8 } catch { }
try { [System.Console]::InputEncoding = $utf8 } catch { }

$wtSession = $env:WT_SESSION
if (-not $wtSession) {
    Write-Host 'Note: Terminal does not report WT_SESSION; renderer may choose ASCII fallback.'
}
