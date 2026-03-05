package tui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// PathPicker provides directory autocomplete functionality
type PathPicker struct {
	input       textinput.Model
	suggestions []string
	cursor      int
	showList    bool
}

// NewPathPicker creates a new path picker
func NewPathPicker() *PathPicker {
	ti := textinput.New()
	ti.Placeholder = "~/projects/myapp"
	ti.CharLimit = 200
	ti.Width = 50

	return &PathPicker{
		input:       ti,
		suggestions: nil,
		cursor:      0,
		showList:    false,
	}
}

// Focus focuses the path picker
func (p *PathPicker) Focus() {
	p.input.Focus()
}

// Blur unfocuses the path picker
func (p *PathPicker) Blur() {
	p.input.Blur()
}

// Value returns the current path value
func (p *PathPicker) Value() string {
	return expandPath(p.input.Value())
}

// SetValue sets the path value
func (p *PathPicker) SetValue(path string) {
	p.input.SetValue(path)
	p.updateSuggestions()
}

// Update handles input
func (p *PathPicker) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if len(p.suggestions) > 0 {
				if p.showList {
					// Cycle through suggestions
					p.cursor = (p.cursor + 1) % len(p.suggestions)
				} else {
					p.showList = true
					p.cursor = 0
				}
				// Apply current suggestion
				p.input.SetValue(p.suggestions[p.cursor])
				p.input.CursorEnd()
			}
			return nil

		case "shift+tab":
			if len(p.suggestions) > 0 && p.showList {
				p.cursor--
				if p.cursor < 0 {
					p.cursor = len(p.suggestions) - 1
				}
				p.input.SetValue(p.suggestions[p.cursor])
				p.input.CursorEnd()
			}
			return nil

		case "ctrl+u":
			p.input.SetValue("")
			p.showList = false
			p.suggestions = nil
			return nil

		default:
			// Hide suggestion list on any other key
			if p.showList && msg.String() != "up" && msg.String() != "down" {
				p.showList = false
			}
		}
	}

	var cmd tea.Cmd
	prevValue := p.input.Value()
	p.input, cmd = p.input.Update(msg)

	// Update suggestions if value changed
	if p.input.Value() != prevValue {
		p.updateSuggestions()
	}

	return cmd
}

func (p *PathPicker) updateSuggestions() {
	path := expandPath(p.input.Value())
	p.suggestions = getPathSuggestions(path)
	p.cursor = 0
}

// View renders the path picker
func (p *PathPicker) View() string {
	var b strings.Builder

	b.WriteString(p.input.View())

	if p.showList && len(p.suggestions) > 0 {
		b.WriteString("\n")
		maxShow := 5
		if len(p.suggestions) < maxShow {
			maxShow = len(p.suggestions)
		}

		for i := 0; i < maxShow; i++ {
			prefix := "  "
			if i == p.cursor {
				prefix = "> "
			}
			b.WriteString(mutedStyle.Render(prefix + shortenPath(p.suggestions[i], 50)))
			b.WriteString("\n")
		}

		if len(p.suggestions) > maxShow {
			b.WriteString(mutedStyle.Render(strings.Repeat(" ", 2) + "..."))
		}
	}

	return b.String()
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		return home
	}
	return path
}

func getPathSuggestions(path string) []string {
	if path == "" {
		home, _ := os.UserHomeDir()
		path = home
	}

	// If path ends with /, list contents
	// Otherwise, list siblings that match prefix
	dir := path
	prefix := ""

	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		dir = filepath.Dir(path)
		prefix = filepath.Base(path)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var suggestions []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue // Skip hidden directories by default
		}
		if prefix != "" && !strings.HasPrefix(strings.ToLower(name), strings.ToLower(prefix)) {
			continue
		}
		suggestions = append(suggestions, filepath.Join(dir, name))
	}

	sort.Strings(suggestions)
	return suggestions
}
