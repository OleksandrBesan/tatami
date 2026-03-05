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
			Name:        "nvim-right",
			Description: "term LEFT, nvim RIGHT",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
			},
		},
		{
			Name:        "nvim-bottom",
			Description: "term TOP, nvim BOTTOM",
			Panes: []Pane{
				{Command: "nvim", Direction: "down"},
			},
		},
		{
			Name:        "lazygit-right",
			Description: "term LEFT, lazygit RIGHT",
			Panes: []Pane{
				{Command: "lazygit", Direction: "right"},
			},
		},
		{
			Name:        "nvim-right-term-below",
			Description: "term LEFT, nvim RIGHT TOP, term RIGHT BOTTOM",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
				{Command: "", Direction: "down"},
			},
		},
		{
			Name:        "dev-stack",
			Description: "term LEFT, nvim RIGHT TOP, lazygit RIGHT BOTTOM",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
				{Command: "lazygit", Direction: "down"},
			},
		},
		{
			Name:        "2-side",
			Description: "term LEFT, term RIGHT",
			Panes: []Pane{
				{Command: "", Direction: "right"},
			},
		},
		{
			Name:        "2-stack",
			Description: "term TOP, term BOTTOM",
			Panes: []Pane{
				{Command: "", Direction: "down"},
			},
		},
		{
			Name:        "3-side",
			Description: "term LEFT, term CENTER, term RIGHT",
			Panes: []Pane{
				{Command: "", Direction: "right"},
				{Command: "", Direction: "right"},
			},
		},
		{
			Name:        "3-stack",
			Description: "term TOP, term MIDDLE, term BOTTOM",
			Panes: []Pane{
				{Command: "", Direction: "down"},
				{Command: "", Direction: "down"},
			},
		},
		{
			Name:        "3-right-stack",
			Description: "term LEFT, term RIGHT TOP, term RIGHT BOTTOM",
			Panes: []Pane{
				{Command: "", Direction: "right"},
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
