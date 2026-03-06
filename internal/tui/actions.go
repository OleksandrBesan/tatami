package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/OleksandrBesan/tatami/internal/workspace"
)

// Action represents a workspace action
type Action int

const (
	ActionCD Action = iota
	ActionNewTab
	ActionNewPane
	ActionWithLayout
	ActionWithTemplate
)

// ActionView displays the action menu
type ActionView struct {
	workspace *workspace.Workspace
	actions   []Action
	cursor    int
	inZellij  bool
	inTmux    bool
}

// NewActionView creates a new action view
func NewActionView(ws *workspace.Workspace, inZellij, inTmux bool) *ActionView {
	var actions []Action

	if inZellij {
		// Saved layout first (if available)
		if ws.Layout.Type == workspace.LayoutZellij && len(ws.Layout.Panes) > 0 {
			actions = append(actions, ActionWithLayout)
		}
		actions = append(actions, ActionWithTemplate, ActionNewPane, ActionNewTab, ActionCD)
	} else if inTmux {
		// Saved layout first (if available)
		if ws.Layout.Type == workspace.LayoutTmux && len(ws.Layout.Panes) > 0 {
			actions = append(actions, ActionWithLayout)
		}
		actions = append(actions, ActionWithTemplate, ActionNewPane, ActionNewTab, ActionCD)
	} else {
		// Outside multiplexer - only CD
		actions = []Action{ActionCD}
	}

	return &ActionView{
		workspace: ws,
		actions:   actions,
		cursor:    0,
		inZellij:  inZellij,
		inTmux:    inTmux,
	}
}

// Selected returns the currently selected action
func (a *ActionView) Selected() Action {
	return a.actions[a.cursor]
}

// Workspace returns the workspace
func (a *ActionView) Workspace() *workspace.Workspace {
	return a.workspace
}

// Update handles input for the action view
func (a *ActionView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if a.cursor < len(a.actions)-1 {
				a.cursor++
			}
		case "k", "up":
			if a.cursor > 0 {
				a.cursor--
			}
		case "1":
			a.cursor = 0
		case "2":
			if len(a.actions) > 1 {
				a.cursor = 1
			}
		case "3":
			if len(a.actions) > 2 {
				a.cursor = 2
			}
		case "4":
			if len(a.actions) > 3 {
				a.cursor = 3
			}
		}
	}
	return nil
}

// View renders the action view
func (a *ActionView) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Open: " + a.workspace.Name))
	b.WriteString("\n\n")

	actionLabels := map[Action]string{
		ActionCD:           "cd here",
		ActionNewTab:       "new tab",
		ActionNewPane:      "new pane",
		ActionWithTemplate: "with template",
		ActionWithLayout:   "with saved layout",
	}

	for i, action := range a.actions {
		cursor := "  "
		style := normalStyle
		if i == a.cursor {
			cursor = "> "
			style = selectedStyle
		}

		label := actionLabels[action]
		b.WriteString(cursor)
		b.WriteString(style.Render(label))
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("\n[enter]select  [esc]back"))

	return boxStyle.Render(b.String())
}
