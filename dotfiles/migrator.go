package dotfiles

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MaxMigrations is a hard limit on the number of files we will migrate
// in a given run.
//
// While we can't always cater for user input errors, this will help to prevent deeply nested
// directories seemingly being parsed infinitely by mistake
//
// Rerunning the command will allow another batch to be migrated
const MaxMigrations = 100

// DryRunMaxMigrations is the hard limit used when running a dry run
//
// Read operations are less of a concern here, so the limit is much larger to provide more details
// if possible - e.g., 300 changes, but only 200 will be applied on next migration
const DryRunMaxMigrations = 1000

type Migrator struct {
	dotfiles *Dotfiles
	config   *MigrateConfig

	migrated int
	existing int
	errors   []string
}

func NewMigrator(dotfiles *Dotfiles) (*Migrator, error) {
	config, err := LoadMigrationConfig(dotfiles.Dir)
	if err != nil {
		return nil, err
	}
	return &Migrator{dotfiles: dotfiles, config: config}, nil
}

// Migrate moves request dotfiles to ~/.dotfiles.
//
// Note this method does NOT then re-sync symlinks, so existing configurations
// will be 'missing' from original locations, until `go-dotfiles` sync is run
//
// It's recommended to back up files manually first
func (m *Migrator) Migrate() error {
	fmt.Printf("Attempting to migrate %d paths\n", len(m.config.Paths))
	for _, path := range m.config.Paths {
		if m.dotfiles.DryRun {
			if m.migrated >= DryRunMaxMigrations {
				fmt.Printf("Reached dry run migration limit of %d files.\n", DryRunMaxMigrations)
				break
			}
		} else {
			if m.migrated >= MaxMigrations {
				fmt.Printf("Reached migration limit of %d files. Run again to migrate more.\n", MaxMigrations)
				break
			}
		}

		source := m.dotfiles.HomePath(path)

		// In case the user has a path in both - ignore takes priority
		if m.shouldIgnore(path) {
			fmt.Printf("⏭️ Skipping ignored file: %s\n", path)
			continue
		}

		sourceInfo, err := os.Lstat(source)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Skipping %s (does not exist)\n", source)
				continue
			}
			return fmt.Errorf("failed to stat %s: %w", source, err)
		}

		if sourceInfo.IsDir() {
			err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if m.dotfiles.DryRun {
					if m.migrated >= DryRunMaxMigrations {
						fmt.Printf("Reached dry run migration limit of %d files.\n", DryRunMaxMigrations)
						return filepath.SkipAll
					}
				} else {
					if m.migrated >= MaxMigrations {
						fmt.Printf("Reached migration limit of %d files. Run again to migrate more.\n", MaxMigrations)
						return filepath.SkipAll
					}
				}

				relPath, err := filepath.Rel(m.dotfiles.HomeDir, path)
				if err != nil {
					return err
				}
				if m.shouldIgnore(relPath) {
					if info.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
				if info.IsDir() {
					return nil
				}

				if err := m.migrateFile(path, m.dotfiles.DotPath(relPath)); err != nil {
					m.errors = append(m.errors, fmt.Sprintf("failed to migrate %s: %s", path, err))
				}
				return nil
			})

			if err != nil {
				m.errors = append(m.errors, fmt.Sprintf("failed to walk dir %s: %s", source, err))
			}
		} else {
			target := m.dotfiles.DotPath(path)
			if err := m.migrateFile(source, target); err != nil {
				m.errors = append(m.errors, fmt.Sprintf("failed to migrate %s: %s", source, err))
			}
		}
	}

	if m.dotfiles.DryRun {
		fmt.Printf("\n[DRY RUN SUMMARY]\n")
		fmt.Printf("Checked %d files for migration.\n", m.migrated)
		if m.migrated > MaxMigrations {
			fmt.Printf("NOTICE: On an actual migration, only the first %d files would be moved (MaxMigrations limit).\n", MaxMigrations)
			fmt.Printf("You would need to run the command multiple times to migrate all %d files.\n", m.migrated)
		} else {
			fmt.Printf("On an actual migration, %d files would be moved.\n", m.migrated)
		}
		fmt.Println()
	}

	fmt.Printf("Migration complete: %d files migrated, %d files already migrated, %d file errors\n", m.migrated, m.existing, len(m.errors))
	if len(m.errors) > 0 {
		return fmt.Errorf("errors:\n%s", strings.Join(m.errors, "\n"))
	}

	return nil
}

// migrateFile moves a single file from source to target
func (m *Migrator) migrateFile(absSource, absTarget string) error {
	// Safety check: warn about critical files
	if IsCriticalFile(absSource) {
		fmt.Printf("⚠️  WARNING: Migrating critical system file: %s\n", absSource)
	}

	// Safety check: validate path
	if err := CheckPathSafety(absSource); err != nil {
		return fmt.Errorf("safety check failed for %s: %w", absSource, err)
	}

	if s, err := os.Lstat(absSource); err == nil {
		if s.IsDir() {
			return m.handleConflict(absSource, absTarget, "can't migrate a directory")
		}
		if s.Mode()&os.ModeSymlink != 0 {
			existing, _ := os.Readlink(absSource)
			if existing == absTarget {
				// It's a dotfile symlink, mark as found existing
				if m.dotfiles.DryRun {
					// Logging just for debugging purposes
					fmt.Printf("[DRY RUN] Skipping existing dotfile symlink: %s -> %s\n", absSource, absTarget)
				}
				m.existing++
				return nil
			}
			return m.handleConflict(absSource, absTarget, "source is a symlink")
		}
	} else {
		if os.IsNotExist(err) {
			return fmt.Errorf("skipping %s (does not exist)", absSource)
		}
		return fmt.Errorf("failed to stat %s: %w", absSource, err)
	}

	if _, err := os.Stat(absTarget); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat %s: %w", absTarget, err)
		}
	} else {
		return m.handleConflict(absSource, absTarget, "already exists in dotfiles")
	}

	if m.dotfiles.DryRun {
		fmt.Printf("[DRY RUN] Would migrate %s -> %s\n", absSource, absTarget)
		if IsCriticalFile(absSource) {
			fmt.Printf("[DRY RUN] ⚠️  This is a critical system file - ensure you have backups!\n")
		}
		m.migrated++
		return nil
	}

	// Safety: Create backup for critical files before migration
	if IsCriticalFile(absSource) {
		backupPath, err := BackupFile(absSource)
		if err != nil {
			return fmt.Errorf("failed to create backup of critical file %s: %w", absSource, err)
		}
		fmt.Printf("📦 Created backup: %s\n", backupPath)
	}

	targetDir := filepath.Dir(absTarget)
	if err := os.MkdirAll(targetDir, FilePermReadWriteExUser); err != nil {
		return err
	}

	fmt.Printf("Migrating %s -> %s\n", absSource, absTarget)
	err := os.Rename(absSource, absTarget)
	if err != nil {
		return fmt.Errorf("failed to rename file %s: %w", absSource, err)
	}

	m.migrated++
	return nil
}

func (m *Migrator) handleConflict(source, target, reason string) error {
	// Provide detailed error message with suggestions
	errMsg := fmt.Sprintf("failed to migrate %s -> %s: %s", source, target, reason)
	
	// Add helpful suggestions
	if IsCriticalFile(source) {
		errMsg += fmt.Sprintf("\n\n⚠️  SAFETY: %s is a critical system file.", source)
		errMsg += "\n   Suggestions:"
		errMsg += "\n   1. Ensure you have a backup of this file"
		errMsg += "\n   2. Review the file contents before migration"
		errMsg += "\n   3. Use '--dry-run' first to preview changes"
	} else {
		errMsg += "\n\n   Suggestions:"
		errMsg += "\n   1. Check if the file should be migrated"
		errMsg += "\n   2. Remove or rename conflicting files if safe"
		errMsg += "\n   3. Use '--dry-run' first to preview changes"
	}
	
	return fmt.Errorf("%s", errMsg)
}

func (m *Migrator) shouldIgnore(relPath string) bool {
	return shouldIgnore(m.config.Ignore, relPath)
}

func shouldIgnore(ignoreList []string, relPath string) bool {
	base := filepath.Base(relPath)
	for _, pattern := range ignoreList {
		matched, err := filepath.Match(pattern, relPath)
		if err != nil {
			panic(fmt.Sprintf("bad pattern found: '%q' : %v", pattern, err))
		}
		if matched {
			return true
		}

		matched, err = filepath.Match(pattern, base)
		if err != nil {
			panic(fmt.Sprintf("bad pattern found: '%q' : %v", pattern, err))
		}
		if matched {
			return true
		}
	}
	return false
}
