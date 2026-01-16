package animations

import (
	"math/rand"
	"time"

	"slip/internal/engine"
)

// Star represents a star in the starfield
type Star struct {
	X, Y  float64
	Speed float64
	Char  rune
}

// StarfieldAnimation is a parallax starfield effect
type StarfieldAnimation struct {
	ctx   *engine.GameContext
	stars []Star
	rng   *rand.Rand
}

// NewStarfield creates a new Starfield animation
func NewStarfield() *StarfieldAnimation {
	return &StarfieldAnimation{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Info returns animation metadata
func (s *StarfieldAnimation) Info() engine.GameInfo {
	return engine.GameInfo{
		ID:          "starfield",
		Name:        "Starfield",
		Description: "Relaxing parallax star effect",
		Author:      "Slip",
		Version:     "1.0.0",
		MinWidth:    20,
		MinHeight:   10,
		Category:    "animation",
	}
}

// Init initializes the animation
func (s *StarfieldAnimation) Init(ctx *engine.GameContext) error {
	s.ctx = ctx

	// Create stars
	numStars := ctx.Width * ctx.Height / 15
	if numStars < 20 {
		numStars = 20
	}
	if numStars > 200 {
		numStars = 200
	}

	s.stars = make([]Star, numStars)
	for i := range s.stars {
		s.stars[i] = s.newStar(true)
	}

	return nil
}

func (s *StarfieldAnimation) newStar(randomX bool) Star {
	// Different star types with different speeds and appearances
	layer := s.rng.Intn(3)

	var char rune
	var speed float64

	switch layer {
	case 0: // Far/slow
		char = '.'
		speed = 3 + s.rng.Float64()*2
	case 1: // Medium
		char = '·'
		speed = 7 + s.rng.Float64()*3
	case 2: // Near/fast
		chars := []rune{'*', '✦', '+'}
		char = chars[s.rng.Intn(len(chars))]
		speed = 12 + s.rng.Float64()*5
	}

	x := float64(s.ctx.Width)
	if randomX {
		x = s.rng.Float64() * float64(s.ctx.Width)
	}

	return Star{
		X:     x,
		Y:     s.rng.Float64() * float64(s.ctx.Height),
		Speed: speed,
		Char:  char,
	}
}

// Update updates the animation
func (s *StarfieldAnimation) Update(dt float64) error {
	for i := range s.stars {
		s.stars[i].X -= s.stars[i].Speed * dt

		// Respawn star on right side when it goes off left
		if s.stars[i].X < 0 {
			s.stars[i] = s.newStar(false)
		}
	}
	return nil
}

// Render draws the animation
func (s *StarfieldAnimation) Render(screen *engine.Screen) error {
	// Draw stars
	for _, star := range s.stars {
		x := int(star.X)
		y := int(star.Y)
		if x >= 0 && x < s.ctx.Width && y >= 0 && y < s.ctx.Height {
			// Color based on speed (depth)
			var style engine.Style
			if star.Speed < 5 {
				style = engine.Style{FG: engine.ColorBrightBlack}
			} else if star.Speed < 10 {
				style = engine.Style{FG: engine.ColorWhite}
			} else {
				style = engine.Style{FG: engine.ColorBrightWhite, Bold: true}
			}
			screen.Set(x, y, star.Char, style)
		}
	}

	// Draw title
	screen.DrawTextCentered(0, "Starfield", engine.Style{FG: engine.ColorBrightBlue, Bold: true})

	return nil
}
