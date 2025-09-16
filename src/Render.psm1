Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

function New-DungeonRenderer {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] [int] $Width,
        [Parameter(Mandatory)] [int] $Height,
        [Parameter(Mandatory)] [bool] $Compat
    )

    $symbols = if ($Compat) {
        [pscustomobject]@{
            Hidden    = '.'
            Flag      = 'F'
            Monster   = 'X'
            Horizontal= '-'
            Vertical  = '|'
            TopLeft   = '+'
            TopRight  = '+'
            BottomLeft= '+'
            BottomRight='+'
            CursorTag = '>'
            TitleDash = '-'
            ControlSep= '·'
        }
    } else {
        [pscustomobject]@{
            Hidden    = '·'
            Flag      = '⚑'
            Monster   = '☠'
            Horizontal= '─'
            Vertical  = '│'
            TopLeft   = '┌'
            TopRight  = '┐'
            BottomLeft= '└'
            BottomRight='┘'
            CursorTag = $null
            TitleDash = '–'
            ControlSep= '·'
        }
    }

    $cellSpan = $Width * 2
    $frameWidth = [Math]::Max($cellSpan + 6, 80)
    $topBorder = "  {0}{1}{2}" -f $symbols.TopLeft, ($symbols.Horizontal * $cellSpan), $symbols.TopRight
    $bottomBorder = "  {0}{1}{2}" -f $symbols.BottomLeft, ($symbols.Horizontal * $cellSpan), $symbols.BottomRight

    [pscustomobject]@{
        Compat       = $Compat
        BoardWidth   = $Width
        BoardHeight  = $Height
        Symbols      = $symbols
        FrameWidth   = $frameWidth
        FrameHeight  = $Height + 7
        TopBorder    = $topBorder
        BottomBorder = $bottomBorder
        Buffer       = [System.Text.StringBuilder]::new()
        CursorAnsi   = if ($Compat) { $null } else { "`e[48;5;239m" }
        ResetAnsi    = if ($Compat) { '' } else { "`e[0m" }
    }
}

function Write-DungeonFrame {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)] $Renderer,
        [Parameter(Mandatory)] $State
    )

    $board = $State.Board
    $symbols = $Renderer.Symbols
    $sb = $Renderer.Buffer
    $null = $sb.Clear()

    $titleDash = $symbols.TitleDash
    $title = "Dungeon Sweeper {0} Seed {1}" -f $titleDash, $State.Seed
    $flagCount = 0
    for ($i = 0; $i -lt $board.Flags.Length; $i++) { if ($board.Flags[$i]) { $flagCount++ } }
    $stats = "Remaining: {0}   Flags: {1}   Status: {2}" -f $board.RemainingSafe, $flagCount, $State.Status
    $message = if ($State.LastAction) { $State.LastAction } else { 'Use arrows or WASD to move. Space reveals, F toggles flag.' }

    $lines = [System.Collections.Generic.List[string]]::new()
    $lines.Add(($title.PadRight($Renderer.FrameWidth)))
    $lines.Add(($stats.PadRight($Renderer.FrameWidth)))
    $lines.Add(($message.PadRight($Renderer.FrameWidth)))
    $lines.Add(($Renderer.TopBorder.PadRight($Renderer.FrameWidth)))

    for ($y = 0; $y -lt $board.Height; $y++) {
        $lineBuilder = [System.Text.StringBuilder]::new()
        $null = $lineBuilder.Append('  ')
        $null = $lineBuilder.Append($symbols.Vertical)
        for ($x = 0; $x -lt $board.Width; $x++) {
            $index = ConvertTo-DungeonIndex -Width $board.Width -Height $board.Height -X $x -Y $y
            $isCursor = ($State.Cursor.X -eq $x -and $State.Cursor.Y -eq $y)
            $isFlagged = $board.Flags[$index]
            $isRevealed = $board.Revealed[$index]
            $isMonster = $board.Monsters[$index]
            $displayMonster = $isMonster -and ($isRevealed -or $State.Status -eq 'Lost')

            $glyph = if ($displayMonster) {
                $symbols.Monster
            } elseif ($isRevealed) {
                $threat = $board.Threats[$index]
                if ($threat -gt 0) { [string]$threat } else { ' ' }
            } elseif ($isFlagged) {
                $symbols.Flag
            } else {
                $symbols.Hidden
            }

            if ($Renderer.Compat) {
                if ($isCursor) {
                    $null = $lineBuilder.Append($symbols.CursorTag)
                } else {
                    $null = $lineBuilder.Append(' ')
                }
                $null = $lineBuilder.Append($glyph)
            } else {
                if ($isCursor) { $null = $lineBuilder.Append($Renderer.CursorAnsi) }
                $null = $lineBuilder.Append(' ')
                $null = $lineBuilder.Append($glyph)
                if ($isCursor) { $null = $lineBuilder.Append($Renderer.ResetAnsi) }
            }
        }
        $null = $lineBuilder.Append($symbols.Vertical)
        $lines.Add(($lineBuilder.ToString().PadRight($Renderer.FrameWidth)))
    }

    $lines.Add(($Renderer.BottomBorder.PadRight($Renderer.FrameWidth)))

    $controls = "Controls: Arrows/WASD move {0} Space reveal {0} F flag {0} Q quit" -f $symbols.ControlSep
    $lines.Add(($controls.PadRight($Renderer.FrameWidth)))

    $outcome = switch ($State.Status) {
        'Won'      { 'You cleared the dungeon! Press Q to exit.' }
        'Lost'     { 'The monsters feast... Press Q to exit.' }
        'Quit'     { 'Run abandoned. Thanks for playing.' }
        Default    { 'Hunt safely. Monsters remain hidden.' }
    }
    $lines.Add(($outcome.PadRight($Renderer.FrameWidth)))

    while ($lines.Count -lt $Renderer.FrameHeight) {
        $lines.Add(''.PadRight($Renderer.FrameWidth))
    }

    foreach ($line in $lines) {
        $null = $sb.AppendLine($line)
    }

    [System.Console]::SetCursorPosition(0, 0)
    [System.Console]::Write($sb.ToString())
}

Export-ModuleMember -Function New-DungeonRenderer, Write-DungeonFrame
