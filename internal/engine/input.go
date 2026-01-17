package engine

import (
	"os"
	"runtime"
	"sync"
	"time"
)

// Key represents a keyboard key
type Key int

const (
	KeyNone Key = iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyEnter
	KeyEscape
	KeySpace
	KeyTab
	KeyBackspace
	KeyDelete
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyCtrlC
	KeyCtrlD
	KeyCtrlQ
	KeyCtrlZ
	KeyRune // For regular characters
)

// Input represents an input event
type Input struct {
	Type InputType
	Key  Key
	Rune rune
}

// InputType represents the type of input event
type InputType int

const (
	InputKeyPress InputType = iota
	InputResize
	InputQuit
)

// IsQuit returns true if this input should quit the application
func (i Input) IsQuit() bool {
	return i.Type == InputQuit ||
		i.Key == KeyCtrlC ||
		i.Key == KeyCtrlD ||
		i.Key == KeyCtrlQ
}

// IsEscape returns true if this is an escape key
func (i Input) IsEscape() bool {
	return i.Key == KeyEscape
}

// IsEnter returns true if this is an enter key
func (i Input) IsEnter() bool {
	return i.Key == KeyEnter
}

// IsArrow returns true if this is an arrow key
func (i Input) IsArrow() bool {
	return i.Key == KeyUp || i.Key == KeyDown || i.Key == KeyLeft || i.Key == KeyRight
}

// Direction returns the direction for arrow keys
func (i Input) Direction() Direction {
	switch i.Key {
	case KeyUp:
		return DirUp
	case KeyDown:
		return DirDown
	case KeyLeft:
		return DirLeft
	case KeyRight:
		return DirRight
	default:
		return DirNone
	}
}

// InputReader reads input from the terminal
type InputReader struct {
	ch       chan Input
	done     chan struct{}
	oldState interface{} // Platform-specific terminal state
	wg       sync.WaitGroup
}

// NewInputReader creates a new input reader
func NewInputReader() *InputReader {
	return &InputReader{
		ch:   make(chan Input, 32),
		done: make(chan struct{}),
	}
}

// Start starts reading input
func (r *InputReader) Start() error {
	if err := r.setRawMode(); err != nil {
		return err
	}

	r.wg.Add(1)
	go r.readLoop()
	return nil
}

// Stop stops reading input
func (r *InputReader) Stop() {
	close(r.done)
	r.wg.Wait() // Wait for readLoop to finish
	r.restoreMode()
}

// Channel returns the input channel
func (r *InputReader) Channel() <-chan Input {
	return r.ch
}

// ReadWithTimeout reads input with a timeout
func (r *InputReader) ReadWithTimeout(timeout time.Duration) (Input, bool) {
	select {
	case input := <-r.ch:
		return input, true
	case <-time.After(timeout):
		return Input{}, false
	}
}

func (r *InputReader) readLoop() {
	defer r.wg.Done()

	// Lock this goroutine to an OS thread for better I/O responsiveness
	// This ensures the goroutine wakes up immediately when stdin has data
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	buf := make([]byte, 16)
	for {
		// Check if we should stop before blocking on read
		select {
		case <-r.done:
			return
		default:
		}

		// Blocking read from stdin
		// By locking to an OS thread, the scheduler will wake this up faster
		n, err := os.Stdin.Read(buf)
		if err != nil {
			continue
		}
		if n == 0 {
			continue
		}

		// Parse and send inputs immediately
		inputs := r.parseInput(buf[:n])
		for _, input := range inputs {
			select {
			case r.ch <- input:
			case <-r.done:
				return
			}
		}
	}
}

func (r *InputReader) parseInput(buf []byte) []Input {
	var inputs []Input

	for i := 0; i < len(buf); {
		// Check for escape sequences
		if buf[i] == 0x1b && i+1 < len(buf) {
			// ESC sequence
			if buf[i+1] == '[' && i+2 < len(buf) {
				// CSI sequence
				switch buf[i+2] {
				case 'A':
					inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyUp})
					i += 3
					continue
				case 'B':
					inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyDown})
					i += 3
					continue
				case 'C':
					inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyRight})
					i += 3
					continue
				case 'D':
					inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyLeft})
					i += 3
					continue
				case 'H':
					inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyHome})
					i += 3
					continue
				case 'F':
					inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyEnd})
					i += 3
					continue
				case '3':
					if i+3 < len(buf) && buf[i+3] == '~' {
						inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyDelete})
						i += 4
						continue
					}
				case '5':
					if i+3 < len(buf) && buf[i+3] == '~' {
						inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyPageUp})
						i += 4
						continue
					}
				case '6':
					if i+3 < len(buf) && buf[i+3] == '~' {
						inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyPageDown})
						i += 4
						continue
					}
				}
				// F1-F4 may be ESC O P/Q/R/S
			} else if buf[i+1] == 'O' && i+2 < len(buf) {
				switch buf[i+2] {
				case 'P':
					inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyF1})
					i += 3
					continue
				case 'Q':
					inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyF2})
					i += 3
					continue
				case 'R':
					inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyF3})
					i += 3
					continue
				case 'S':
					inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyF4})
					i += 3
					continue
				}
			}
			// Plain escape
			inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyEscape})
			i++
			continue
		}

		// Control characters
		switch buf[i] {
		case 0x03: // Ctrl+C
			inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyCtrlC})
		case 0x04: // Ctrl+D
			inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyCtrlD})
		case 0x09: // Tab
			inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyTab})
		case 0x0d, 0x0a: // Enter
			inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyEnter})
		case 0x11: // Ctrl+Q
			inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyCtrlQ})
		case 0x1a: // Ctrl+Z
			inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyCtrlZ})
		case 0x1b: // Escape (standalone)
			inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyEscape})
		case 0x7f: // Backspace
			inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyBackspace})
		case ' ':
			inputs = append(inputs, Input{Type: InputKeyPress, Key: KeySpace})
		default:
			if buf[i] >= 0x20 && buf[i] < 0x7f {
				inputs = append(inputs, Input{Type: InputKeyPress, Key: KeyRune, Rune: rune(buf[i])})
			}
		}
		i++
	}

	return inputs
}

// Platform-specific raw mode handling
func (r *InputReader) setRawMode() error {
	if runtime.GOOS == "windows" {
		return r.setRawModeWindows()
	}
	return r.setRawModeUnix()
}

func (r *InputReader) restoreMode() {
	if runtime.GOOS == "windows" {
		r.restoreModeWindows()
	} else {
		r.restoreModeUnix()
	}
}
