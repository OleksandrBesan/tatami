package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/OleksandrBesan/tatami/internal/shell"
	"github.com/OleksandrBesan/tatami/internal/workspace"
)

// View represents the current view state
type View int

const (
	ViewList View = iota
	ViewCreate
	ViewActions
	ViewLayout
	ViewTemplates
)

// Result represents the outcome of the TUI session
type Result struct {
	Action    Action
	Workspace *workspace.Workspace
	Template  *workspace.Template
}

// App is the main Bubbletea model
type App struct {
	store        *workspace.Store
	zellij       *shell.ZellijRunner
	tmux         *shell.TmuxRunner
	currentView  View
	previousView View
	listView     *ListView
	createView   *CreateView
	actionsView  *ActionView
	layoutEditor *LayoutEditor
	templateView *TemplateView
	result       *Result
	width        int
	height       int
	err          error
}

// NewApp creates a new App
func NewApp(store *workspace.Store) *App {
	zellij := shell.NewZellijRunner()
	tmux := shell.NewTmuxRunner()

	return &App{
		store:        store,
		zellij:       zellij,
		tmux:         tmux,
		currentView:  ViewList,
		listView:     NewListView(store.List()),
		createView:   NewCreateView(),
		layoutEditor: NewLayoutEditor(),
	}
}

// Result returns the result of the TUI session
func (a *App) Result() *Result {
	return a.result
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.listView.SetSize(msg.Width, msg.Height)
		return a, nil

	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			return a, tea.Quit
		}

		// View-specific handling
		switch a.currentView {
		case ViewList:
			return a.updateList(msg)
		case ViewCreate:
			return a.updateCreate(msg)
		case ViewActions:
			return a.updateActions(msg)
		case ViewLayout:
			return a.updateLayout(msg)
		case ViewTemplates:
			return a.updateTemplates(msg)
		}
	}

	return a, nil
}

func (a *App) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle filter mode
	if a.listView.IsFiltering() {
		switch msg.String() {
		case "enter":
			a.listView.StopFiltering()
			return a, nil
		case "esc":
			a.listView.ClearFilter()
			return a, nil
		default:
			return a, a.listView.Update(msg)
		}
	}

	switch msg.String() {
	case "q", "esc":
		return a, tea.Quit

	case "enter":
		if ws := a.listView.Selected(); ws != nil {
			a.actionsView = NewActionView(ws, a.zellij.IsInsideSession(), a.tmux.IsInsideSession())
			a.currentView = ViewActions
		}
		return a, nil

	case "n":
		a.createView.Reset()
		a.currentView = ViewCreate
		return a, nil

	case "e":
		if ws := a.listView.Selected(); ws != nil {
			a.createView.EditWorkspace(ws)
			a.currentView = ViewCreate
		}
		return a, nil

	case "d":
		if ws := a.listView.Selected(); ws != nil {
			if err := a.store.Delete(ws.Name); err == nil {
				a.listView.SetWorkspaces(a.store.List())
			}
		}
		return a, nil

	default:
		return a, a.listView.Update(msg)
	}
}

func (a *App) updateCreate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// If in layout editor, go back to create form
		if a.currentView == ViewLayout {
			a.currentView = ViewCreate
			return a, nil
		}
		a.currentView = ViewList
		return a, nil

	case "enter":
		// Don't submit if we're in layout mode
		if a.currentView == ViewLayout && a.layoutEditor.IsEditing() {
			return a, a.layoutEditor.Update(msg)
		}

		a.createView.Validate()
		if a.createView.errorMsg != "" {
			return a, nil
		}

		ws := a.createView.GetWorkspace()

		var err error
		if a.createView.IsEditing() {
			err = a.store.Update(a.createView.EditingName(), ws)
		} else {
			err = a.store.Create(ws)
		}

		if err != nil {
			a.createView.SetError(err.Error())
			return a, nil
		}

		a.listView.SetWorkspaces(a.store.List())
		a.currentView = ViewList
		return a, nil

	case "ctrl+l":
		// Open layout editor
		ws := a.createView.GetWorkspace()
		a.layoutEditor.SetPanes(ws.Layout.Panes)
		a.currentView = ViewLayout
		return a, nil

	case "ctrl+t":
		// Open template picker
		a.templateView = NewTemplateView()
		a.previousView = ViewCreate
		a.currentView = ViewTemplates
		return a, nil

	default:
		return a, a.createView.Update(msg)
	}
}

func (a *App) updateActions(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		a.currentView = ViewList
		return a, nil

	case "enter":
		action := a.actionsView.Selected()
		ws := a.actionsView.Workspace()

		// If template action, show template picker
		if action == ActionWithTemplate {
			a.templateView = NewTemplateView()
			a.previousView = ViewActions
			a.currentView = ViewTemplates
			return a, nil
		}

		a.result = &Result{
			Action:    action,
			Workspace: ws,
		}
		return a, tea.Quit

	default:
		return a, a.actionsView.Update(msg)
	}
}

func (a *App) updateLayout(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if a.layoutEditor.IsEditing() {
			// Cancel pane edit
			return a, a.layoutEditor.Update(msg)
		}
		// Save panes and go back
		ws := a.createView.GetWorkspace()
		ws.Layout.Panes = a.layoutEditor.GetPanes()
		a.createView.EditWorkspace(ws)
		a.createView.editing = true // Keep edit mode
		a.currentView = ViewCreate
		return a, nil

	default:
		return a, a.layoutEditor.Update(msg)
	}
}

func (a *App) updateTemplates(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		a.currentView = a.previousView
		return a, nil

	case "enter":
		tmpl := a.templateView.Selected()

		// If came from create view, apply template and go back
		if a.previousView == ViewCreate {
			a.createView.ApplyTemplate(tmpl)
			a.currentView = ViewCreate
			return a, nil
		}

		// If came from actions view, execute with template
		ws := a.actionsView.Workspace()
		a.result = &Result{
			Action:    ActionWithTemplate,
			Workspace: ws,
			Template:  tmpl,
		}
		return a, tea.Quit

	default:
		return a, a.templateView.Update(msg)
	}
}

// View implements tea.Model
func (a *App) View() string {
	if a.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v", a.err))
	}

	switch a.currentView {
	case ViewList:
		return a.listView.View()
	case ViewCreate:
		return a.createView.View()
	case ViewActions:
		return a.actionsView.View()
	case ViewLayout:
		return boxStyle.Render(a.layoutEditor.View())
	case ViewTemplates:
		return a.templateView.View()
	default:
		return ""
	}
}
