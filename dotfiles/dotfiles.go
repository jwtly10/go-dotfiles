package dotfiles

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed templates/README.md
var readmeTemplate []byte

const (
	FilePermReadWriteUser   = 0644
	FilePermReadWriteExUser = 0755
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

// IsInitialised checks if dir exists, and has a valid dotfiles.yaml marker
func (d *Dotfiles) IsInitialised() bool {
	markerPath := filepath.Join(d.Dir, ConfigFile)
	_, err := os.Stat(markerPath)
	return err == nil
}

func (d *Dotfiles) Init() error {
	if d.IsInitialised() {
		return fmt.Errorf("dotfiles already initialised in %s", d.Dir)
	}

	// Safety check: if .dotfiles exists, check if it's a file (not directory) or has existing content
	if info, err := os.Stat(d.Dir); err == nil {
		if !info.IsDir() {
			return fmt.Errorf("safety check failed: %s exists but is not a directory (it's a file). Please remove it first", d.Dir)
		}
		// Directory exists but not initialized - check if it has content
		entries, err := os.ReadDir(d.Dir)
		if err != nil {
			return fmt.Errorf("failed to read existing directory %s: %w", d.Dir, err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("safety check failed: %s already exists and contains %d item(s). This may be an existing dotfiles directory. If you want to reinitialize, please remove it first or ensure it's empty", d.Dir, len(entries))
		}
	}

	if d.DryRun {
		fmt.Println("[DRY RUN] Would create directory:", d.Dir)
		fmt.Println("[DRY RUN] Would create file:", filepath.Join(d.Dir, ConfigFile))
		fmt.Println("[DRY RUN] Would create file:", filepath.Join(d.Dir, MigrationFile))
		fmt.Println("[DRY RUN] Would create file:", filepath.Join(d.Dir, ".gitignore"))
		fmt.Println("[DRY RUN] Would create file:", filepath.Join(d.Dir, "README.md"))
		if _, err := exec.LookPath("git"); err == nil {
			fmt.Println("[DRY RUN] Would run: git init")
		}
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

	readmeFile := filepath.Join(d.Dir, "README.md")
	if err := os.WriteFile(readmeFile, readmeTemplate, FilePermReadWriteUser); err != nil {
		return fmt.Errorf("failed to create README.md: %w", err)
	}

	if gitPath, err := exec.LookPath("git"); err == nil {
		cmd := exec.Command(gitPath, "init")
		cmd.Dir = d.Dir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run git init: %w", err)
		}
		fmt.Println("✅ Initialized git repository")
	}

	return nil
}

// HomePath returns the absolute path in home dir
func (d *Dotfiles) HomePath(relativePath string) string {
	return filepath.Join(d.HomeDir, relativePath)
}

// DotPath returns the absolute path in dotfiles dir
func (d *Dotfiles) DotPath(relativePath string) string {
	return filepath.Join(d.Dir, relativePath)
}

func defaultGitignore() string {
	return fmt.Sprintf(`# go-dotfiles manager files
%s

# OS files
.DS_Store
*.log
`, MigrationFile)
}
