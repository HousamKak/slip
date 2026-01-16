package animations

import (
	"math/rand"
	"time"

	"slip/internal/engine"
)

// FireAnimation is a classic demoscene fire effect
type FireAnimation struct {
	ctx    *engine.GameContext
	buffer [][]int
	rng    *rand.Rand
}

// Fire palette characters (from coolest to hottest)
var fireChars = []rune{' ', '.', ':', '*', 's', 'S', '#', '$', '%', '@'}

// NewFire creates a new Fire animation
func NewFire() *FireAnimation {
	return &FireAnimation{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Info returns animation metadata
func (f *FireAnimation) Info() engine.GameInfo {
	return engine.GameInfo{
		ID:          "fire",
		Name:        "Fire Effect",
		Description: "Classic demoscene fire simulation",
		Author:      "Slip",
		Version:     "1.0.0",
		MinWidth:    20,
		MinHeight:   10,
		Category:    "animation",
	}
}

// Init initializes the animation
func (f *FireAnimation) Init(ctx *engine.GameContext) error {
	f.ctx = ctx

	// Create buffer (add extra row at bottom for heat source)
	f.buffer = make([][]int, ctx.Height+1)
	for i := range f.buffer {
		f.buffer[i] = make([]int, ctx.Width)
	}

	return nil
}

// Update updates the animation
func (f *FireAnimation) Update(dt float64) error {
	height := len(f.buffer)
	width := len(f.buffer[0])

	// Set bottom row on fire (random heat values)
	for x := 0; x < width; x++ {
		// More heat in the center
		centerDist := float64(x-width/2) / float64(width/2)
		if centerDist < 0 {
			centerDist = -centerDist
		}
		maxHeat := int(float64(len(fireChars)-1) * (1 - centerDist*0.5))

		if f.rng.Intn(4) == 0 {
			f.buffer[height-1][x] = f.rng.Intn(maxHeat + 1)
		} else {
			f.buffer[height-1][x] = maxHeat
		}
	}

	// Propagate fire upwards with cooling
	for y := 0; y < height-1; y++ {
		for x := 0; x < width; x++ {
			// Sample from below (with horizontal spread)
			sum := 0
			count := 0

			// Current cell
			if y+1 < height {
				sum += f.buffer[y+1][x]
				count++
			}

			// Below-left
			if y+1 < height && x > 0 {
				sum += f.buffer[y+1][x-1]
				count++
			}

			// Below-right
			if y+1 < height && x < width-1 {
				sum += f.buffer[y+1][x+1]
				count++
			}

			// Two rows below
			if y+2 < height {
				sum += f.buffer[y+2][x]
				count++
			}

			// Average and cool down
			if count > 0 {
				avg := sum / count
				// Random cooling
				cooling := f.rng.Intn(2)
				f.buffer[y][x] = avg - cooling
				if f.buffer[y][x] < 0 {
					f.buffer[y][x] = 0
				}
			}
		}
	}

	return nil
}

// Render draws the animation
func (f *FireAnimation) Render(screen *engine.Screen) error {
	for y := 0; y < f.ctx.Height; y++ {
		for x := 0; x < f.ctx.Width; x++ {
			heat := f.buffer[y][x]
			if heat < 0 {
				heat = 0
			}
			if heat >= len(fireChars) {
				heat = len(fireChars) - 1
			}

			char := fireChars[heat]
			if char == ' ' {
				continue
			}

			// Color based on heat
			var style engine.Style
			if heat < 3 {
				style = engine.Style{FG: engine.ColorRed}
			} else if heat < 5 {
				style = engine.Style{FG: engine.ColorBrightRed}
			} else if heat < 7 {
				style = engine.Style{FG: engine.ColorYellow}
			} else {
				style = engine.Style{FG: engine.ColorBrightYellow, Bold: true}
			}

			screen.Set(x, y, char, style)
		}
	}

	// Draw title at top
	screen.DrawTextCentered(0, "Fire Effect", engine.Style{FG: engine.ColorBrightWhite, Bold: true})

	return nil
}
