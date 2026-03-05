# Tatami

Terminal workspace manager with Zellij/Tmux integration.

## Installation

```bash
go install github.com/oleslab/tatami/cmd/tatami@latest
```

Or build from source:

```bash
go build -o tatami ./cmd/tatami
```

## Shell Integration

For `cd` functionality to work, add the shell wrapper to your shell config:

```bash
# Add to ~/.zshrc or ~/.bashrc
source /path/to/tatami/scripts/tatami.sh
```

This allows the `tatami` command to change your shell's working directory when you select "cd here".

## Usage

Run `tatami` to open the workspace manager.

### Keyboard Shortcuts

#### List View
- `j/k` or `↓/↑` - Navigate workspaces
- `Enter` - Open action menu
- `n` - Create new workspace
- `e` - Edit selected workspace
- `d` - Delete selected workspace
- `/` - Filter workspaces
- `q` or `Esc` - Quit

#### Create/Edit View
- `Tab` - Next field
- `Shift+Tab` - Previous field
- `Enter` - Save workspace
- `Esc` - Cancel
- `Ctrl+L` - Open layout editor

#### Path Picker
- `Tab` - Autocomplete / cycle suggestions
- `Shift+Tab` - Previous suggestion
- `Ctrl+U` - Clear input

#### Action Menu
- `j/k` or `↓/↑` - Navigate actions
- `Enter` - Execute action
- `Esc` - Back to list

### Actions

- **cd here** - Change directory (requires shell wrapper)
- **new tab** - Open workspace in a new Zellij tab or Tmux window
- **new pane** - Open workspace in a new pane
- **with layout** - Open workspace with configured pane layout

## Configuration

Workspaces are stored in `~/.config/tatami/workspaces.json`.

### Layout Configuration

Each workspace can have a layout with multiple panes:

```json
{
  "workspaces": [
    {
      "name": "myproject",
      "path": "/home/user/projects/myproject",
      "layout": {
        "type": "zellij",
        "panes": [
          { "command": "nvim", "direction": "down" },
          { "command": "lazygit", "direction": "right" }
        ]
      }
    }
  ]
}
```

Layout types:
- `none` - No layout (default)
- `zellij` - Use Zellij for pane management
- `tmux` - Use Tmux for pane management

Pane directions:
- `down` - Split horizontally (new pane below)
- `right` - Split vertically (new pane to the right)
