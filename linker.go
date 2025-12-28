package main

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
	if relPath == ".git" || relPath == ".gitignore" || relPath == "dotfiles.yaml" || relPath == "migrate.yaml" {
		return true
	}

	ignoreList := l.config.Ignore
	return shouldIgnore(ignoreList, relPath)
}

func (l *Linker) handleConflict(source, target, reason string) error {
	// TODO: Let the user decide what to do here
	return fmt.Errorf("failed to link %s -> %s: %s", target, source, reason)
}
