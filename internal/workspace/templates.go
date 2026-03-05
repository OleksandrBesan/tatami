package workspace

// Template represents a predefined layout template
type Template struct {
	Name        string
	Description string
	MainCmd     string // Command to run in the original (left/top) pane
	Panes       []Pane
}

// GetTemplates returns all available layout templates
func GetTemplates() []Template {
	return []Template{
		// nvim LEFT layouts
		{
			Name:        "nvim-left",
			Description: "nvim LEFT, term RIGHT",
			MainCmd:     "nvim",
			Panes: []Pane{
				{Command: "", Direction: "right"},
			},
		},
		{
			Name:        "nvim-left-2term",
			Description: "nvim LEFT, term RIGHT TOP, term RIGHT BOTTOM",
			MainCmd:     "nvim",
			Panes: []Pane{
				{Command: "", Direction: "right"},
				{Command: "", Direction: "down"},
			},
		},
		{
			Name:        "nvim-left-lazygit",
			Description: "nvim LEFT, lazygit RIGHT",
			MainCmd:     "nvim",
			Panes: []Pane{
				{Command: "lazygit", Direction: "right"},
			},
		},
		{
			Name:        "nvim-top",
			Description: "nvim TOP, term BOTTOM",
			MainCmd:     "nvim",
			Panes: []Pane{
				{Command: "", Direction: "down"},
			},
		},
		// term LEFT layouts
		{
			Name:        "term-left-nvim",
			Description: "term LEFT, nvim RIGHT",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
			},
		},
		{
			Name:        "term-left-lazygit",
			Description: "term LEFT, lazygit RIGHT",
			Panes: []Pane{
				{Command: "lazygit", Direction: "right"},
			},
		},
		{
			Name:        "term-left-nvim-lazygit",
			Description: "term LEFT, nvim RIGHT TOP, lazygit RIGHT BOTTOM",
			Panes: []Pane{
				{Command: "nvim", Direction: "right"},
				{Command: "lazygit", Direction: "down"},
			},
		},
		// terminal only layouts
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
