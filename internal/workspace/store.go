package workspace

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/OleksandrBesan/tatami/internal/config"
)

var (
	ErrWorkspaceNotFound = errors.New("workspace not found")
	ErrWorkspaceExists   = errors.New("workspace already exists")
)

// Store manages workspace persistence
type Store struct {
	paths      *config.Paths
	workspaces []Workspace
}

type storeData struct {
	Workspaces []Workspace `json:"workspaces"`
}

// NewStore creates a new workspace store
func NewStore(paths *config.Paths) (*Store, error) {
	s := &Store{
		paths:      paths,
		workspaces: []Workspace{},
	}

	if err := s.load(); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Store) load() error {
	data, err := os.ReadFile(s.paths.WorkspacesFile)
	if err != nil {
		if os.IsNotExist(err) {
			s.workspaces = []Workspace{}
			return nil
		}
		return err
	}

	var sd storeData
	if err := json.Unmarshal(data, &sd); err != nil {
		return err
	}

	s.workspaces = sd.Workspaces
	return nil
}

func (s *Store) save() error {
	sd := storeData{Workspaces: s.workspaces}
	data, err := json.MarshalIndent(sd, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.paths.WorkspacesFile, data, 0644)
}

// List returns all workspaces
func (s *Store) List() []Workspace {
	return s.workspaces
}

// Get returns a workspace by name
func (s *Store) Get(name string) (*Workspace, error) {
	for i := range s.workspaces {
		if s.workspaces[i].Name == name {
			return &s.workspaces[i], nil
		}
	}
	return nil, ErrWorkspaceNotFound
}

// Create adds a new workspace
func (s *Store) Create(ws *Workspace) error {
	for _, existing := range s.workspaces {
		if existing.Name == ws.Name {
			return ErrWorkspaceExists
		}
	}

	s.workspaces = append(s.workspaces, *ws)
	return s.save()
}

// Update modifies an existing workspace
func (s *Store) Update(name string, ws *Workspace) error {
	for i := range s.workspaces {
		if s.workspaces[i].Name == name {
			s.workspaces[i] = *ws
			return s.save()
		}
	}
	return ErrWorkspaceNotFound
}

// Delete removes a workspace
func (s *Store) Delete(name string) error {
	for i := range s.workspaces {
		if s.workspaces[i].Name == name {
			s.workspaces = append(s.workspaces[:i], s.workspaces[i+1:]...)
			return s.save()
		}
	}
	return ErrWorkspaceNotFound
}

// Filter returns workspaces matching the query
func (s *Store) Filter(query string) []Workspace {
	if query == "" {
		return s.workspaces
	}

	query = strings.ToLower(query)
	var filtered []Workspace
	for _, ws := range s.workspaces {
		if strings.Contains(strings.ToLower(ws.Name), query) ||
			strings.Contains(strings.ToLower(ws.Path), query) {
			filtered = append(filtered, ws)
		}
	}
	return filtered
}

// QuickAccess returns workspaces marked as quick access
func (s *Store) QuickAccess() []Workspace {
	var result []Workspace
	for _, ws := range s.workspaces {
		if ws.QuickAccess {
			result = append(result, ws)
		}
	}
	return result
}

// ToggleQuickAccess toggles quick access for a workspace
func (s *Store) ToggleQuickAccess(name string) error {
	for i := range s.workspaces {
		if s.workspaces[i].Name == name {
			s.workspaces[i].QuickAccess = !s.workspaces[i].QuickAccess
			return s.save()
		}
	}
	return ErrWorkspaceNotFound
}

// ListFolders returns all unique folder paths
func (s *Store) ListFolders() []string {
	folderSet := make(map[string]bool)
	for _, ws := range s.workspaces {
		if ws.Folder != "" {
			// Add folder and all parent folders
			parts := strings.Split(ws.Folder, "/")
			path := ""
			for _, part := range parts {
				if part == "" {
					continue
				}
				if path == "" {
					path = part
				} else {
					path = path + "/" + part
				}
				folderSet[path] = true
			}
		}
	}

	var folders []string
	for f := range folderSet {
		folders = append(folders, f)
	}
	return folders
}

// ListInFolder returns workspaces in a specific folder (not subfolders)
func (s *Store) ListInFolder(folder string) []Workspace {
	var result []Workspace
	for _, ws := range s.workspaces {
		if ws.Folder == folder {
			result = append(result, ws)
		}
	}
	return result
}

// ListSubfolders returns immediate subfolders of a folder
func (s *Store) ListSubfolders(folder string) []string {
	subfolderSet := make(map[string]bool)
	prefix := folder
	if prefix != "" {
		prefix = prefix + "/"
	}

	for _, ws := range s.workspaces {
		if ws.Folder == "" {
			continue
		}
		wsFolder := ws.Folder
		if prefix == "" {
			// Root level - get first part of folder
			parts := strings.SplitN(wsFolder, "/", 2)
			if len(parts) > 0 && parts[0] != "" {
				subfolderSet[parts[0]] = true
			}
		} else if strings.HasPrefix(wsFolder, prefix) {
			// Get next level
			rest := strings.TrimPrefix(wsFolder, prefix)
			parts := strings.SplitN(rest, "/", 2)
			if len(parts) > 0 && parts[0] != "" {
				subfolderSet[parts[0]] = true
			}
		}
	}

	var subfolders []string
	for f := range subfolderSet {
		subfolders = append(subfolders, f)
	}
	return subfolders
}

// ListRootWorkspaces returns workspaces not in any folder
func (s *Store) ListRootWorkspaces() []Workspace {
	var result []Workspace
	for _, ws := range s.workspaces {
		if ws.Folder == "" {
			result = append(result, ws)
		}
	}
	return result
}
