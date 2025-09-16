param(
  [Parameter(Mandatory)] [string] $InFile,
  [Parameter(Mandatory)] [string] $OutFile
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

if ([string]::IsNullOrWhiteSpace($InFile)) {
  throw 'Input parameter was empty.'
}
if ([string]::IsNullOrWhiteSpace($OutFile)) {
  throw 'Output parameter was empty.'
}

$scriptRoot = Split-Path -Parent $PSCommandPath
$repoRoot = [System.IO.Path]::GetFullPath((Join-Path $scriptRoot '..'))

function Resolve-RepoPath {
  param([Parameter(Mandatory)][string] $Path)

  if ([System.IO.Path]::IsPathRooted($Path)) {
    return [System.IO.Path]::GetFullPath($Path)
  }

  $normalized = $Path.Trim()
  $normalized = $normalized -replace '/', '\\'

  return [System.IO.Path]::GetFullPath((Join-Path $repoRoot $normalized))
}

function Ensure-PS2EXE {
  if (Get-Command Invoke-ps2exe -ErrorAction SilentlyContinue) { return }
  if (-not (Get-Module -ListAvailable -Name PowerShellGet)) {
    throw 'PowerShellGet module not available. Install it from https://aka.ms/install-powershellget before packaging.'
  }
  Import-Module PowerShellGet -ErrorAction Stop
  try {
    Install-Module PS2EXE -Scope CurrentUser -Force -AllowClobber
  } catch {
    throw 'Installing PS2EXE failed: ' + $_.Exception.Message
  }
}

$inputPath = Resolve-RepoPath -Path $InFile
if (-not (Test-Path -LiteralPath $inputPath)) {
  throw "Input file '$inputPath' not found."
}

$outputPath = Resolve-RepoPath -Path $OutFile

Ensure-PS2EXE

$destination = Split-Path -Parent $outputPath
if (-not (Test-Path -LiteralPath $destination)) {
  $null = New-Item -ItemType Directory -Force -Path $destination
}

Invoke-ps2exe -inputFile $inputPath -outputFile $outputPath -noConsole -title 'DungeonSweeper' -verbose
Write-Host "OK: $outputPath"
