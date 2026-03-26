package git

import (
	"bufio"
	"os/exec"
	"path/filepath"
	"strings"
)

// Worktree represents a git worktree
type Worktree struct {
	Path   string // Absolute path to worktree
	Branch string // Branch name (e.g., "feature/auth")
	Commit string // Current commit hash
	IsMain bool   // True if this is the main worktree
}

// IsGitRepo checks if the given path is inside a git repository
func IsGitRepo(path string) bool {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// GetRepoRoot returns the root directory of the git repository
func GetRepoRoot(path string) (string, error) {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// ListWorktrees returns all worktrees for a repository
func ListWorktrees(repoPath string) ([]Worktree, error) {
	cmd := exec.Command("git", "-C", repoPath, "worktree", "list", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var worktrees []Worktree
	var current Worktree
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = Worktree{}
			}
			continue
		}

		if strings.HasPrefix(line, "worktree ") {
			current.Path = strings.TrimPrefix(line, "worktree ")
		} else if strings.HasPrefix(line, "HEAD ") {
			current.Commit = strings.TrimPrefix(line, "HEAD ")
		} else if strings.HasPrefix(line, "branch ") {
			// Format: "branch refs/heads/main"
			branch := strings.TrimPrefix(line, "branch refs/heads/")
			current.Branch = branch
		} else if line == "bare" {
			current.IsMain = true
		}
	}

	// Add the last worktree if any
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	// Mark the first worktree as main (if not bare repo)
	if len(worktrees) > 0 && !worktrees[0].IsMain {
		worktrees[0].IsMain = true
	}

	return worktrees, nil
}

// CreateWorktree creates a new worktree for the given branch
// Location: <repoPath>/.worktrees/<sanitized-branch-name>
func CreateWorktree(repoPath, branch string) (Worktree, error) {
	sanitized := SanitizeBranchName(branch)
	worktreePath := filepath.Join(repoPath, ".worktrees", sanitized)

	// Check if branch exists locally or remotely
	localExists := branchExists(repoPath, branch, false)
	remoteExists := branchExists(repoPath, branch, true)

	var cmd *exec.Cmd
	if localExists {
		// Branch exists locally, just create worktree
		cmd = exec.Command("git", "-C", repoPath, "worktree", "add", worktreePath, branch)
	} else if remoteExists {
		// Branch exists on remote, create tracking branch
		cmd = exec.Command("git", "-C", repoPath, "worktree", "add", worktreePath, "-b", branch, "origin/"+branch)
	} else {
		// New branch, create from current HEAD
		cmd = exec.Command("git", "-C", repoPath, "worktree", "add", "-b", branch, worktreePath)
	}

	if err := cmd.Run(); err != nil {
		return Worktree{}, err
	}

	// Get commit hash
	commitCmd := exec.Command("git", "-C", worktreePath, "rev-parse", "HEAD")
	out, _ := commitCmd.Output()
	commit := strings.TrimSpace(string(out))

	return Worktree{
		Path:   worktreePath,
		Branch: branch,
		Commit: commit,
		IsMain: false,
	}, nil
}

// RemoveWorktree removes a worktree
func RemoveWorktree(repoPath, worktreePath string) error {
	// First try normal remove
	cmd := exec.Command("git", "-C", repoPath, "worktree", "remove", worktreePath)
	if err := cmd.Run(); err != nil {
		// Try force remove if normal fails
		cmd = exec.Command("git", "-C", repoPath, "worktree", "remove", "--force", worktreePath)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	// Prune stale worktree info
	pruneCmd := exec.Command("git", "-C", repoPath, "worktree", "prune")
	return pruneCmd.Run()
}

// ListBranches returns all branches (local + remote)
func ListBranches(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoPath, "branch", "-a", "--format=%(refname:short)")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var branches []string
	seen := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		branch := scanner.Text()
		// Remove "origin/" prefix for remote branches
		if strings.HasPrefix(branch, "origin/") {
			branch = strings.TrimPrefix(branch, "origin/")
			// Skip HEAD pointer
			if branch == "HEAD" {
				continue
			}
		}
		// Deduplicate (same branch may exist locally and remotely)
		if !seen[branch] {
			seen[branch] = true
			branches = append(branches, branch)
		}
	}

	return branches, nil
}

// SanitizeBranchName converts branch name to filesystem-safe name
// e.g., "feature/auth" -> "feature-auth"
func SanitizeBranchName(branch string) string {
	// Replace common path separators
	result := strings.ReplaceAll(branch, "/", "-")
	result = strings.ReplaceAll(result, "\\", "-")
	// Remove other problematic characters
	result = strings.ReplaceAll(result, ":", "-")
	result = strings.ReplaceAll(result, "*", "-")
	result = strings.ReplaceAll(result, "?", "-")
	result = strings.ReplaceAll(result, "\"", "-")
	result = strings.ReplaceAll(result, "<", "-")
	result = strings.ReplaceAll(result, ">", "-")
	result = strings.ReplaceAll(result, "|", "-")
	// Collapse multiple dashes
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}

// branchExists checks if a branch exists locally or remotely
func branchExists(repoPath, branch string, remote bool) bool {
	var ref string
	if remote {
		ref = "refs/remotes/origin/" + branch
	} else {
		ref = "refs/heads/" + branch
	}
	cmd := exec.Command("git", "-C", repoPath, "show-ref", "--verify", "--quiet", ref)
	return cmd.Run() == nil
}
