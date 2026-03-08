package shell

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/OleksandrBesan/tatami/internal/workspace"
)

// TmuxRunner executes Tmux commands
type TmuxRunner struct{}

// NewTmuxRunner creates a new Tmux runner
func NewTmuxRunner() *TmuxRunner {
	return &TmuxRunner{}
}

// IsAvailable checks if Tmux is installed
func (t *TmuxRunner) IsAvailable() bool {
	_, err := exec.LookPath("tmux")
	return err == nil
}

// IsInsideSession checks if we're inside a Tmux session
func (t *TmuxRunner) IsInsideSession() bool {
	return os.Getenv("TMUX") != ""
}

// SendKeys sends keys to the current pane
func (t *TmuxRunner) SendKeys(text string) error {
	cmd := exec.Command("tmux", "send-keys", text, "Enter")
	return cmd.Run()
}

// NewWindow opens a new window in the current Tmux session
func (t *TmuxRunner) NewWindow(path, name string) error {
	args := []string{"new-window", "-c", path}
	if name != "" {
		args = append(args, "-n", name)
	}
	cmd := exec.Command("tmux", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// NewPane opens a new pane in the current Tmux session
func (t *TmuxRunner) NewPane(path string, direction string) error {
	args := []string{"split-window", "-c", path}
	switch direction {
	case "down":
		args = append(args, "-v")
	case "right":
		args = append(args, "-h")
	}
	cmd := exec.Command("tmux", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunWithLayout opens a workspace with its configured layout
func (t *TmuxRunner) RunWithLayout(ws *workspace.Workspace) error {
	// First, create a new window
	if err := t.NewWindow(ws.Path, ws.Name); err != nil {
		return fmt.Errorf("failed to create window: %w", err)
	}

	// Then create panes for the layout
	for _, pane := range ws.Layout.Panes {
		if err := t.NewPane(ws.Path, pane.Direction); err != nil {
			return fmt.Errorf("failed to create pane: %w", err)
		}

		// Run command in the new pane if specified
		if pane.Command != "" {
			cmd := exec.Command("tmux", "send-keys", pane.Command, "Enter")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("failed to run command in pane: %w", err)
			}
		}
	}

	return nil
}

// buildTmuxSSHCmd builds an SSH command string with optional key and command
func buildTmuxSSHCmd(host, key, remotePath, command string) string {
	var sshPart string
	if key != "" {
		sshPart = fmt.Sprintf("ssh -i %s %s", key, host)
	} else {
		sshPart = fmt.Sprintf("ssh %s", host)
	}

	if command != "" {
		return fmt.Sprintf("%s -t 'cd %s && %s'", sshPart, remotePath, command)
	}
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "$SHELL"
	}
	return fmt.Sprintf("%s -t 'cd %s && %s'", sshPart, remotePath, shell)
}

// NewWindowSSH opens a new window with an SSH session to a remote host
func (t *TmuxRunner) NewWindowSSH(host, key, remotePath, name string) error {
	sshCmd := buildTmuxSSHCmd(host, key, remotePath, "")

	args := []string{"new-window"}
	if name != "" {
		args = append(args, "-n", name)
	}
	args = append(args, "sh", "-c", sshCmd)

	cmd := exec.Command("tmux", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// NewPaneSSH opens a new pane with an SSH session to a remote host
func (t *TmuxRunner) NewPaneSSH(host, key, remotePath, direction string) error {
	sshCmd := buildTmuxSSHCmd(host, key, remotePath, "")

	args := []string{"split-window"}
	switch direction {
	case "down":
		args = append(args, "-v")
	case "right":
		args = append(args, "-h")
	}
	args = append(args, "sh", "-c", sshCmd)

	cmd := exec.Command("tmux", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunPaneSSH opens a new pane with SSH and runs a command
func (t *TmuxRunner) RunPaneSSH(host, key, remotePath, direction, command string) error {
	sshCmd := buildTmuxSSHCmd(host, key, remotePath, command)

	args := []string{"split-window"}
	switch direction {
	case "down":
		args = append(args, "-v")
	case "right":
		args = append(args, "-h")
	}
	args = append(args, "sh", "-c", sshCmd)

	cmd := exec.Command("tmux", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunWithLayoutSSH opens a remote workspace with its configured layout via SSH
func (t *TmuxRunner) RunWithLayoutSSH(ws *workspace.Workspace) error {
	host := ws.Remote.Host
	key := ws.Remote.Key
	remotePath := ws.Remote.Path

	// First, create a new window with SSH
	if err := t.NewWindowSSH(host, key, remotePath, ws.Name); err != nil {
		return fmt.Errorf("failed to create window: %w", err)
	}

	// Then create panes for the layout
	for _, pane := range ws.Layout.Panes {
		if err := t.RunPaneSSH(host, key, remotePath, pane.Direction, pane.Command); err != nil {
			return fmt.Errorf("failed to create pane: %w", err)
		}
	}

	return nil
}
