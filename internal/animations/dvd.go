package animations

import (
	"math/rand"
	"time"

	"slip/internal/engine"
)

// DVDAnimation is the classic bouncing DVD logo
type DVDAnimation struct {
	ctx        *engine.GameContext
	x, y       float64
	vx, vy     float64
	color      engine.Color
	colorIndex int
	cornerHits int
	rng        *rand.Rand
}

// NewDVD creates a new DVD bounce animation
func NewDVD() *DVDAnimation {
	return &DVDAnimation{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Info returns animation metadata
func (d *DVDAnimation) Info() engine.GameInfo {
	return engine.GameInfo{
		ID:          "dvd",
		Name:        "DVD Bounce",
		Description: "The classic bouncing DVD logo - wait for the corner!",
		Author:      "Slip",
		Version:     "1.0.0",
		MinWidth:    20,
		MinHeight:   10,
		Category:    "animation",
	}
}

var dvdColors = []engine.Color{
	engine.ColorRed,
	engine.ColorGreen,
	engine.ColorBlue,
	engine.ColorYellow,
	engine.ColorMagenta,
	engine.ColorCyan,
}

var dvdLogo = []string{
	"╔═══════╗",
	"║  DVD  ║",
	"╚═══════╝",
}

// Init initializes the animation
func (d *DVDAnimation) Init(ctx *engine.GameContext) error {
	d.ctx = ctx

	// Start in a random position
	logoW := len(dvdLogo[0])
	logoH := len(dvdLogo)

	d.x = float64(d.rng.Intn(ctx.Width - logoW))
	d.y = float64(d.rng.Intn(ctx.Height - logoH))

	// Random velocity
	d.vx = 10 + d.rng.Float64()*5
	d.vy = 6 + d.rng.Float64()*3
	if d.rng.Intn(2) == 0 {
		d.vx = -d.vx
	}
	if d.rng.Intn(2) == 0 {
		d.vy = -d.vy
	}

	d.colorIndex = d.rng.Intn(len(dvdColors))
	d.color = dvdColors[d.colorIndex]
	d.cornerHits = 0

	return nil
}

// Update updates the animation
func (d *DVDAnimation) Update(dt float64) error {
	logoW := float64(len(dvdLogo[0]))
	logoH := float64(len(dvdLogo))

	d.x += d.vx * dt
	d.y += d.vy * dt

	hitX := false
	hitY := false

	// Bounce off walls
	if d.x <= 0 {
		d.x = 0
		d.vx = -d.vx
		hitX = true
	}
	if d.x >= float64(d.ctx.Width)-logoW {
		d.x = float64(d.ctx.Width) - logoW
		d.vx = -d.vx
		hitX = true
	}
	if d.y <= 0 {
		d.y = 0
		d.vy = -d.vy
		hitY = true
	}
	if d.y >= float64(d.ctx.Height)-logoH {
		d.y = float64(d.ctx.Height) - logoH
		d.vy = -d.vy
		hitY = true
	}

	// Change color on any wall hit
	if hitX || hitY {
		d.colorIndex = (d.colorIndex + 1) % len(dvdColors)
		d.color = dvdColors[d.colorIndex]
	}

	// Track corner hits
	if hitX && hitY {
		d.cornerHits++
	}

	return nil
}

// Render draws the animation
func (d *DVDAnimation) Render(screen *engine.Screen) error {
	// Draw the logo
	style := engine.Style{FG: d.color, Bold: true}

	for dy, line := range dvdLogo {
		for dx, char := range line {
			screen.Set(int(d.x)+dx, int(d.y)+dy, char, style)
		}
	}

	// Draw corner hit counter if any
	if d.cornerHits > 0 {
		text := "Corner hits: "
		screen.DrawText(2, d.ctx.Height-1, text, engine.Style{FG: engine.ColorBrightBlack})
		screen.DrawText(2+len(text), d.ctx.Height-1, string(rune('0'+d.cornerHits%10)),
			engine.Style{FG: engine.ColorBrightGreen, Bold: true})
	}

	return nil
}
