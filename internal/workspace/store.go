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
