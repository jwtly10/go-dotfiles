package cmd

import (
	"fmt"

	"github.com/jwtly10/go-dotfiles/dotfiles"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate existing dotfiles to ~/.dotfiles",
	RunE:  runMigration,
}

func init() {
	migrateCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Log file operations instead of executing them")
}

func runMigration(cmd *cobra.Command, args []string) error {
	df, err := dotfiles.New()
	if err != nil {
		return err
	}

	if !df.IsInitialised() {
		return fmt.Errorf("dotfiles not initialized (run 'dotfiles init' first)")
	}

	df.DryRun = dryRun

	if !dryRun {
		if !confirmAction("This will MIGRATE your dotfiles (NOT a dry run). Continue?") {
			return nil
		}
	}

	migrator, err := dotfiles.NewMigrator(df)
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
