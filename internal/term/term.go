package term

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// ANSI escape sequences
const (
	ESC = "\x1b"

	// Cursor control
	CursorHide    = ESC + "[?25l"
	CursorShow    = ESC + "[?25h"
	CursorHome    = ESC + "[H"
	CursorSave    = ESC + "[s"
	CursorRestore = ESC + "[u"

	// Screen control
	ClearScreen     = ESC + "[2J"
	ClearLine       = ESC + "[2K"
	ClearToEnd      = ESC + "[J"
	ClearLineToEnd  = ESC + "[K"

	// Colors
	Reset     = ESC + "[0m"
	Bold      = ESC + "[1m"
	Dim       = ESC + "[2m"
	Italic    = ESC + "[3m"
	Underline = ESC + "[4m"
	Blink     = ESC + "[5m"
	Reverse   = ESC + "[7m"

	// Foreground colors
	FgBlack   = ESC + "[30m"
	FgRed     = ESC + "[31m"
	FgGreen   = ESC + "[32m"
	FgYellow  = ESC + "[33m"
	FgBlue    = ESC + "[34m"
	FgMagenta = ESC + "[35m"
	FgCyan    = ESC + "[36m"
	FgWhite   = ESC + "[37m"
	FgDefault = ESC + "[39m"

	// Bright foreground colors
	FgBrightBlack   = ESC + "[90m"
	FgBrightRed     = ESC + "[91m"
	FgBrightGreen   = ESC + "[92m"
	FgBrightYellow  = ESC + "[93m"
	FgBrightBlue    = ESC + "[94m"
	FgBrightMagenta = ESC + "[95m"
	FgBrightCyan    = ESC + "[96m"
	FgBrightWhite   = ESC + "[97m"

	// Background colors
	BgBlack   = ESC + "[40m"
	BgRed     = ESC + "[41m"
	BgGreen   = ESC + "[42m"
	BgYellow  = ESC + "[43m"
	BgBlue    = ESC + "[44m"
	BgMagenta = ESC + "[45m"
	BgCyan    = ESC + "[46m"
	BgWhite   = ESC + "[47m"
	BgDefault = ESC + "[49m"
)

// HideCursor hides the terminal cursor
func HideCursor(w io.Writer) {
	fmt.Fprint(w, CursorHide)
}

// ShowCursor shows the terminal cursor
func ShowCursor(w io.Writer) {
	fmt.Fprint(w, CursorShow)
}

// ClearScreenFull clears the entire screen
func ClearScreenFull(w io.Writer) {
	fmt.Fprint(w, ClearScreen)
}

// MoveHome moves cursor to top-left
func MoveHome(w io.Writer) {
	fmt.Fprint(w, CursorHome)
}

// MoveTo moves cursor to specific position (1-indexed)
func MoveTo(w io.Writer, row, col int) {
	fmt.Fprintf(w, ESC+"[%d;%dH", row, col)
}

// MoveToZero moves cursor to specific position (0-indexed)
func MoveToZero(w io.Writer, row, col int) {
	fmt.Fprintf(w, ESC+"[%d;%dH", row+1, col+1)
}

// SetFgColor sets foreground color using 256-color mode
func SetFgColor(w io.Writer, color int) {
	fmt.Fprintf(w, ESC+"[38;5;%dm", color)
}

// SetBgColor sets background color using 256-color mode
func SetBgColor(w io.Writer, color int) {
	fmt.Fprintf(w, ESC+"[48;5;%dm", color)
}

// SetFgRGB sets foreground color using 24-bit RGB
func SetFgRGB(w io.Writer, r, g, b int) {
	fmt.Fprintf(w, ESC+"[38;2;%d;%d;%dm", r, g, b)
}

// SetBgRGB sets background color using 24-bit RGB
func SetBgRGB(w io.Writer, r, g, b int) {
	fmt.Fprintf(w, ESC+"[48;2;%d;%d;%dm", r, g, b)
}

// ResetStyle resets all text attributes
func ResetStyle(w io.Writer) {
	fmt.Fprint(w, Reset)
}

// Size returns the terminal size (width, height)
func Size() (int, int, error) {
	if runtime.GOOS == "windows" {
		return sizeWindows()
	}
	return sizeUnix()
}

func sizeUnix() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		// Fallback to tput
		return sizeTput()
	}
	parts := strings.Fields(string(out))
	if len(parts) != 2 {
		return sizeTput()
	}
	rows, err := strconv.Atoi(parts[0])
	if err != nil {
		return sizeTput()
	}
	cols, err := strconv.Atoi(parts[1])
	if err != nil {
		return sizeTput()
	}
	return cols, rows, nil
}

func sizeTput() (int, int, error) {
	colsCmd := exec.Command("tput", "cols")
	colsCmd.Stdin = os.Stdin
	colsOut, err := colsCmd.Output()
	if err != nil {
		return 80, 24, nil // Default fallback
	}
	rowsCmd := exec.Command("tput", "lines")
	rowsCmd.Stdin = os.Stdin
	rowsOut, err := rowsCmd.Output()
	if err != nil {
		return 80, 24, nil
	}
	cols, _ := strconv.Atoi(strings.TrimSpace(string(colsOut)))
	rows, _ := strconv.Atoi(strings.TrimSpace(string(rowsOut)))
	if cols == 0 {
		cols = 80
	}
	if rows == 0 {
		rows = 24
	}
	return cols, rows, nil
}

func sizeWindows() (int, int, error) {
	cmd := exec.Command("cmd", "/c", "mode", "con")
	out, err := cmd.Output()
	if err != nil {
		return 80, 24, nil
	}
	lines := strings.Split(string(out), "\n")
	var cols, rows int
	for _, line := range lines {
		line = strings.ToLower(line)
		if strings.Contains(line, "columns") || strings.Contains(line, "cols") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				cols, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
			}
		}
		if strings.Contains(line, "lines") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				rows, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
			}
		}
	}
	if cols == 0 {
		cols = 80
	}
	if rows == 0 {
		rows = 24
	}
	return cols, rows, nil
}

// Bell sounds the terminal bell
func Bell(w io.Writer) {
	fmt.Fprint(w, "\a")
}

// EnableAlternateScreen switches to alternate screen buffer
func EnableAlternateScreen(w io.Writer) {
	fmt.Fprint(w, ESC+"[?1049h")
}

// DisableAlternateScreen returns to main screen buffer
func DisableAlternateScreen(w io.Writer) {
	fmt.Fprint(w, ESC+"[?1049l")
}

// EnableMouse enables mouse tracking
func EnableMouse(w io.Writer) {
	fmt.Fprint(w, ESC+"[?1000h"+ESC+"[?1002h"+ESC+"[?1015h"+ESC+"[?1006h")
}

// DisableMouse disables mouse tracking
func DisableMouse(w io.Writer) {
	fmt.Fprint(w, ESC+"[?1006l"+ESC+"[?1015l"+ESC+"[?1002l"+ESC+"[?1000l")
}
