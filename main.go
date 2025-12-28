package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var dryRun bool

var rootCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "Dotfiles Manager",
	Long:  `A lightweight tool to manage your dotfiles using symlinks`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dotfiles structure",
	RunE:  runInit,
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync dotfiles to home directory",
	RunE:  runSync,
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate existing dotfiles to ~/.dotfiles",
	RunE:  runMigration,
}

func init() {
	initCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Log file operations instead of executing them")
	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Log file operations instead of executing them")
	migrateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Log file operations instead of executing them")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(migrateCmd)
}

func runSync(cmd *cobra.Command, args []string) error {
	df, err := New()
	if err != nil {
		return err
	}

	if !df.isInitialised() {
		return fmt.Errorf("dotfiles not initialized (run 'dotfiles init' first)")
	}

	df.DryRun = dryRun

	if !dryRun {
		if !confirmAction("This will SYNC your dotfiles (NOT a dry run). Continue?") {
			return nil
		}
	}

	linker, err := NewLinker(df)
	if err != nil {
		return err
	}

	fmt.Println("Syncing dotfiles...")
	if err := linker.Sync(); err != nil {
		return err
	}

	if !dryRun {
		fmt.Println("\n✅ Sync complete!")
	}

	return nil
}

func runMigration(cmd *cobra.Command, args []string) error {
	df, err := New()
	if err != nil {
		return err
	}

	if !df.isInitialised() {
		return fmt.Errorf("dotfiles not initialized (run 'dotfiles init' first)")
	}

	df.DryRun = dryRun

	if !dryRun {
		if !confirmAction("This will MIGRATE your dotfiles (NOT a dry run). Continue?") {
			return nil
		}
	}

	migrator, err := NewMigrator(df)
	if err != nil {
		return err
	}

	fmt.Println("Migrating declared config files to ~/.dotfiles ...")
	if err := migrator.Migrate(); err != nil {
		return err
	}

	if !dryRun {
		fmt.Println("\n✅ Migration complete!")
	}

	return nil
}

func runInit(cmd *cobra.Command, args []string) error {
	df, err := New()
	if err != nil {
		return err
	}
	df.DryRun = dryRun

	fmt.Println("Initializing dotfiles...")

	if err := df.Init(); err != nil {
		return err
	}

	fmt.Println("✅ Created", df.Dir)
	fmt.Println("✅ Created dotfiles.yaml")
	fmt.Println("✅ Created migrate.yaml")
	fmt.Println("✅ Created .gitignore")
	fmt.Println("✅ Created README.md")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. Edit", df.Dir+"/migrate.yaml", "with your files")
	fmt.Println("  2. dotfiles migrate")
	fmt.Println("  3. cd", df.Dir)
	if _, err := exec.LookPath("git"); err == nil {
		fmt.Println("  4. git add . && git commit -m 'Initial dotfiles'")
	} else {
		fmt.Println("  4. git init && git add . && git commit -m 'Initial dotfiles'")
	}

	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func confirmAction(message string) bool {
	fmt.Printf("%s [y/N]: ", message)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func main() {
	Execute()
}
