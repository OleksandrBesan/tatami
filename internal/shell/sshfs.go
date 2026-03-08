package shell

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// SSHFSRunner handles SSHFS operations
type SSHFSRunner struct {
	mountBase string // Base directory for mounts
}

// NewSSHFSRunner creates a new SSHFS runner
func NewSSHFSRunner() *SSHFSRunner {
	home, _ := os.UserHomeDir()
	return &SSHFSRunner{
		mountBase: filepath.Join(home, ".local", "share", "tatami", "mounts"),
	}
}

// IsInstalled checks if sshfs is available
func (s *SSHFSRunner) IsInstalled() bool {
	_, err := exec.LookPath("sshfs")
	return err == nil
}

// InstallInstructions returns OS-specific install instructions
func (s *SSHFSRunner) InstallInstructions() string {
	var b strings.Builder
	b.WriteString("SSHFS is required for remote workspaces.\n\n")

	switch runtime.GOOS {
	case "darwin":
		b.WriteString("Install on macOS:\n")
		b.WriteString("  brew install macfuse\n")
		b.WriteString("  brew install sshfs\n")
		b.WriteString("\n")
		b.WriteString("Note: You may need to allow the kernel extension\n")
		b.WriteString("in System Preferences > Security & Privacy.\n")
	case "linux":
		b.WriteString("Install on Linux:\n")
		b.WriteString("  Ubuntu/Debian: sudo apt install sshfs\n")
		b.WriteString("  Fedora/RHEL:   sudo dnf install fuse-sshfs\n")
		b.WriteString("  Arch:          sudo pacman -S sshfs\n")
	default:
		b.WriteString("Install sshfs for your operating system.\n")
	}

	b.WriteString("\nAfter installing, run tatami again.")
	return b.String()
}

// GetMountPoint returns the mount point for a remote workspace
func (s *SSHFSRunner) GetMountPoint(host, remotePath string) string {
	// Create a safe directory name from host and path
	safeName := strings.ReplaceAll(host, "@", "_at_")
	safeName = strings.ReplaceAll(safeName, ":", "_")
	safeName = strings.ReplaceAll(safeName, "/", "_")
	pathSafe := strings.ReplaceAll(remotePath, "/", "_")
	return filepath.Join(s.mountBase, safeName+pathSafe)
}

// IsMounted checks if a path is already mounted
func (s *SSHFSRunner) IsMounted(mountPoint string) bool {
	// Check if mount point exists and has files
	entries, err := os.ReadDir(mountPoint)
	if err != nil {
		return false
	}
	return len(entries) > 0
}

// Mount mounts a remote path via SSHFS
func (s *SSHFSRunner) Mount(host, remotePath string) (string, error) {
	mountPoint := s.GetMountPoint(host, remotePath)

	// Check if already mounted
	if s.IsMounted(mountPoint) {
		return mountPoint, nil
	}

	// Create mount point directory
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return "", fmt.Errorf("failed to create mount point: %w", err)
	}

	// Mount via sshfs
	remote := fmt.Sprintf("%s:%s", host, remotePath)
	cmd := exec.Command("sshfs", remote, mountPoint,
		"-o", "reconnect",
		"-o", "ServerAliveInterval=15",
		"-o", "ServerAliveCountMax=3",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		// Clean up empty mount point on failure
		os.Remove(mountPoint)
		return "", fmt.Errorf("sshfs mount failed: %w", err)
	}

	return mountPoint, nil
}

// Unmount unmounts a SSHFS mount point
func (s *SSHFSRunner) Unmount(mountPoint string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("umount", mountPoint)
	} else {
		cmd = exec.Command("fusermount", "-u", mountPoint)
	}
	return cmd.Run()
}

// UnmountAll unmounts all tatami SSHFS mounts
func (s *SSHFSRunner) UnmountAll() error {
	entries, err := os.ReadDir(s.mountBase)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			mountPoint := filepath.Join(s.mountBase, entry.Name())
			s.Unmount(mountPoint)
		}
	}
	return nil
}
