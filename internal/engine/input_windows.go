//go:build windows

package engine

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	kernel32           = syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode     = kernel32.NewProc("GetConsoleMode")
	setConsoleMode     = kernel32.NewProc("SetConsoleMode")
	getStdHandle       = kernel32.NewProc("GetStdHandle")
)

const (
	stdInputHandle     = ^uintptr(0) - 10 + 1 // STD_INPUT_HANDLE = -10
	enableLineInput    = 0x0002
	enableEchoInput    = 0x0004
	enableProcessedInput = 0x0001
	enableVirtualTerminalInput = 0x0200
)

type windowsState struct {
	mode uint32
}

func (r *InputReader) setRawModeWindows() error {
	handle, _, _ := getStdHandle.Call(stdInputHandle)
	if handle == 0 {
		// Fallback - try direct file descriptor
		return nil
	}

	var mode uint32
	ret, _, _ := getConsoleMode.Call(handle, uintptr(unsafe.Pointer(&mode)))
	if ret == 0 {
		return nil
	}

	r.oldState = &windowsState{mode: mode}

	// Disable line input, echo, and processed input
	newMode := mode &^ (enableLineInput | enableEchoInput | enableProcessedInput)
	// Enable virtual terminal input for escape sequences
	newMode |= enableVirtualTerminalInput

	setConsoleMode.Call(handle, uintptr(newMode))
	return nil
}

func (r *InputReader) restoreModeWindows() {
	if r.oldState == nil {
		return
	}
	state, ok := r.oldState.(*windowsState)
	if !ok {
		return
	}

	handle, _, _ := getStdHandle.Call(stdInputHandle)
	if handle == 0 {
		return
	}
	setConsoleMode.Call(handle, uintptr(state.mode))
}

func (r *InputReader) setRawModeUnix() error {
	return nil
}

func (r *InputReader) restoreModeUnix() {
}

// Ensure os is used
var _ = os.Stdin
