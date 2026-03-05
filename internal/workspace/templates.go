package workspace

// Template represents a predefined layout template
type Template struct {
	Name        string
	Description string
	Panes       []Pane
}

// GetTemplates returns all available layout templates
// Main pane (left/top) is always a terminal
// Direction: "right" = new pane RIGHT, "down" = new pane BELOW
func GetTemplates() []Template {
	return []Template{
		{
			Name:        "term+nvim",
			Description: "[term | nvim]",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
			},
		},
		{
			Name:        "term+lazygit",
			Description: "[term | lazygit]",
			Panes: []Pane{
				{Command: "lazygit", Direction: "right"},
			},
		},
		{
			Name:        "term/nvim",
			Description: "[term / nvim]",
			Panes: []Pane{
				{Command: "nvim", Direction: "down"},
			},
		},
		{
			Name:        "term+nvim/term",
			Description: "[term | nvim / term]",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
				{Command: "", Direction: "down"},
			},
		},
		{
			Name:        "term+term/term",
			Description: "[term | term / term]",
			Panes: []Pane{
				{Command: "", Direction: "right"},
				{Command: "", Direction: "down"},
			},
		},
		{
			Name:        "term+nvim/lazygit",
			Description: "[term | nvim / lazygit]",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
				{Command: "lazygit", Direction: "down"},
			},
		},
		{
			Name:        "2h",
			Description: "[term | term]",
			Panes: []Pane{
				{Command: "", Direction: "right"},
			},
		},
		{
			Name:        "2v",
			Description: "[term / term]",
			Panes: []Pane{
				{Command: "", Direction: "down"},
			},
		},
		{
			Name:        "3h",
			Description: "[term | term | term]",
			Panes: []Pane{
				{Command: "", Direction: "right"},
				{Command: "", Direction: "right"},
			},
		},
		{
			Name:        "3v",
			Description: "[term / term / term]",
			Panes: []Pane{
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
