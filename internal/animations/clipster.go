package animations

import (
	"slip/internal/engine"
)

// ClipsterAnimation is the animated paperclip buddy
type ClipsterAnimation struct {
	ctx        *engine.GameContext
	frames     [][]string
	frameIndex int
	timer      float64
	frameDelay float64
}

// NewClipster creates a new Clipster animation
func NewClipster() *ClipsterAnimation {
	return &ClipsterAnimation{
		frameDelay: 0.3,
		frames:     clipsterFrames(),
	}
}

// Info returns animation metadata
func (c *ClipsterAnimation) Info() engine.GameInfo {
	return engine.GameInfo{
		ID:          "clipster",
		Name:        "Clipster",
		Description: "Your friendly terminal paperclip companion",
		Author:      "Slip",
		Version:     "1.0.0",
		MinWidth:    16,
		MinHeight:   12,
		Category:    "animation",
	}
}

// Init initializes the animation
func (c *ClipsterAnimation) Init(ctx *engine.GameContext) error {
	c.ctx = ctx
	c.frameIndex = 0
	c.timer = 0
	return nil
}

// Update updates the animation
func (c *ClipsterAnimation) Update(dt float64) error {
	c.timer += dt
	if c.timer >= c.frameDelay {
		c.timer = 0
		c.frameIndex = (c.frameIndex + 1) % len(c.frames)
	}
	return nil
}

// Render draws the animation
func (c *ClipsterAnimation) Render(screen *engine.Screen) error {
	frame := c.frames[c.frameIndex]

	// Center the clipster
	maxWidth := 0
	for _, line := range frame {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	x := (c.ctx.Width - maxWidth) / 2
	y := (c.ctx.Height - len(frame)) / 2

	// Draw the clipster
	style := engine.Style{FG: engine.ColorBrightCyan}
	eyeStyle := engine.Style{FG: engine.ColorBrightWhite, Bold: true}

	for dy, line := range frame {
		for dx, char := range line {
			if char == ' ' {
				continue
			}
			s := style
			if char == 'o' || char == '-' || char == '^' {
				s = eyeStyle
			}
			screen.Set(x+dx, y+dy, char, s)
		}
	}

	// Draw title
	screen.DrawTextCentered(1, "Clipster", engine.Style{FG: engine.ColorBrightCyan, Bold: true})
	screen.DrawTextCentered(c.ctx.Height-2, "Your terminal buddy", engine.Style{FG: engine.ColorBrightBlack})

	return nil
}

func clipsterFrames() [][]string {
	return [][]string{
		// Frame 1: Normal
		{
			"    .---.    ",
			"   /     \\   ",
			"  | o   o |  ",
			"  |   =   |  ",
			"   \\_____/   ",
			"      |      ",
			"     /|\\     ",
			"    / | \\    ",
			"      |      ",
			"     / \\     ",
		},
		// Frame 2: Blink
		{
			"    .---.    ",
			"   /     \\   ",
			"  | -   - |  ",
			"  |   =   |  ",
			"   \\_____/   ",
			"      |      ",
			"     /|\\     ",
			"    / | \\    ",
			"      |      ",
			"     / \\     ",
		},
		// Frame 3: Normal
		{
			"    .---.    ",
			"   /     \\   ",
			"  | o   o |  ",
			"  |   =   |  ",
			"   \\_____/   ",
			"      |      ",
			"     /|\\     ",
			"    / | \\    ",
			"      |      ",
			"     / \\     ",
		},
		// Frame 4: Look left
		{
			"    .---.    ",
			"   /     \\   ",
			"  |o    o |  ",
			"  |   =   |  ",
			"   \\_____/   ",
			"      |      ",
			"     /|\\     ",
			"    / | \\    ",
			"      |      ",
			"     / \\     ",
		},
		// Frame 5: Normal
		{
			"    .---.    ",
			"   /     \\   ",
			"  | o   o |  ",
			"  |   =   |  ",
			"   \\_____/   ",
			"      |      ",
			"     /|\\     ",
			"    / | \\    ",
			"      |      ",
			"     / \\     ",
		},
		// Frame 6: Look right
		{
			"    .---.    ",
			"   /     \\   ",
			"  | o    o|  ",
			"  |   =   |  ",
			"   \\_____/   ",
			"      |      ",
			"     /|\\     ",
			"    / | \\    ",
			"      |      ",
			"     / \\     ",
		},
		// Frame 7: Wave
		{
			"    .---.    ",
			"   /     \\   ",
			"  | ^   ^ |  ",
			"  |   =   |  ",
			"   \\_____/ \\ ",
			"      |     \\",
			"     /|      ",
			"    / |      ",
			"      |      ",
			"     / \\     ",
		},
		// Frame 8: Wave 2
		{
			"    .---.    ",
			"   /     \\   ",
			"  | ^   ^ |  ",
			"  |   =   |  ",
			"   \\_____/  /",
			"      |    / ",
			"     /|      ",
			"    / |      ",
			"      |      ",
			"     / \\     ",
		},
	}
}
