package animations

import (
	"math/rand"
	"time"

	"slip/internal/engine"
)

// Column represents a falling column of characters
type Column struct {
	X       int
	Y       float64
	Speed   float64
	Length  int
	Chars   []rune
	Active  bool
	Delay   float64 // Delay before starting
}

// MatrixAnimation is the Matrix-style digital rain effect
type MatrixAnimation struct {
	ctx     *engine.GameContext
	columns []Column
	rng     *rand.Rand
	chars   []rune
}

// NewMatrix creates a new Matrix animation
func NewMatrix() *MatrixAnimation {
	// Matrix-style characters (katakana + numbers + symbols)
	chars := []rune{
		'ア', 'イ', 'ウ', 'エ', 'オ', 'カ', 'キ', 'ク', 'ケ', 'コ',
		'サ', 'シ', 'ス', 'セ', 'ソ', 'タ', 'チ', 'ツ', 'テ', 'ト',
		'ナ', 'ニ', 'ヌ', 'ネ', 'ノ', 'ハ', 'ヒ', 'フ', 'ヘ', 'ホ',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		'$', '+', '-', '*', '%', ':', '.', '=', '<', '>',
	}

	return &MatrixAnimation{
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
		chars: chars,
	}
}

// Info returns animation metadata
func (m *MatrixAnimation) Info() engine.GameInfo {
	return engine.GameInfo{
		ID:          "matrix",
		Name:        "Matrix Rain",
		Description: "Digital rain from The Matrix",
		Author:      "Slip",
		Version:     "1.0.0",
		MinWidth:    20,
		MinHeight:   10,
		Category:    "animation",
	}
}

// Init initializes the animation
func (m *MatrixAnimation) Init(ctx *engine.GameContext) error {
	m.ctx = ctx

	// Create a column for each x position (spaced out)
	m.columns = make([]Column, ctx.Width/2)
	for i := range m.columns {
		m.columns[i] = m.newColumn(i * 2)
		// Stagger start times
		m.columns[i].Delay = m.rng.Float64() * 2
		m.columns[i].Active = false
	}

	return nil
}

func (m *MatrixAnimation) newColumn(x int) Column {
	length := 5 + m.rng.Intn(15)
	if length > m.ctx.Height {
		length = m.ctx.Height
	}

	chars := make([]rune, length)
	for i := range chars {
		chars[i] = m.randomChar()
	}

	return Column{
		X:      x,
		Y:      0,
		Speed:  8 + m.rng.Float64()*10,
		Length: length,
		Chars:  chars,
		Active: true,
		Delay:  0,
	}
}

func (m *MatrixAnimation) randomChar() rune {
	return m.chars[m.rng.Intn(len(m.chars))]
}

// Update updates the animation
func (m *MatrixAnimation) Update(dt float64) error {
	for i := range m.columns {
		col := &m.columns[i]

		// Handle delay
		if col.Delay > 0 {
			col.Delay -= dt
			continue
		}
		col.Active = true

		// Move column down
		col.Y += col.Speed * dt

		// Randomly change characters
		if m.rng.Float64() < 0.1 {
			changeIdx := m.rng.Intn(len(col.Chars))
			col.Chars[changeIdx] = m.randomChar()
		}

		// Reset column when it goes off screen
		if int(col.Y)-col.Length > m.ctx.Height {
			newCol := m.newColumn(col.X)
			newCol.Delay = m.rng.Float64() * 1.5
			newCol.Active = false
			m.columns[i] = newCol
		}
	}
	return nil
}

// Render draws the animation
func (m *MatrixAnimation) Render(screen *engine.Screen) error {
	for _, col := range m.columns {
		if !col.Active {
			continue
		}

		headY := int(col.Y)

		for i, char := range col.Chars {
			y := headY - i
			if y < 0 || y >= m.ctx.Height {
				continue
			}

			var style engine.Style
			if i == 0 {
				// Head is bright white
				style = engine.Style{FG: engine.ColorBrightWhite, Bold: true}
			} else if i < 3 {
				// Near head is bright green
				style = engine.Style{FG: engine.ColorBrightGreen}
			} else if i < col.Length/2 {
				// Middle is regular green
				style = engine.Style{FG: engine.ColorGreen}
			} else {
				// Tail fades to dark
				style = engine.Style{FG: engine.ColorBrightBlack}
			}

			screen.Set(col.X, y, char, style)
		}
	}

	return nil
}
