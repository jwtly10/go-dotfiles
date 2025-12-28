# go-dotfiles

A minimal dotfiles manager written in Go. It uses symlinks to keep your configuration files organized in a single directory while keeping them in their expected locations.

- **Built in migration**: Move existing config files from your home directory into your dotfiles repo.
- **Symlink Management**: Automatically sync files from your repository to your home directory.
- **Dry Run Support**: Preview all operations before they happen.
- **Ignore Patterns**: Support for glob patterns to ignore specific files or directories.

## 🚀 Quick Start

### 1. Installation

```bash
go install github.com/jwtly10/go-dotfiles@latest
```

### 2. Initialization

Initialize your dotfiles directory (defaults to `~/.dotfiles`):

```bash
go-dotfiles init
```

This creates the directory and some initial configuration files:
- `dotfiles.yaml`: General settings and global ignore list.
- `migrate.yaml`: List of files you want to move into your dotfiles.
- `.gitignore`: Basic ignore rules for your dotfile repository.
- `README.md`: A template for your new dotfiles repo.

### 3. Migrating Files

Add paths to `migrate.yaml` and run:

```bash
# Preview changes
go-dotfiles migrate --dry-run
# Commit changes
go-dotfiles migrate
```

Example `migrate.yaml`:
```yaml
paths:
  - .zshrc
  - .config/nvim
  - .gitconfig
```

### 4. Syncing to a New Machine

On a new machine, clone your dotfiles repo to `~/.dotfiles` and run:

```bash
go-dotfiles sync
```

## 🛠️ Commands

- `init`: Setup the initial structure.
- `sync`: Symlink files from the dotfiles directory to the home directory.
- `migrate`: Move files defined in `migrate.yaml` into the dotfiles directory and symlink them back.

## ⚙️ Configuration

### `dotfiles.yaml`

Used to define patterns that should be ignored during sync.

```yaml
ignore:
  - .DS_Store
  - "*.log"
  - ".git"
```