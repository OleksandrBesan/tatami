package workspace

// Template represents a predefined layout template
type Template struct {
	Name        string
	Description string
	Panes       []Pane
}

// GetTemplates returns all available layout templates
func GetTemplates() []Template {
	return []Template{
		{
			Name:        "editor-right",
			Description: "Editor left, terminal right",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
			},
		},
		{
			Name:        "editor-left",
			Description: "Terminal left, editor right",
			Panes: []Pane{
				{Command: "", Direction: "right"},
				{Command: "nvim", Direction: ""},
			},
		},
		{
			Name:        "editor-top",
			Description: "Editor top, terminal bottom",
			Panes: []Pane{
				{Command: "nvim", Direction: "down"},
			},
		},
		{
			Name:        "editor-bottom",
			Description: "Terminal top, editor bottom",
			Panes: []Pane{
				{Command: "", Direction: "down"},
				{Command: "nvim", Direction: ""},
			},
		},
		{
			Name:        "editor-2term",
			Description: "Editor left, 2 terminals right",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
				{Command: "", Direction: "down"},
			},
		},
		{
			Name:        "lazygit-right",
			Description: "Editor left, lazygit right",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
				{Command: "lazygit", Direction: ""},
			},
		},
		{
			Name:        "dev-full",
			Description: "Editor, terminal, lazygit",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
				{Command: "", Direction: "down"},
				{Command: "lazygit", Direction: ""},
			},
		},
		{
			Name:        "terminals-h",
			Description: "2 terminals side by side",
			Panes: []Pane{
				{Command: "", Direction: "right"},
			},
		},
		{
			Name:        "terminals-v",
			Description: "2 terminals stacked",
			Panes: []Pane{
				{Command: "", Direction: "down"},
			},
		},
		{
			Name:        "terminals-quad",
			Description: "4 terminals in grid",
			Panes: []Pane{
				{Command: "", Direction: "right"},
				{Command: "", Direction: "down"},
				{Command: "", Direction: "down"},
			},
		},
	}
}

// GetTemplate returns a template by name
func GetTemplate(name string) *Template {
	for _, t := range GetTemplates() {
		if t.Name == name {
			return &t
		}
	}
	return nil
}
