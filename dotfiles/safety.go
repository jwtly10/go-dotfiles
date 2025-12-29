package dotfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// CriticalFiles is a list of critical system files that should be handled with extra care
var CriticalFiles = []string{
	".bashrc",
	".bash_profile",
	".zshrc",
	".zprofile",
	".profile",
	".bash_history",
	".zsh_history",
	".ssh/config",
	".ssh/authorized_keys",
	".ssh/id_rsa",
	".ssh/id_ed25519",
	".gitconfig",
}

// IsCriticalFile checks if a file path matches any critical system files
func IsCriticalFile(path string) bool {
	base := filepath.Base(path)
	for _, critical := range CriticalFiles {
		if base == filepath.Base(critical) || path == critical {
			return true
		}
	}
	return false
}

// BackupFile creates a backup of a file with a timestamp
func BackupFile(filePath string) (string, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to stat file for backup: %w", err)
	}

	if info.IsDir() {
		return "", fmt.Errorf("cannot backup directory: %s", filePath)
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupPath := fmt.Sprintf("%s.backup_%s", filePath, timestamp)

	// Read original file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file for backup: %w", err)
	}

	// Write backup
	if err := os.WriteFile(backupPath, data, info.Mode()); err != nil {
		return "", fmt.Errorf("failed to write backup: %w", err)
	}

	return backupPath, nil
}

// CheckPathSafety performs safety checks on a path before operations
func CheckPathSafety(path string) error {
	// Check if path is absolute and outside home directory (extra safety)
	if filepath.IsAbs(path) {
		// This is a basic check - in practice, paths should be relative to home
		if path == "/" || path == "/root" || path == "/etc" {
			return fmt.Errorf("safety check failed: path %s appears to be a system directory", path)
		}
	}
	return nil
}

