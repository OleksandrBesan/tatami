package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/oleslab/tatami/internal/workspace"
)

// TemplateView displays the template picker
type TemplateView struct {
	templates []workspace.Template
	cursor    int
}

// NewTemplateView creates a new template view
func NewTemplateView() *TemplateView {
	return &TemplateView{
		templates: workspace.GetTemplates(),
		cursor:    0,
	}
}

// Selected returns the currently selected template
func (t *TemplateView) Selected() *workspace.Template {
	if len(t.templates) == 0 {
		return nil
	}
	return &t.templates[t.cursor]
}

// Update handles input for the template view
func (t *TemplateView) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if t.cursor < len(t.templates)-1 {
				t.cursor++
			}
		case "k", "up":
			if t.cursor > 0 {
				t.cursor--
			}
		case "g":
			t.cursor = 0
		case "G":
			if len(t.templates) > 0 {
				t.cursor = len(t.templates) - 1
			}
		}
	}
	return nil
}

// View renders the template view
func (t *TemplateView) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Choose Layout Template"))
	b.WriteString("\n\n")

	for i, tmpl := range t.templates {
		cursor := "  "
		style := normalStyle
		if i == t.cursor {
			cursor = "> "
			style = selectedStyle
		}

		name := style.Render(tmpl.Name)
		desc := mutedStyle.Render(" - " + tmpl.Description)
		b.WriteString(cursor + name + desc + "\n")
	}

	b.WriteString(helpStyle.Render("\n[enter]select  [esc]cancel"))

	return boxStyle.Render(b.String())
}
