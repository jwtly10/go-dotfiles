# .dotfiles

Personal configuration files, managed by [go-dotfiles](https://github.com/jwtly10/go-dotfiles).

##  Getting Started

If you are setting this up on a new machine:

1. Clone this repository to `~/.dotfiles`.
2. Install `go-dotfiles`.
3. Run the sync command:

```bash
go-dotfiles sync
```

This generates symlinks to the files in `~/.dotfiles/` to your home directory. Configuration files should be managed directly with the dotfiles directory, and `go-dotfiles sync` should be run to update your home directory with the latest changes.

## 📦 Adding new configuration

Any new dotfile should be added to `~/.dotfiles/` as a mirror of `~`

Then you can re-run

```bash
# To preview changes
go-dotfiles sync --dry-run

# To commit changes
go-dotfiles sync
```

## ⚙️ Configuration

- `dotfiles.yaml`: Define patterns that should be ignored during sync between the dotfiles directory and your home directory.
- `migrate.yaml`: Define files to be migrated from your home directory to this repository. 

## 🛠️ Commands

- `go-dotfiles sync`: Creates symlinks for all files in this repository to your home directory.
- `go-dotfiles migrate`: Moves files listed in `migrate.yaml` into this repository in preparation for syncing.