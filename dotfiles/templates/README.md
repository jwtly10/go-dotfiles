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

## 📦 Adding new configuration

Any new dotfile should be added to `~/.dotfiles/` as a mirror of `~`

Then you can re-run

```bash
# To preview
go-dotfiles sync --dry-run
# To commit
go-dotfiles sync
```

To sync the file to the userspace, and due to the symlink, future changes will be automatically applied

## ⚙️ Configuration

- `dotfiles.yaml`: Define global ignore patterns.
- `migrate.yaml`: Define files to be migrated from your home directory to this repository.

## 🛠️ Commands

- `go-dotfiles sync`: Creates symlinks for all files in this repository to your home directory.
- `go-dotfiles migrate`: Moves files listed in `migrate.yaml` into this repository and symlinks them back.
