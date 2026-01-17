package tmux

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const slipPaneOption = "@slip_pane"

// InTmux returns true if running inside a tmux session
func InTmux() bool {
	return os.Getenv("TMUX") != ""
}

// HasTmux returns true if tmux is available in PATH
func HasTmux() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

// CurrentPaneID returns the current pane ID
func CurrentPaneID() (string, error) {
	return runTmux("display-message", "-p", "#{pane_id}")
}

// CurrentSessionName returns the current session name
func CurrentSessionName() (string, error) {
	return runTmux("display-message", "-p", "#{session_name}")
}

// FindSlipPane finds the Slip pane if it exists
func FindSlipPane() (string, bool, error) {
	out, err := runTmux("list-panes", "-s", "-F", "#{pane_id} #{"+slipPaneOption+"}")
	if err != nil {
		return "", false, err
	}
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		if fields[1] == "1" {
			return fields[0], true, nil
		}
	}
	return "", false, nil
}

// KillSlipPane kills the Slip pane if it exists
func KillSlipPane() error {
	paneID, ok, err := FindSlipPane()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	_, err = runTmux("kill-pane", "-t", paneID)
	return err
}

// SplitSlip creates a new pane for Slip
func SplitSlip(dock string, size int, command []string) (string, error) {
	if len(command) == 0 {
		return "", errors.New("missing command")
	}

	origPane, _ := CurrentPaneID()

	args := []string{"split-window", "-P", "-F", "#{pane_id}"}

	switch strings.ToLower(dock) {
	case "bottom":
		args = append(args, "-v")
	case "left":
		args = append(args, "-h", "-b")
	case "top":
		args = append(args, "-v", "-b")
	default: // right
		args = append(args, "-h")
	}

	if size > 0 {
		args = append(args, "-l", strconv.Itoa(size))
	}

	// Quote each part of the command and join them
	// This handles paths with spaces correctly
	var quotedParts []string
	for _, part := range command {
		// Quote parts that contain spaces
		if strings.Contains(part, " ") {
			quotedParts = append(quotedParts, "\""+part+"\"")
		} else {
			quotedParts = append(quotedParts, part)
		}
	}
	cmdStr := strings.Join(quotedParts, " ")
	args = append(args, cmdStr)

	paneID, err := runTmux(args...)
	if err != nil {
		return "", fmt.Errorf("split-window failed: %w (command: %s)", err, cmdStr)
	}
	if paneID == "" {
		return "", errors.New("no pane id returned")
	}
	paneID = strings.TrimSpace(paneID)

	// Tag the pane
	if _, err := runTmux("set-option", "-p", "-t", paneID, slipPaneOption, "1"); err != nil {
		// Non-fatal, pane was still created
		return paneID, nil
	}

	// Return focus to original pane
	if origPane != "" {
		_, _ = runTmux("select-pane", "-t", origPane)
	}

	return paneID, nil
}

// ResizePane resizes a pane
func ResizePane(paneID string, width, height int) error {
	if width > 0 {
		_, err := runTmux("resize-pane", "-t", paneID, "-x", strconv.Itoa(width))
		if err != nil {
			return err
		}
	}
	if height > 0 {
		_, err := runTmux("resize-pane", "-t", paneID, "-y", strconv.Itoa(height))
		if err != nil {
			return err
		}
	}
	return nil
}

// SelectPane selects/focuses a pane
func SelectPane(paneID string) error {
	_, err := runTmux("select-pane", "-t", paneID)
	return err
}

// SendKeys sends keys to a pane
func SendKeys(paneID string, keys string) error {
	_, err := runTmux("send-keys", "-t", paneID, keys)
	return err
}

// GetPaneSize returns the size of a pane
func GetPaneSize(paneID string) (int, int, error) {
	out, err := runTmux("display-message", "-t", paneID, "-p", "#{pane_width} #{pane_height}")
	if err != nil {
		return 0, 0, err
	}
	parts := strings.Fields(out)
	if len(parts) != 2 {
		return 0, 0, errors.New("invalid pane size output")
	}
	width, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	height, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}
	return width, height, nil
}

// ListPanes lists all panes in the current session
func ListPanes() ([]PaneInfo, error) {
	out, err := runTmux("list-panes", "-s", "-F", "#{pane_id}\t#{pane_width}\t#{pane_height}\t#{pane_current_command}\t#{"+slipPaneOption+"}")
	if err != nil {
		return nil, err
	}

	var panes []PaneInfo
	for _, line := range strings.Split(out, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 5 {
			continue
		}
		width, _ := strconv.Atoi(parts[1])
		height, _ := strconv.Atoi(parts[2])
		panes = append(panes, PaneInfo{
			ID:      parts[0],
			Width:   width,
			Height:  height,
			Command: parts[3],
			IsSlip:  parts[4] == "1",
		})
	}
	return panes, nil
}

// PaneInfo contains information about a tmux pane
type PaneInfo struct {
	ID      string
	Width   int
	Height  int
	Command string
	IsSlip  bool
}

// Version returns the tmux version
func Version() (string, error) {
	return runTmux("-V")
}

// runTmux executes a tmux command and returns the output
func runTmux(args ...string) (string, error) {
	cmd := exec.Command("tmux", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(output)), err
	}
	return strings.TrimSpace(string(output)), nil
}
