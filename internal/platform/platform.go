package platform

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// IsWindows returns true if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsMacOS returns true if running on macOS
func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// IsLinux returns true if running on Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsWSL returns true if running in Windows Subsystem for Linux
func IsWSL() bool {
	if runtime.GOOS != "linux" {
		return false
	}
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	lower := strings.ToLower(string(data))
	return strings.Contains(lower, "microsoft") || strings.Contains(lower, "wsl")
}

// HasWindowsTerminal checks if Windows Terminal (wt.exe) is available
func HasWindowsTerminal() bool {
	if !IsWindows() {
		return false
	}
	// Check WT_SESSION environment variable
	if os.Getenv("WT_SESSION") != "" {
		return true
	}
	// Check if wt.exe exists in PATH
	_, err := exec.LookPath("wt.exe")
	return err == nil
}

// InWindowsTerminal returns true if currently running inside Windows Terminal
func InWindowsTerminal() bool {
	return os.Getenv("WT_SESSION") != ""
}

// LaunchWindowsTerminal opens a new Windows Terminal tab/pane with the given command
func LaunchWindowsTerminal(command []string, dock string) error {
	wtPath, err := exec.LookPath("wt.exe")
	if err != nil {
		return err
	}

	args := []string{}

	// Determine split direction based on dock
	if InWindowsTerminal() {
		// Split the current window
		switch strings.ToLower(dock) {
		case "bottom":
			args = append(args, "split-pane", "-V")
		case "left":
			args = append(args, "split-pane", "-H", "-d", ".")
		default: // right
			args = append(args, "split-pane", "-H")
		}
	} else {
		// Open new tab
		args = append(args, "new-tab", "--title", "Slip")
	}

	// Add the command
	args = append(args, "--")
	args = append(args, command...)

	cmd := exec.Command(wtPath, args...)
	return cmd.Start()
}

// OpenNewTerminal opens a new terminal window with the given command
func OpenNewTerminal(command []string) error {
	if IsWindows() {
		if HasWindowsTerminal() {
			return LaunchWindowsTerminal(command, "right")
		}
		// Fallback to cmd
		args := []string{"/c", "start", "cmd", "/k"}
		args = append(args, command...)
		cmd := exec.Command("cmd", args...)
		return cmd.Start()
	}

	if IsMacOS() {
		// Use osascript to open Terminal
		script := `tell application "Terminal" to do script "` + strings.Join(command, " ") + `"`
		cmd := exec.Command("osascript", "-e", script)
		return cmd.Start()
	}

	// Linux - try common terminal emulators
	terminals := []string{"gnome-terminal", "konsole", "xterm", "xfce4-terminal"}
	for _, term := range terminals {
		if _, err := exec.LookPath(term); err == nil {
			var cmd *exec.Cmd
			switch term {
			case "gnome-terminal":
				cmd = exec.Command(term, "--", command[0])
				if len(command) > 1 {
					cmd.Args = append(cmd.Args, command[1:]...)
				}
			case "konsole":
				cmd = exec.Command(term, "-e", strings.Join(command, " "))
			default:
				cmd = exec.Command(term, "-e", strings.Join(command, " "))
			}
			return cmd.Start()
		}
	}

	return exec.ErrNotFound
}

// GetArch returns the current architecture
func GetArch() string {
	return runtime.GOARCH
}

// GetOS returns the current OS
func GetOS() string {
	return runtime.GOOS
}

// GetPlatformString returns a string like "linux/amd64"
func GetPlatformString() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}
