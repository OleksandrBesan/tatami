package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/OleksandrBesan/tatami/internal/workspace"
)

// ListView displays the list of workspaces
type ListView struct {
	workspaces []workspace.Workspace
	filtered   []workspace.Workspace
	cursor     int
	filter     textinput.Model
	filtering  bool
	width      int
	height     int
}

// NewListView creates a new list view
func NewListView(workspaces []workspace.Workspace) *ListView {
	ti := textinput.New()
	ti.Placeholder = "Filter workspaces..."
	ti.CharLimit = 50

	return &ListView{
		workspaces: workspaces,
		filtered:   workspaces,
		cursor:     0,
		filter:     ti,
		filtering:  false,
	}
}

// SetWorkspaces updates the workspace list
func (l *ListView) SetWorkspaces(workspaces []workspace.Workspace) {
	l.workspaces = workspaces
	l.applyFilter()
}

// SetSize sets the view dimensions
func (l *ListView) SetSize(width, height int) {
	l.width = width
	l.height = height
}

// Selected returns the currently selected workspace
func (l *ListView) Selected() *workspace.Workspace {
	if len(l.filtered) == 0 {
		return nil
	}
	return &l.filtered[l.cursor]
}

func (l *ListView) applyFilter() {
	query := l.filter.Value()
	if query == "" {
		l.filtered = l.workspaces
	} else {
		query = strings.ToLower(query)
		l.filtered = nil
		for _, ws := range l.workspaces {
			if strings.Contains(strings.ToLower(ws.Name), query) ||
				strings.Contains(strings.ToLower(ws.Path), query) {
				l.filtered = append(l.filtered, ws)
			}
		}
	}

	// Adjust cursor if needed
	if l.cursor >= len(l.filtered) {
		l.cursor = max(0, len(l.filtered)-1)
	}
}

// Update handles input for the list view
func (l *ListView) Update(msg tea.Msg) tea.Cmd {
	if l.filtering {
		var cmd tea.Cmd
		l.filter, cmd = l.filter.Update(msg)
		l.applyFilter()
		return cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if l.cursor < len(l.filtered)-1 {
				l.cursor++
			}
		case "k", "up":
			if l.cursor > 0 {
				l.cursor--
			}
		case "g":
			l.cursor = 0
		case "G":
			l.cursor = max(0, len(l.filtered)-1)
		case "/":
			l.filtering = true
			l.filter.Focus()
			return nil
		}
	}
	return nil
}

// StopFiltering exits filter mode
func (l *ListView) StopFiltering() {
	l.filtering = false
	l.filter.Blur()
}

// ClearFilter resets the filter
func (l *ListView) ClearFilter() {
	l.filter.SetValue("")
	l.applyFilter()
	l.StopFiltering()
}

// IsFiltering returns whether filter mode is active
func (l *ListView) IsFiltering() bool {
	return l.filtering
}

// View renders the list view
func (l *ListView) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("TATAMI - Workspaces"))
	b.WriteString("\n\n")

	// Filter input (if active)
	if l.filtering {
		b.WriteString(l.filter.View())
		b.WriteString("\n\n")
	}

	// Workspace list
	if len(l.filtered) == 0 {
		if len(l.workspaces) == 0 {
			b.WriteString(mutedStyle.Render("No workspaces yet. Press 'n' to create one."))
		} else {
			b.WriteString(mutedStyle.Render("No matching workspaces."))
		}
	} else {
		// Calculate available height for list
		listHeight := l.height - 10 // Reserve space for title, help, etc.
		if listHeight < 5 {
			listHeight = 5
		}

		// Determine visible range
		start := 0
		end := len(l.filtered)
		if end > listHeight {
			// Scroll to keep cursor visible
			if l.cursor >= listHeight {
				start = l.cursor - listHeight + 1
			}
			end = start + listHeight
			if end > len(l.filtered) {
				end = len(l.filtered)
				start = end - listHeight
			}
		}

		for i := start; i < end; i++ {
			ws := l.filtered[i]
			cursor := "  "
			style := normalStyle
			if i == l.cursor {
				cursor = "> "
				style = selectedStyle
			}

			// Format: name + path
			name := style.Render(ws.Name)
			path := mutedStyle.Render(shortenPath(ws.Path, 40))
			line := fmt.Sprintf("%s%-20s %s", cursor, name, path)

			// Add layout indicator
			if ws.Layout.Type != workspace.LayoutNone && len(ws.Layout.Panes) > 0 {
				indicator := mutedStyle.Render(fmt.Sprintf(" [%s:%d]", ws.Layout.Type, len(ws.Layout.Panes)))
				line += indicator
			}

			b.WriteString(line)
			b.WriteString("\n")
		}

		// Show scroll indicator if needed
		if len(l.filtered) > listHeight {
			scrollInfo := fmt.Sprintf(" (%d/%d)", l.cursor+1, len(l.filtered))
			b.WriteString(mutedStyle.Render(scrollInfo))
			b.WriteString("\n")
		}
	}

	// Help text
	help := "[n]ew  [e]dit  [d]elete  [/]filter  [q]uit"
	if l.filtering {
		help = "[enter]confirm  [esc]cancel"
	}
	b.WriteString(helpStyle.Render(help))

	return lipgloss.NewStyle().Padding(1, 2).Render(b.String())
}

func shortenPath(path string, maxLen int) string {
	if len(path) <= maxLen {
		return path
	}

	// Try to replace home directory with ~
	home, _ := strings.CutPrefix(path, "/Users/")
	if home != path {
		parts := strings.SplitN(home, "/", 2)
		if len(parts) == 2 {
			path = "~/" + parts[1]
		}
	}

	if len(path) <= maxLen {
		return path
	}

	// Truncate from the beginning
	return "..." + path[len(path)-maxLen+3:]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
