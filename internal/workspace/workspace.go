package workspace

// LayoutType represents the multiplexer type
type LayoutType string

const (
	LayoutNone   LayoutType = "none"
	LayoutZellij LayoutType = "zellij"
	LayoutTmux   LayoutType = "tmux"
)

// Pane represents a single pane in a layout
type Pane struct {
	Command   string `json:"command"`
	Direction string `json:"direction"` // "down", "right"
}

// Layout represents the pane arrangement for a workspace
type Layout struct {
	Type    LayoutType `json:"type"`
	MainCmd string     `json:"main_cmd,omitempty"` // Command for the original (left/top) pane
	Panes   []Pane     `json:"panes"`
}

// Workspace represents a terminal workspace
type Workspace struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Layout Layout `json:"layout"`
}

// NewWorkspace creates a new workspace with default values
func NewWorkspace(name, path string) *Workspace {
	return &Workspace{
		Name: name,
		Path: path,
		Layout: Layout{
			Type:  LayoutNone,
			Panes: []Pane{},
		},
	}
}
