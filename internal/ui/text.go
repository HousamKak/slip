package ui

import (
	"strings"

	"slip/internal/engine"
)

// WrapText wraps text to fit within a maximum width
func WrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return nil
	}

	var lines []string
	paragraphs := strings.Split(text, "\n")

	for _, para := range paragraphs {
		if para == "" {
			lines = append(lines, "")
			continue
		}

		words := strings.Fields(para)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		currentLine := words[0]
		for _, word := range words[1:] {
			if len(currentLine)+1+len(word) <= maxWidth {
				currentLine += " " + word
			} else {
				lines = append(lines, currentLine)
				currentLine = word
			}
		}
		lines = append(lines, currentLine)
	}

	return lines
}

// CenterText centers text within a given width
func CenterText(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}
	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}

// PadRight pads text to the right to fill width
func PadRight(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}
	return text + strings.Repeat(" ", width-len(text))
}

// PadLeft pads text to the left to fill width
func PadLeft(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}
	return strings.Repeat(" ", width-len(text)) + text
}

// TruncateText truncates text and adds ellipsis if needed
func TruncateText(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	if maxLen <= 3 {
		return text[:maxLen]
	}
	return text[:maxLen-3] + "..."
}

// ProgressBar renders a progress bar
type ProgressBar struct {
	X, Y     int
	Width    int
	Progress float64 // 0.0 to 1.0
	Style    ProgressBarStyle
}

// ProgressBarStyle defines the visual style of a progress bar
type ProgressBarStyle struct {
	Border   engine.Style
	Fill     engine.Style
	Empty    engine.Style
	FillChar rune
	EmptyChar rune
	LeftCap  rune
	RightCap rune
}

// DefaultProgressBarStyle returns the default progress bar style
func DefaultProgressBarStyle() ProgressBarStyle {
	return ProgressBarStyle{
		Border:    engine.Style{FG: engine.ColorWhite},
		Fill:      engine.Style{FG: engine.ColorGreen},
		Empty:     engine.Style{FG: engine.ColorBrightBlack},
		FillChar:  '█',
		EmptyChar: '░',
		LeftCap:   '[',
		RightCap:  ']',
	}
}

// NewProgressBar creates a new progress bar
func NewProgressBar(x, y, width int) *ProgressBar {
	return &ProgressBar{
		X:        x,
		Y:        y,
		Width:    width,
		Progress: 0,
		Style:    DefaultProgressBarStyle(),
	}
}

// SetProgress sets the progress (0.0 to 1.0)
func (p *ProgressBar) SetProgress(progress float64) {
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}
	p.Progress = progress
}

// Render draws the progress bar
func (p *ProgressBar) Render(screen *engine.Screen) {
	innerWidth := p.Width - 2 // Subtract caps
	filled := int(float64(innerWidth) * p.Progress)

	screen.Set(p.X, p.Y, p.Style.LeftCap, p.Style.Border)

	for i := 0; i < innerWidth; i++ {
		if i < filled {
			screen.Set(p.X+1+i, p.Y, p.Style.FillChar, p.Style.Fill)
		} else {
			screen.Set(p.X+1+i, p.Y, p.Style.EmptyChar, p.Style.Empty)
		}
	}

	screen.Set(p.X+p.Width-1, p.Y, p.Style.RightCap, p.Style.Border)
}

// Spinner is an animated loading spinner
type Spinner struct {
	X, Y   int
	Frame  int
	Frames []rune
	Style  engine.Style
}

// DefaultSpinnerFrames returns the default spinner frames
func DefaultSpinnerFrames() []rune {
	return []rune{'⠋', '⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
}

// ASCIISpinnerFrames returns ASCII-only spinner frames
func ASCIISpinnerFrames() []rune {
	return []rune{'|', '/', '-', '\\'}
}

// NewSpinner creates a new spinner
func NewSpinner(x, y int) *Spinner {
	return &Spinner{
		X:      x,
		Y:      y,
		Frame:  0,
		Frames: DefaultSpinnerFrames(),
		Style:  engine.Style{FG: engine.ColorCyan},
	}
}

// Advance advances the spinner to the next frame
func (s *Spinner) Advance() {
	s.Frame = (s.Frame + 1) % len(s.Frames)
}

// Render draws the spinner
func (s *Spinner) Render(screen *engine.Screen) {
	screen.Set(s.X, s.Y, s.Frames[s.Frame], s.Style)
}

// Logo renders the SLIP logo (simple ASCII - works everywhere)
func DrawLogo(screen *engine.Screen, x, y int, style engine.Style) {
	logo := []string{
		" ____  _     ___ ____  ",
		"/ ___|| |   |_ _|  _ \\ ",
		"\\___ \\| |    | || |_) |",
		" ___) | |___ | ||  __/ ",
		"|____/|_____|___|_|    ",
	}
	for i, line := range logo {
		screen.DrawText(x, y+i, line, style)
	}
}

// DrawLogoCompact renders a compact SLIP logo for smaller spaces
func DrawLogoCompact(screen *engine.Screen, x, y int, style engine.Style) {
	logo := []string{
		"╔═╗╦  ╦╔═╗",
		"╚═╗║  ║╠═╝",
		"╚═╝╩═╝╩╩  ",
	}
	for i, line := range logo {
		screen.DrawText(x, y+i, line, style)
	}
}

// DrawLogoLarge renders a larger SLIP logo
func DrawLogoLarge(screen *engine.Screen, x, y int, style engine.Style) {
	logo := []string{
		"  ██████  ██▓     ██▓ ██▓███  ",
		"▒██    ▒ ▓██▒    ▓██▒▓██░  ██▒",
		"░ ▓██▄   ▒██░    ▒██▒▓██░ ██▓▒",
		"  ▒   ██▒▒██░    ░██░▒██▄█▓▒ ▒",
		"▒██████▒▒░██████▒░██░▒██▒ ░  ░",
		"▒ ▒▓▒ ▒ ░░ ▒░▓  ░░▓  ▒▓▒░ ░  ░",
		"░ ░▒  ░ ░░ ░ ▒  ░ ▒ ░░▒ ░     ",
		"░  ░  ░    ░ ░    ▒ ░░░       ",
		"      ░      ░  ░ ░           ",
	}
	for i, line := range logo {
		screen.DrawText(x, y+i, line, style)
	}
}
