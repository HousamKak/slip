package engine

import (
	"fmt"
	"io"
	"strings"

	"slip/internal/term"
)

// Screen represents a terminal screen buffer
type Screen struct {
	width    int
	height   int
	buffer   [][]Cell
	previous [][]Cell
	dirty    bool
	writer   io.Writer
}

// NewScreen creates a new screen with the given dimensions
func NewScreen(width, height int, writer io.Writer) *Screen {
	s := &Screen{
		width:  width,
		height: height,
		writer: writer,
	}
	s.buffer = s.makeBuffer()
	s.previous = s.makeBuffer()
	return s
}

func (s *Screen) makeBuffer() [][]Cell {
	buffer := make([][]Cell, s.height)
	for i := range buffer {
		buffer[i] = make([]Cell, s.width)
		for j := range buffer[i] {
			buffer[i][j] = EmptyCell()
		}
	}
	return buffer
}

// Width returns the screen width
func (s *Screen) Width() int {
	return s.width
}

// Height returns the screen height
func (s *Screen) Height() int {
	return s.height
}

// Resize resizes the screen
func (s *Screen) Resize(width, height int) {
	s.width = width
	s.height = height
	s.buffer = s.makeBuffer()
	s.previous = s.makeBuffer()
	s.dirty = true
}

// Clear clears the screen buffer
func (s *Screen) Clear() {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			s.buffer[y][x] = EmptyCell()
		}
	}
	s.dirty = true
}

// Fill fills the screen with a character
func (s *Screen) Fill(r rune, style Style) {
	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			s.buffer[y][x] = Cell{Rune: r, Style: style}
		}
	}
	s.dirty = true
}

// Set sets a cell at the given position
func (s *Screen) Set(x, y int, r rune, style Style) {
	if x >= 0 && x < s.width && y >= 0 && y < s.height {
		s.buffer[y][x] = Cell{Rune: r, Style: style}
		s.dirty = true
	}
}

// SetCell sets a cell at the given position
func (s *Screen) SetCell(x, y int, cell Cell) {
	if x >= 0 && x < s.width && y >= 0 && y < s.height {
		s.buffer[y][x] = cell
		s.dirty = true
	}
}

// Get gets a cell at the given position
func (s *Screen) Get(x, y int) Cell {
	if x >= 0 && x < s.width && y >= 0 && y < s.height {
		return s.buffer[y][x]
	}
	return EmptyCell()
}

// DrawText draws text at the given position
func (s *Screen) DrawText(x, y int, text string, style Style) {
	for i, r := range text {
		s.Set(x+i, y, r, style)
	}
}

// DrawTextCentered draws text centered at the given y position
func (s *Screen) DrawTextCentered(y int, text string, style Style) {
	x := (s.width - len(text)) / 2
	s.DrawText(x, y, text, style)
}

// DrawTextRight draws text right-aligned
func (s *Screen) DrawTextRight(x, y int, text string, style Style) {
	s.DrawText(x-len(text), y, text, style)
}

// DrawBox draws a box
func (s *Screen) DrawBox(x, y, w, h int, boxStyle BoxStyle, style Style) {
	if w < 2 || h < 2 {
		return
	}
	chars := boxStyle.BoxChars()

	// Corners
	s.Set(x, y, chars.TopLeft, style)
	s.Set(x+w-1, y, chars.TopRight, style)
	s.Set(x, y+h-1, chars.BottomLeft, style)
	s.Set(x+w-1, y+h-1, chars.BottomRight, style)

	// Horizontal lines
	for i := 1; i < w-1; i++ {
		s.Set(x+i, y, chars.Horizontal, style)
		s.Set(x+i, y+h-1, chars.Horizontal, style)
	}

	// Vertical lines
	for i := 1; i < h-1; i++ {
		s.Set(x, y+i, chars.Vertical, style)
		s.Set(x+w-1, y+i, chars.Vertical, style)
	}
}

// DrawFilledBox draws a filled box
func (s *Screen) DrawFilledBox(x, y, w, h int, boxStyle BoxStyle, borderStyle Style, fillRune rune, fillStyle Style) {
	// Fill interior
	for dy := 1; dy < h-1; dy++ {
		for dx := 1; dx < w-1; dx++ {
			s.Set(x+dx, y+dy, fillRune, fillStyle)
		}
	}
	// Draw border
	s.DrawBox(x, y, w, h, boxStyle, borderStyle)
}

// DrawHLine draws a horizontal line
func (s *Screen) DrawHLine(x, y, length int, r rune, style Style) {
	for i := 0; i < length; i++ {
		s.Set(x+i, y, r, style)
	}
}

// DrawVLine draws a vertical line
func (s *Screen) DrawVLine(x, y, length int, r rune, style Style) {
	for i := 0; i < length; i++ {
		s.Set(x, y+i, r, style)
	}
}

// DrawRect draws a rectangle (outline only)
func (s *Screen) DrawRect(r Rect, char rune, style Style) {
	s.DrawHLine(r.X, r.Y, r.Width, char, style)
	s.DrawHLine(r.X, r.Y+r.Height-1, r.Width, char, style)
	s.DrawVLine(r.X, r.Y, r.Height, char, style)
	s.DrawVLine(r.X+r.Width-1, r.Y, r.Height, char, style)
}

// FillRect fills a rectangle
func (s *Screen) FillRect(r Rect, char rune, style Style) {
	for y := r.Y; y < r.Y+r.Height; y++ {
		for x := r.X; x < r.X+r.Width; x++ {
			s.Set(x, y, char, style)
		}
	}
}

// DrawSprite draws multi-line ASCII art at position
func (s *Screen) DrawSprite(x, y int, lines []string, style Style) {
	for dy, line := range lines {
		for dx, r := range line {
			if r != ' ' {
				s.Set(x+dx, y+dy, r, style)
			}
		}
	}
}

// DrawSpriteColored draws multi-line ASCII art with per-line styles
func (s *Screen) DrawSpriteColored(x, y int, lines []string, styles []Style) {
	for dy, line := range lines {
		style := DefaultStyle()
		if dy < len(styles) {
			style = styles[dy]
		}
		for dx, r := range line {
			if r != ' ' {
				s.Set(x+dx, y+dy, r, style)
			}
		}
	}
}

// Flush renders the buffer to the terminal
func (s *Screen) Flush() {
	if !s.dirty {
		return
	}

	var sb strings.Builder
	sb.WriteString(term.CursorHome)

	lastStyle := DefaultStyle()
	needStyleReset := false

	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			cell := s.buffer[y][x]

			// Check if style changed
			if cell.Style != lastStyle || needStyleReset {
				sb.WriteString(term.Reset)
				sb.WriteString(s.styleToANSI(cell.Style))
				lastStyle = cell.Style
				needStyleReset = false
			}

			sb.WriteRune(cell.Rune)
		}
		if y < s.height-1 {
			sb.WriteString("\n")
			needStyleReset = true
		}
	}

	sb.WriteString(term.Reset)
	fmt.Fprint(s.writer, sb.String())

	// Swap buffers
	s.buffer, s.previous = s.previous, s.buffer
	s.dirty = false
}

// FlushDiff renders only changed cells (more efficient but complex)
func (s *Screen) FlushDiff() {
	var sb strings.Builder
	lastStyle := DefaultStyle()
	lastX, lastY := -1, -1

	for y := 0; y < s.height; y++ {
		for x := 0; x < s.width; x++ {
			cell := s.buffer[y][x]
			prev := s.previous[y][x]

			if cell.Rune == prev.Rune && cell.Style == prev.Style {
				continue
			}

			// Move cursor if needed
			if lastY != y || lastX != x-1 {
				sb.WriteString(fmt.Sprintf("\x1b[%d;%dH", y+1, x+1))
			}

			// Update style if needed
			if cell.Style != lastStyle {
				sb.WriteString(term.Reset)
				sb.WriteString(s.styleToANSI(cell.Style))
				lastStyle = cell.Style
			}

			sb.WriteRune(cell.Rune)
			lastX = x
			lastY = y
		}
	}

	if sb.Len() > 0 {
		sb.WriteString(term.Reset)
		fmt.Fprint(s.writer, sb.String())
	}

	// Copy current to previous
	for y := 0; y < s.height; y++ {
		copy(s.previous[y], s.buffer[y])
	}
	s.dirty = false
}

func (s *Screen) styleToANSI(style Style) string {
	var sb strings.Builder

	if style.Bold {
		sb.WriteString(term.Bold)
	}
	if style.Dim {
		sb.WriteString(term.Dim)
	}
	if style.Italic {
		sb.WriteString(term.Italic)
	}
	if style.Underline {
		sb.WriteString(term.Underline)
	}
	if style.Blink {
		sb.WriteString(term.Blink)
	}
	if style.Reverse {
		sb.WriteString(term.Reverse)
	}

	sb.WriteString(s.colorToFgANSI(style.FG))
	sb.WriteString(s.colorToBgANSI(style.BG))

	return sb.String()
}

func (s *Screen) colorToFgANSI(c Color) string {
	switch c {
	case ColorBlack:
		return term.FgBlack
	case ColorRed:
		return term.FgRed
	case ColorGreen:
		return term.FgGreen
	case ColorYellow:
		return term.FgYellow
	case ColorBlue:
		return term.FgBlue
	case ColorMagenta:
		return term.FgMagenta
	case ColorCyan:
		return term.FgCyan
	case ColorWhite:
		return term.FgWhite
	case ColorBrightBlack:
		return term.FgBrightBlack
	case ColorBrightRed:
		return term.FgBrightRed
	case ColorBrightGreen:
		return term.FgBrightGreen
	case ColorBrightYellow:
		return term.FgBrightYellow
	case ColorBrightBlue:
		return term.FgBrightBlue
	case ColorBrightMagenta:
		return term.FgBrightMagenta
	case ColorBrightCyan:
		return term.FgBrightCyan
	case ColorBrightWhite:
		return term.FgBrightWhite
	default:
		return term.FgDefault
	}
}

func (s *Screen) colorToBgANSI(c Color) string {
	switch c {
	case ColorBlack:
		return term.BgBlack
	case ColorRed:
		return term.BgRed
	case ColorGreen:
		return term.BgGreen
	case ColorYellow:
		return term.BgYellow
	case ColorBlue:
		return term.BgBlue
	case ColorMagenta:
		return term.BgMagenta
	case ColorCyan:
		return term.BgCyan
	case ColorWhite:
		return term.BgWhite
	default:
		return term.BgDefault
	}
}
