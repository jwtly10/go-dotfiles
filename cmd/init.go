package cmd

import (
	"fmt"
	"os/exec"

	"github.com/jwtly10/go-dotfiles/dotfiles"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize dotfiles structure",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Log file operations instead of executing them")
}

func runInit(cmd *cobra.Command, args []string) error {
	df, err := dotfiles.New()
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
