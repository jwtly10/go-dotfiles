package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	FilePermReadWriteUser = 0644
)

type Dotfiles struct {
	Dir     string // eg ~/.dotfiles
	HomeDir string // eg. ~
	DryRun  bool
}

func New() (*Dotfiles, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(home, ".dotfiles")
	return &Dotfiles{Dir: dir, HomeDir: home}, nil
}

// isInitialised checks if dir exists, and has a valid dotfiles.yaml marker
func (d *Dotfiles) isInitialised() bool {
	markerPath := filepath.Join(d.Dir, ConfigFile)
	_, err := os.Stat(markerPath)
	return err == nil
}

func (d *Dotfiles) Init() error {
	if d.isInitialised() {
		return fmt.Errorf("dotfiles already initialised in %s", d.Dir)
	}

	if d.DryRun {
		fmt.Println("[DRY RUN] Would create directory:", d.Dir)
		fmt.Println("[DRY RUN] Would create file:", filepath.Join(d.Dir, ConfigFile))
		fmt.Println("[DRY RUN] Would create file:", filepath.Join(d.Dir, MigrationFile))
		fmt.Println("[DRY RUN] Would create file:", filepath.Join(d.Dir, ".gitignore"))
		return nil
	}

	if err := os.MkdirAll(d.Dir, 0755); err != nil {
		return fmt.Errorf("failed to create dotfiles dir: %w", err)
	}

	err := DefaultDotfilesConfig().Save(d.Dir)
	if err != nil {
		return fmt.Errorf("failed to save default dotfiles.yaml: %w", err)
	}

	err = DefaultMigrateConfig().Save(d.Dir)
	if err != nil {
		return fmt.Errorf("failed to save default migrate.yaml: %w", err)
	}

	gitIgnoreFile := filepath.Join(d.Dir, ".gitignore")
	if err := os.WriteFile(gitIgnoreFile, []byte(defaultGitignore()), FilePermReadWriteUser); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	return nil
}

// HomePath returns the absolute path in home dir
func (d *Dotfiles) HomePath(relativePath string) string {
	return filepath.Join(d.HomeDir, relativePath)
}

func defaultGitignore() string {
	return fmt.Sprintf(`# go-dotfiles manager files
%s

# OS files
.DS_Store
*.log
`, MigrationFile)
}
