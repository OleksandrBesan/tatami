# Tatami

Terminal workspace manager with Zellij/Tmux integration. Quickly switch between projects with predefined layouts.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap OleksandrBesan/tap
brew install tatami
```

### Go Install

```bash
go install github.com/OleksandrBesan/tatami/cmd/tatami@latest
```

### Build from Source

```bash
git clone https://github.com/OleksandrBesan/tatami.git
cd tatami
make install
```

### Download Binary

Download from [Releases](https://github.com/OleksandrBesan/tatami/releases).

## Shell Integration

For `cd` to work in the current terminal, add to `~/.zshrc` or `~/.bashrc`:

```bash
tatami() {
  local output
  output=$(TATAMI_WRAPPER=1 command tatami "$@")
  local exit_code=$?
  if [[ $exit_code -eq 0 && -d "$output" ]]; then
    cd "$output"
  elif [[ -n "$output" ]]; then
    echo "$output"
  fi
  return $exit_code
}
```

Without the wrapper, `cd` will type the command in Zellij or copy to clipboard.

## Usage

```bash
tatami
```

### Keyboard Shortcuts

#### List View
| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `Enter` | Open action menu |
| `n` | New workspace |
| `e` | Edit workspace |
| `d` | Delete workspace |
| `/` | Filter workspaces |
| `q` / `Esc` | Quit |

#### Create/Edit View
| Key | Action |
|-----|--------|
| `Tab` | Autocomplete path / Next field |
| `Ctrl+J` | Next field |
| `Ctrl+K` | Previous field |
| `Ctrl+T` | Choose template |
| `←` / `→` | Change layout type |
| `Enter` | Save |
| `Esc` | Cancel |

#### Action Menu
| Key | Action |
|-----|--------|
| `j` / `k` | Navigate |
| `Enter` | Execute |
| `Esc` | Back |

### Actions

| Action | Description |
|--------|-------------|
| **cd here** | Change to workspace directory |
| **new tab** | Open in new Zellij tab / Tmux window |
| **new pane** | Open in new pane |
| **with template** | Open with a layout template |
| **with saved layout** | Open with workspace's saved layout |

## Layout Templates

Select a template when creating a workspace (`Ctrl+T`) or when opening (`with template`).

| Template | Layout |
|----------|--------|
| `nvim-left` | nvim LEFT, term RIGHT |
| `nvim-left-2term` | nvim LEFT, term RIGHT TOP, term RIGHT BOTTOM |
| `nvim-left-lazygit` | nvim LEFT, lazygit RIGHT |
| `nvim-top` | nvim TOP, term BOTTOM |
| `term-left-nvim` | term LEFT, nvim RIGHT |
| `term-left-lazygit` | term LEFT, lazygit RIGHT |
| `term-left-nvim-lazygit` | term LEFT, nvim RIGHT TOP, lazygit RIGHT BOTTOM |
| `2-side` | term LEFT, term RIGHT |
| `2-stack` | term TOP, term BOTTOM |
| `3-right-stack` | term LEFT, term RIGHT TOP, term RIGHT BOTTOM |

## Configuration

Workspaces are stored in `~/.config/tatami/workspaces.json`:

```json
{
  "workspaces": [
    {
      "name": "myproject",
      "path": "/home/user/projects/myproject",
      "layout": {
        "type": "zellij",
        "main_cmd": "nvim",
        "panes": [
          { "command": "", "direction": "right" },
          { "command": "", "direction": "down" }
        ]
      }
    }
  ]
}
```

### Layout Fields

| Field | Description |
|-------|-------------|
| `type` | `none`, `zellij`, or `tmux` |
| `main_cmd` | Command to run in the original (left/top) pane |
| `panes` | Array of additional panes |
| `panes[].command` | Command to run (empty = shell) |
| `panes[].direction` | `right` or `down` |

## Requirements

- **Zellij** or **Tmux** (for tab/pane features)
- Works without them for basic `cd` functionality

## License

MIT
