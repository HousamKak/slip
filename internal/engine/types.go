package engine

// Color represents a terminal color
type Color int

const (
	ColorDefault Color = iota
	ColorBlack
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBrightBlack
	ColorBrightRed
	ColorBrightGreen
	ColorBrightYellow
	ColorBrightBlue
	ColorBrightMagenta
	ColorBrightCyan
	ColorBrightWhite
)

// Style represents text styling
type Style struct {
	FG        Color
	BG        Color
	Bold      bool
	Dim       bool
	Italic    bool
	Underline bool
	Blink     bool
	Reverse   bool
}

// DefaultStyle returns the default style
func DefaultStyle() Style {
	return Style{FG: ColorDefault, BG: ColorDefault}
}

// WithFG returns a copy with foreground color set
func (s Style) WithFG(c Color) Style {
	s.FG = c
	return s
}

// WithBG returns a copy with background color set
func (s Style) WithBG(c Color) Style {
	s.BG = c
	return s
}

// WithBold returns a copy with bold set
func (s Style) WithBold() Style {
	s.Bold = true
	return s
}

// WithDim returns a copy with dim set
func (s Style) WithDim() Style {
	s.Dim = true
	return s
}

// Cell represents a single character cell
type Cell struct {
	Rune  rune
	Style Style
}

// EmptyCell returns an empty cell
func EmptyCell() Cell {
	return Cell{Rune: ' ', Style: DefaultStyle()}
}

// Point represents a 2D point
type Point struct {
	X, Y int
}

// Add adds two points
func (p Point) Add(other Point) Point {
	return Point{X: p.X + other.X, Y: p.Y + other.Y}
}

// Sub subtracts two points
func (p Point) Sub(other Point) Point {
	return Point{X: p.X - other.X, Y: p.Y - other.Y}
}

// Equals checks if two points are equal
func (p Point) Equals(other Point) bool {
	return p.X == other.X && p.Y == other.Y
}

// Rect represents a rectangle
type Rect struct {
	X, Y, Width, Height int
}

// Contains checks if a point is inside the rectangle
func (r Rect) Contains(p Point) bool {
	return p.X >= r.X && p.X < r.X+r.Width && p.Y >= r.Y && p.Y < r.Y+r.Height
}

// Intersects checks if two rectangles intersect
func (r Rect) Intersects(other Rect) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

// BoxStyle defines the characters used for box drawing
type BoxStyle int

const (
	BoxStyleNone BoxStyle = iota
	BoxStyleSingle
	BoxStyleDouble
	BoxStyleRounded
	BoxStyleBold
	BoxStyleASCII
)

// BoxChars returns the box drawing characters for a style
func (bs BoxStyle) BoxChars() BoxChars {
	switch bs {
	case BoxStyleSingle:
		return BoxChars{'┌', '┐', '└', '┘', '─', '│', '├', '┤', '┬', '┴', '┼'}
	case BoxStyleDouble:
		return BoxChars{'╔', '╗', '╚', '╝', '═', '║', '╠', '╣', '╦', '╩', '╬'}
	case BoxStyleRounded:
		return BoxChars{'╭', '╮', '╰', '╯', '─', '│', '├', '┤', '┬', '┴', '┼'}
	case BoxStyleBold:
		return BoxChars{'┏', '┓', '┗', '┛', '━', '┃', '┣', '┫', '┳', '┻', '╋'}
	case BoxStyleASCII:
		return BoxChars{'+', '+', '+', '+', '-', '|', '+', '+', '+', '+', '+'}
	default:
		return BoxChars{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
	}
}

// BoxChars holds box drawing characters
type BoxChars struct {
	TopLeft, TopRight, BottomLeft, BottomRight rune
	Horizontal, Vertical                       rune
	LeftT, RightT, TopT, BottomT, Cross        rune
}

// Direction represents a cardinal direction
type Direction int

const (
	DirNone Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
)

// Opposite returns the opposite direction
func (d Direction) Opposite() Direction {
	switch d {
	case DirUp:
		return DirDown
	case DirDown:
		return DirUp
	case DirLeft:
		return DirRight
	case DirRight:
		return DirLeft
	default:
		return DirNone
	}
}

// Delta returns the point delta for this direction
func (d Direction) Delta() Point {
	switch d {
	case DirUp:
		return Point{0, -1}
	case DirDown:
		return Point{0, 1}
	case DirLeft:
		return Point{-1, 0}
	case DirRight:
		return Point{1, 0}
	default:
		return Point{0, 0}
	}
}
