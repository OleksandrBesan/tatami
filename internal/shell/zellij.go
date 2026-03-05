package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/oleslab/tatami/internal/workspace"
)

// ZellijRunner executes Zellij commands
type ZellijRunner struct{}

// NewZellijRunner creates a new Zellij runner
func NewZellijRunner() *ZellijRunner {
	return &ZellijRunner{}
}

// IsAvailable checks if Zellij is installed
func (z *ZellijRunner) IsAvailable() bool {
	_, err := exec.LookPath("zellij")
	return err == nil
}

// IsInsideSession checks if we're inside a Zellij session
func (z *ZellijRunner) IsInsideSession() bool {
	return os.Getenv("ZELLIJ") != ""
}

// WriteChars writes text to the current pane
func (z *ZellijRunner) WriteChars(text string) error {
	cmd := exec.Command("zellij", "action", "write-chars", text)
	return cmd.Run()
}

// NewTab opens a new tab in the current Zellij session
func (z *ZellijRunner) NewTab(path, name string) error {
	args := []string{"action", "new-tab", "--cwd", path}
	if name != "" {
		args = append(args, "--name", name)
	}
	cmd := exec.Command("zellij", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// NewPane opens a new pane in the current Zellij session
func (z *ZellijRunner) NewPane(path string, direction string) error {
	// Use "zellij run" which properly supports --cwd
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	args := []string{"run", "--cwd", path}
	if direction != "" {
		args = append(args, "--direction", direction)
	}
	args = append(args, "--", shell)

	cmd := exec.Command("zellij", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunWithLayout opens a workspace with its configured layout
func (z *ZellijRunner) RunWithLayout(ws *workspace.Workspace) error {
	// First, create a new tab
	if err := z.NewTab(ws.Path, ws.Name); err != nil {
		return fmt.Errorf("failed to create tab: %w", err)
	}

	// Then create panes for the layout
	for _, pane := range ws.Layout.Panes {
		if err := z.RunPane(ws.Path, pane.Direction, pane.Command); err != nil {
			return fmt.Errorf("failed to create pane: %w", err)
		}
	}

	// If there's a main command, focus first pane and run it
	if ws.Layout.MainCmd != "" {
		// Focus the first pane
		if err := z.FocusFirstPane(); err != nil {
			return fmt.Errorf("failed to focus first pane: %w", err)
		}
		// Write the command
		if err := z.WriteChars(ws.Layout.MainCmd + "\n"); err != nil {
			return fmt.Errorf("failed to run main command: %w", err)
		}
	}

	return nil
}

// FocusFirstPane focuses the first pane in the current tab
func (z *ZellijRunner) FocusFirstPane() error {
	// Move to first pane by going left/up multiple times
	for i := 0; i < 10; i++ {
		exec.Command("zellij", "action", "move-focus", "left").Run()
	}
	for i := 0; i < 10; i++ {
		exec.Command("zellij", "action", "move-focus", "up").Run()
	}
	return nil
}

// RunPane opens a new pane with an optional command
func (z *ZellijRunner) RunPane(path, direction, command string) error {
	args := []string{"run", "--cwd", path}
	if direction != "" {
		args = append(args, "--direction", direction)
	}
	args = append(args, "--")

	if command != "" {
		args = append(args, "sh", "-c", command)
	} else {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
		args = append(args, shell)
	}

	cmd := exec.Command("zellij", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
