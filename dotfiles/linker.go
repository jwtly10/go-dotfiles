package dotfiles

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Linker struct {
	dotfiles *Dotfiles
	config   *DotfilesConfig
}

func NewLinker(dotfiles *Dotfiles) (*Linker, error) {
	config, err := LoadDotfilesConfig(dotfiles.Dir)
	if err != nil {
		return nil, fmt.Errorf("failed to load dotfiles config: %w", err)
	}
	return &Linker{dotfiles: dotfiles, config: config}, nil
}

// Sync walks files in ~/.dotfiles and creates symlinks in ~
func (l *Linker) Sync() error {
	return filepath.Walk(l.dotfiles.Dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Don't do anything with the dotfiles dir itself
		if path == l.dotfiles.Dir {
			return nil
		}

		relPath, err := filepath.Rel(l.dotfiles.Dir, path)
		if err != nil {
			return err
		}

		// Skip any ignored files in the config
		if l.shouldIgnore(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			fmt.Printf("⏭️ Skipping ignored file: %s\n", relPath)
			return nil
		}

		targetPath := l.dotfiles.HomePath(relPath)

		return l.createSymlink(path, targetPath)
	})
}

func (l *Linker) createSymlink(source, target string) error {
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	// Dirs will be created on file symlink
	if sourceInfo.IsDir() {
		return nil
	}

	// Safety check: warn about critical files
	if IsCriticalFile(target) {
		fmt.Printf("⚠️  WARNING: %s is a critical system file. Proceed with caution.\n", target)
	}

	targetInfo, err := os.Lstat(target)
	if err == nil { // Target already exists
		if targetInfo.Mode()&os.ModeSymlink != 0 {
			// Existing symlink
			existing, _ := os.Readlink(target)
			if existing == source {
				return nil
			}

			fmt.Printf("⚠️ %s is an external symlink to %s (expected %s)\n", target, existing, source)
			if l.dotfiles.DryRun {
				fmt.Println("[DRY RUN] Would request resolution to conflicting external symlink:", target)
				return nil
			}
			return l.handleConflict(source, target, "symlink exists to different location")
		}

		// Safety: for critical files, provide more context
		if IsCriticalFile(target) {
			fmt.Printf("⚠️  CRITICAL: %s exists and is not a symlink. This is a critical system file.\n", target)
			fmt.Printf("   Consider backing up before proceeding: dotfiles will need to replace this file.\n")
		}

		if l.dotfiles.DryRun {
			fmt.Println("[DRY RUN] Would request resolution to conflicting existing file:", target)
			return nil
		}
		return l.handleConflict(source, target, "file exists and is not a symlink")
	}

	if l.dotfiles.DryRun {
		fmt.Printf("[DRY RUN] Would create symlink: %s -> %s\n", target, source)
		return nil
	}

	targetDir := filepath.Dir(target)
	if err := os.MkdirAll(targetDir, FilePermReadWriteExUser); err != nil {
		return fmt.Errorf("failed to create parent dir: %w", err)
	}

	if err := os.Symlink(source, target); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	fmt.Printf("✅ Linked: %s\n", target)
	return nil
}

func (l *Linker) shouldIgnore(relPath string) bool {
	// Ignore some of the system generated files when linking
	if relPath == ".git" || relPath == ".gitignore" || relPath == "dotfiles.yaml" || relPath == "migrate.yaml" || relPath == "README.md" {
		return true
	}

	ignoreList := l.config.Ignore
	return shouldIgnore(ignoreList, relPath)
}

func (l *Linker) handleConflict(source, target, reason string) error {
	// Provide detailed error message with suggestions
	errMsg := fmt.Sprintf("failed to link %s -> %s: %s", target, source, reason)
	
	// Add helpful suggestions
	if IsCriticalFile(target) {
		errMsg += fmt.Sprintf("\n\n⚠️  SAFETY: %s is a critical system file.", target)
		errMsg += "\n   Suggestions:"
		errMsg += "\n   1. Backup the file manually before proceeding"
		errMsg += "\n   2. Review the file contents to ensure nothing important will be lost"
		errMsg += "\n   3. Consider using '--dry-run' first to preview changes"
	} else {
		errMsg += "\n\n   Suggestions:"
		errMsg += "\n   1. Backup or rename the existing file"
		errMsg += "\n   2. Remove the file if it's safe to do so"
		errMsg += "\n   3. Use '--dry-run' first to preview changes"
	}
	
	return fmt.Errorf("%s", errMsg)
}
