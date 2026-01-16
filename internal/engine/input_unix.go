//go:build !windows

package engine

import (
	"os"

	"golang.org/x/term"
)

func (r *InputReader) setRawModeUnix() error {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	r.oldState = oldState
	return nil
}

func (r *InputReader) restoreModeUnix() {
	if r.oldState != nil {
		fd := int(os.Stdin.Fd())
		term.Restore(fd, r.oldState.(*term.State))
	}
}

func (r *InputReader) setRawModeWindows() error {
	return nil
}

func (r *InputReader) restoreModeWindows() {
}
