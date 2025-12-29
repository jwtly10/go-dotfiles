# go-dotfiles

A minimal dotfiles manager written in Go. It uses symlinks to keep configuration files organized in a single directory while keeping them in their expected locations.

- **Built in migration**: Move existing config files from your home directory into your dotfiles repo.
- **Symlink Management**: Creates symlinks for all enabled files in your user home directory.
- **Dry Run Support**: Preview all operations before they happen.
- **Ignore Patterns**: Support for glob patterns to ignore specific files or directories.

## Quick Start

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

> **⚠️ It's recommended to backup your important files _before_ migrating them with `go-dotfiles`.** And to never add credentials to the dotfiles directory.

Add paths to `migrate.yaml` and run:

```bash
# Preview changes first
go-dotfiles migrate --dry-run

# Commit changes
go-dotfiles migrate
```

Example `migrate.yaml`:
```yaml
paths:
  - .zshrc
  - .config/nvim
  - .config/starship.toml
  - .config/.gitconfig

ignore:
  - .config/.gitconfig
```

This will then create the following dotfiles structure, and move the original files to the dotfiles directory:
```
~/.dotfiles/
├── .config
│   ├── .gitconfig
│   ├── starship.toml
│   └── nvim
│       ├── init.vim
│       ├── init.lua
│       ├── ... 
│       │   └── ... other nested files
├── .zshrc
```
### 4. Sync Symlinks

You can now run 

```bash
# Preview changes first
go-dotfiles sync --dry-run

# Commit changes
go-dotfiles sync
```

which will create the following symlink structure:

```
~/.zshrc -> ~/.dotfiles/.zshrc
~/.config/nvim -> ~/.dotfiles/.config/nvim
~/.config/starship.toml -> ~/.dotfiles/.config/starship.toml
~/.config/.gitconfig -> ~/.dotfiles/.config/.gitconfig
... etc
```

### 5. Syncing to a New Machine

On a new machine, with go-dotfiles installed, clone your dotfiles repo to `~/.dotfiles` and run:

```bash
go-dotfiles sync
```

This will generate all the symlinks for the files in your dotfiles directory to your home directory.

## 🛠️ Commands

- `init`: Setup the initial configuration.
- `migrate`: Move files defined in `migrate.yaml` into the dotfiles directory
- `sync`: Symlink files from the dotfiles directory to the home directory.
- `version`: Show the commit hash this binary was built from.

You can also create a shell alias for git - allowing you to quickly commit changes from any working dir

``` bash
alias dfg="git -C ~/.dotfiles"
```


## ⚙️ Configuration

### `dotfiles.yaml`

Used to define patterns that should be ignored during sync between the dotfiles directory and your home directory.

```yaml
ignore:
  - .DS_Store
  - "*.log"
  - ".git"
```


### `migrate.yaml`

Used to define files that should be migrated from your home directory to your dotfiles directory, and any patterns that should be ignored during the migration process.

```yaml
paths:
  - .zshrc
  - .config/nvim
  - .config/starship.toml
  - .config/.gitconfig

ignore:
  - .config/.gitconfig

``` 