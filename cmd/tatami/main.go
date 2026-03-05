package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/OleksandrBesan/tatami/internal/config"
	"github.com/OleksandrBesan/tatami/internal/shell"
	"github.com/OleksandrBesan/tatami/internal/tui"
	"github.com/OleksandrBesan/tatami/internal/workspace"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("tatami %s\n", version)
		return
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load config paths
	paths, err := config.GetPaths()
	if err != nil {
		return fmt.Errorf("failed to get config paths: %w", err)
	}

	// Initialize workspace store
	store, err := workspace.NewStore(paths)
	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

	// Create and run the TUI app
	app := tui.NewApp(store)
	p := tea.NewProgram(app, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	// Handle result
	finalApp, ok := finalModel.(*tui.App)
	if !ok {
		return nil
	}

	result := finalApp.Result()
	if result == nil {
		return nil
	}

	return handleResult(result)
}

func copyToClipboard(text string) error {
	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

func handleResult(result *tui.Result) error {
	ws := result.Workspace
	zellij := shell.NewZellijRunner()
	tmux := shell.NewTmuxRunner()

	switch result.Action {
	case tui.ActionCD:
		// Check if shell wrapper is active
		if os.Getenv("TATAMI_WRAPPER") == "1" {
			// Wrapper will handle cd
			fmt.Println(ws.Path)
			return nil
		}
		// If inside Zellij, write cd command to current pane
		if zellij.IsInsideSession() {
			cdCmd := fmt.Sprintf("cd %s\n", ws.Path)
			return zellij.WriteChars(cdCmd)
		}
		// If inside Tmux, send keys to current pane
		if tmux.IsInsideSession() {
			cdCmd := fmt.Sprintf("cd %s", ws.Path)
			return tmux.SendKeys(cdCmd)
		}
		// Fallback - copy to clipboard
		cdCmd := fmt.Sprintf("cd %s", ws.Path)
		if err := copyToClipboard(cdCmd); err == nil {
			fmt.Printf("%s  (copied to clipboard, paste to run)\n", cdCmd)
		} else {
			fmt.Println(cdCmd)
		}
		return nil

	case tui.ActionNewTab:
		if zellij.IsInsideSession() {
			return zellij.NewTab(ws.Path, ws.Name)
		}
		if tmux.IsInsideSession() {
			return tmux.NewWindow(ws.Path, ws.Name)
		}
		fmt.Fprintf(os.Stderr, "Not inside a Zellij or Tmux session\n")
		return nil

	case tui.ActionNewPane:
		if zellij.IsInsideSession() {
			return zellij.NewPane(ws.Path, "down")
		}
		if tmux.IsInsideSession() {
			return tmux.NewPane(ws.Path, "down")
		}
		fmt.Fprintf(os.Stderr, "Not inside a Zellij or Tmux session\n")
		return nil

	case tui.ActionWithLayout:
		if zellij.IsInsideSession() && ws.Layout.Type == workspace.LayoutZellij {
			return zellij.RunWithLayout(ws)
		}
		if tmux.IsInsideSession() && ws.Layout.Type == workspace.LayoutTmux {
			return tmux.RunWithLayout(ws)
		}
		fmt.Fprintf(os.Stderr, "Layout type mismatch or not inside session\n")
		return nil

	case tui.ActionWithTemplate:
		if result.Template == nil {
			return fmt.Errorf("no template selected")
		}
		// Create a temporary workspace with template panes
		tmplWs := &workspace.Workspace{
			Name: ws.Name,
			Path: ws.Path,
			Layout: workspace.Layout{
				MainCmd: result.Template.MainCmd,
				Panes:   result.Template.Panes,
			},
		}
		if zellij.IsInsideSession() {
			tmplWs.Layout.Type = workspace.LayoutZellij
			return zellij.RunWithLayout(tmplWs)
		}
		if tmux.IsInsideSession() {
			tmplWs.Layout.Type = workspace.LayoutTmux
			return tmux.RunWithLayout(tmplWs)
		}
		fmt.Fprintf(os.Stderr, "Not inside a Zellij or Tmux session\n")
		return nil
	}

	return nil
}
