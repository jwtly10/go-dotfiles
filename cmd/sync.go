package cmd

import (
	"fmt"

	"github.com/jwtly10/go-dotfiles/dotfiles"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync dotfiles to home directory",
	RunE:  runSync,
}

func init() {
	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Log file operations instead of executing them")
}

func runSync(cmd *cobra.Command, args []string) error {
	df, err := dotfiles.New()
	if err != nil {
		return err
	}

	if !df.IsInitialised() {
		return fmt.Errorf("dotfiles not initialized (run 'dotfiles init' first)")
	}

	df.DryRun = dryRun

	if !dryRun {
		if !confirmAction("This will SYNC your dotfiles (NOT a dry run). Continue?") {
			return nil
		}
	}

	linker, err := dotfiles.NewLinker(df)
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
