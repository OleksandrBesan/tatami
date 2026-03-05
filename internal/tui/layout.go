package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/oleslab/tatami/internal/workspace"
)

type layoutField int

const (
	layoutFieldCommand layoutField = iota
	layoutFieldDirection
)

// LayoutEditor handles editing workspace layouts
type LayoutEditor struct {
	panes         []workspace.Pane
	cursor        int
	editing       bool
	activeField   layoutField
	commandInput  textinput.Model
	directionOpts []string
	directionIdx  int
}

// NewLayoutEditor creates a new layout editor
func NewLayoutEditor() *LayoutEditor {
	cmdInput := textinput.New()
	cmdInput.Placeholder = "nvim"
	cmdInput.CharLimit = 100
	cmdInput.Width = 30

	return &LayoutEditor{
		panes:         nil,
		cursor:        0,
		editing:       false,
		activeField:   layoutFieldCommand,
		commandInput:  cmdInput,
		directionOpts: []string{"down", "right"},
		directionIdx:  0,
	}
}

// SetPanes sets the panes to edit
func (l *LayoutEditor) SetPanes(panes []workspace.Pane) {
	l.panes = make([]workspace.Pane, len(panes))
	copy(l.panes, panes)
	l.cursor = 0
	l.editing = false
}

// GetPanes returns the current panes
func (l *LayoutEditor) GetPanes() []workspace.Pane {
	return l.panes
}

// IsEditing returns whether currently editing a pane
func (l *LayoutEditor) IsEditing() bool {
	return l.editing
}

func (l *LayoutEditor) startEdit() {
	if len(l.panes) == 0 {
		return
	}
	l.editing = true
	l.activeField = layoutFieldCommand
	l.commandInput.SetValue(l.panes[l.cursor].Command)
	l.commandInput.Focus()

	// Set direction index
	for i, dir := range l.directionOpts {
		if dir == l.panes[l.cursor].Direction {
			l.directionIdx = i
			break
		}
	}
}

func (l *LayoutEditor) stopEdit(save bool) {
	if save && len(l.panes) > l.cursor {
		l.panes[l.cursor].Command = l.commandInput.Value()
		l.panes[l.cursor].Direction = l.directionOpts[l.directionIdx]
	}
	l.editing = false
	l.commandInput.Blur()
}

func (l *LayoutEditor) addPane() {
	l.panes = append(l.panes, workspace.Pane{
		Command:   "",
		Direction: "down",
	})
	l.cursor = len(l.panes) - 1
	l.startEdit()
}

func (l *LayoutEditor) deletePane() {
	if len(l.panes) == 0 {
		return
	}
	l.panes = append(l.panes[:l.cursor], l.panes[l.cursor+1:]...)
	if l.cursor >= len(l.panes) && l.cursor > 0 {
		l.cursor--
	}
}

// Update handles input for the layout editor
func (l *LayoutEditor) Update(msg tea.Msg) tea.Cmd {
	if l.editing {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				l.stopEdit(true)
				return nil
			case "esc":
				l.stopEdit(false)
				return nil
			case "tab":
				l.activeField = (l.activeField + 1) % 2
				if l.activeField == layoutFieldCommand {
					l.commandInput.Focus()
				} else {
					l.commandInput.Blur()
				}
				return nil
			case "left", "right":
				if l.activeField == layoutFieldDirection {
					if msg.String() == "right" {
						l.directionIdx = (l.directionIdx + 1) % len(l.directionOpts)
					} else {
						l.directionIdx = (l.directionIdx + len(l.directionOpts) - 1) % len(l.directionOpts)
					}
					return nil
				}
			}
		}

		if l.activeField == layoutFieldCommand {
			var cmd tea.Cmd
			l.commandInput, cmd = l.commandInput.Update(msg)
			return cmd
		}
		return nil
	}

	// Not editing - handle list navigation
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if l.cursor < len(l.panes)-1 {
				l.cursor++
			}
		case "k", "up":
			if l.cursor > 0 {
				l.cursor--
			}
		case "enter", "e":
			l.startEdit()
		case "a":
			l.addPane()
		case "d", "x":
			l.deletePane()
		}
	}
	return nil
}

// View renders the layout editor
func (l *LayoutEditor) View() string {
	var b strings.Builder

	b.WriteString(labelStyle.Render("Layout Panes"))
	b.WriteString("\n\n")

	if len(l.panes) == 0 {
		b.WriteString(mutedStyle.Render("  No panes. Press 'a' to add one."))
		b.WriteString("\n")
	} else {
		for i, pane := range l.panes {
			cursor := "  "
			style := normalStyle
			if i == l.cursor {
				cursor = "> "
				style = selectedStyle
			}

			cmd := pane.Command
			if cmd == "" {
				cmd = "(empty)"
			}
			line := fmt.Sprintf("%s%s [%s]", cursor, cmd, pane.Direction)
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}
	}

	// Edit form
	if l.editing && len(l.panes) > 0 {
		b.WriteString("\n")
		b.WriteString(labelStyle.Render("Edit Pane"))
		b.WriteString("\n")

		cmdLabel := "  Command"
		if l.activeField == layoutFieldCommand {
			cmdLabel = "> Command"
		}
		b.WriteString(mutedStyle.Render(cmdLabel))
		b.WriteString("\n  ")
		b.WriteString(l.commandInput.View())
		b.WriteString("\n")

		dirLabel := "  Direction"
		if l.activeField == layoutFieldDirection {
			dirLabel = "> Direction"
		}
		b.WriteString(mutedStyle.Render(dirLabel))
		b.WriteString("\n  ")

		for i, dir := range l.directionOpts {
			style := mutedStyle
			if i == l.directionIdx {
				style = selectedStyle
			}
			b.WriteString(style.Render("[" + dir + "]"))
			b.WriteString("  ")
		}
		b.WriteString("\n")
	}

	// Help
	help := "[a]dd  [e]dit  [d]elete  [esc]done"
	if l.editing {
		help = "[tab]next  [enter]save  [esc]cancel"
	}
	b.WriteString(helpStyle.Render("\n" + help))

	return b.String()
}
