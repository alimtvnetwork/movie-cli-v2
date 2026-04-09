# Movie CLI

A cross-platform CLI tool for managing a personal movie and TV show library. Scans local folders for video files, cleans messy filenames, fetches metadata from TMDb, stores everything in SQLite, and organizes files into configured directories.

## Features

- **Scan** local folders for video files and auto-fetch metadata from TMDb
- **Search** TMDb directly and save results to your library
- **List** your library with paginated, interactive browsing
- **Move** and **rename** files with clean naming (`Title (Year).ext`)
- **Undo** any move/rename operation
- **Play** media files with your system's default player
- **Suggest** new content based on your library's genre patterns or trending
- **Stats** with genre charts and average ratings
- **Self-update** via `git pull --ff-only`

## Commands

```
movie-cli
├── hello                         # Greeting with version
├── version                       # Version, commit, build date
├── self-update                   # git pull --ff-only
└── movie
    ├── config [get|set] [key]    # View/set configuration
    ├── scan [folder]             # Scan folder → DB + TMDb metadata
    ├── ls                        # Paginated library list (file-backed only)
    ├── search <name>             # Live TMDb search → save to DB
    ├── info <id|title>           # Detail view (local DB → TMDb fallback)
    ├── suggest [N]               # Recommendations + trending
    ├── move [directory]          # Browse, select, move with clean name
    ├── rename                    # Batch rename to clean format
    ├── undo                      # Revert last move/rename
    ├── play <id>                 # Open with default video player
    ├── stats                     # Counts, storage, genre chart, avg ratings
    └── tag [add|remove|list]     # Manage user-defined tags
```

## Installation

### Prerequisites

| Requirement | Minimum | Check |
|---|---|---|
| **Go** | 1.22+ | `go version` |
| **Git** | 2.x | `git --version` |
| **PowerShell** | 5.1+ (Win) / 7+ (Unix) | `$PSVersionTable.PSVersion` |

### One-Liner Install

**Windows (PowerShell)**
```powershell
git clone https://github.com/mahin/mahin-cli-v1.git; cd mahin-cli-v1; .\run.ps1
```

**macOS / Linux**
```bash
git clone https://github.com/mahin/mahin-cli-v1.git && cd mahin-cli-v1 && pwsh run.ps1
```

### Using the Installer

```powershell
# Fresh install (clones repo if needed, builds, deploys)
pwsh install.ps1

# Custom deploy path
pwsh install.ps1 -DeployPath ~/bin
```

### Verify

```bash
mahin version
# v1.x.x (commit: abc1234, built: 2024-06-01T12:00:00+08:00)
```

> **Tip**: If `mahin` is not found, add the deploy path to your `PATH`.  
> Default: `E:\bin-run` (Windows) or `/usr/local/bin` (Unix).

## Quick Start

```bash
# Set your TMDb API key
mahin movie config set tmdb_api_key YOUR_KEY

# Scan a folder
mahin movie scan ~/Downloads

# Browse your library
mahin movie ls

# Search TMDb directly
mahin movie search "Inception"

# Get suggestions
mahin movie suggest 5
```

## Build & Deploy (run.ps1)

The `run.ps1` script is the single-entry automation for pull → build → deploy → run.

```powershell
.\run.ps1                           # Full pipeline
.\run.ps1 -NoPull                   # Skip git pull
.\run.ps1 -NoPull -NoDeploy        # Build only
.\run.ps1 -R movie scan D:\movies  # Build + run scan
.\run.ps1 -t                       # Run all unit tests
.\run.ps1 -ForcePull               # CI mode: discard changes + pull
```

See [spec/03-general/04-run-guide.md](spec/03-general/04-run-guide.md) for the full usage guide.

---

## Command Reference

### `movie-cli hello`

Print a greeting with the current version.

```bash
movie-cli hello
# 👋 Hello from Movie CLI! v1.2.0
```

### `movie-cli version`

Show version, commit hash, and build date (injected via `-ldflags`).

```bash
movie-cli version
# v1.2.0 (commit: abc1234, built: 2024-06-01)
```

### `movie-cli self-update`

Pull latest code via `git pull --ff-only`. Requires a clean working tree.

```bash
movie-cli self-update
# ✅ Updated abc1234 → def5678
```

### `movie-cli movie config`

View or update configuration settings. API keys are masked in output.

```bash
movie-cli movie config                          # Show all settings
movie-cli movie config get movies_dir           # Get a single key
movie-cli movie config set movies_dir ~/Movies  # Set a key
movie-cli movie config set tmdb_api_key KEY     # Set API key
movie-cli movie config set page_size 30         # Items per page
```

| Key | Default | Purpose |
|---|---|---|
| `movies_dir` | `~/Movies` | Movie file destination |
| `tv_dir` | `~/TVShows` | TV show destination |
| `archive_dir` | `~/Archive` | Archive destination |
| `scan_dir` | `~/Downloads` | Default scan source |
| `tmdb_api_key` | *(none)* | TMDb API key |
| `page_size` | `20` | Items per page in `ls` |

### `movie-cli movie scan [folder]`

Scan a directory for video files, clean filenames, fetch TMDb metadata, and save to the database. Falls back to `scan_dir` config if no folder is given. Skips duplicates by file path.

```bash
movie-cli movie scan ~/Downloads
# 📁 Scanning: /home/user/Downloads
# 🎬 Found: Inception (2010) — ★ 8.4
# 📺 Found: Breaking Bad S01E01 — ★ 9.5
# ✅ Done: 15 files, 12 movies, 3 TV shows
```

### `movie-cli movie ls`

Paginated, interactive list of file-backed media (items with a local file on disk). Records from `search` or `info` without files are excluded. Navigate with `N`/`P`/`Q` or enter a number for detail view.

```bash
movie-cli movie ls
#  1. 🎬 Inception (2010)               ★ 8.4
#  2. 🎬 The Dark Knight (2008)         ★ 9.0
#  3. 📺 Breaking Bad (2008)            ★ 9.5
# Page 1/3 — [N]ext [P]rev [Q]uit or number:
```

### `movie-cli movie search <name>`

Search TMDb live, select a result, fetch full details + poster, and save to the database. Does **not** require a local file — catalogs metadata only.

```bash
movie-cli movie search "Inception"
#  1. 🎬 Inception (2010) ★ 8.4
#  2. 🎬 Inception: The Cobol Job (2010) ★ 7.2
# Select (0 to cancel): 1
# ✅ Saved: Inception (2010)
```

### `movie-cli movie info <id|title>`

Show detailed metadata for a media item. Looks up by numeric ID or title string in the local DB first, then falls back to TMDb API search (auto-saves if found).

```bash
movie-cli movie info 1
movie-cli movie info "Inception"
# 🎬 Inception (2010)
# ★ IMDb 8.8 / TMDb 8.4
# 🎭 Action, Sci-Fi, Thriller
# 🎬 Christopher Nolan
# 👥 Leonardo DiCaprio, Joseph Gordon-Levitt, ...
```

### `movie-cli movie suggest [N]`

Get movie/TV recommendations based on your library's genre patterns or TMDb trending. Choose between Movie, TV, or Random categories interactively.

```bash
movie-cli movie suggest 5
# 🎯 Your top genres: Action, Sci-Fi, Thriller
#  1. 🎬 Tenet (2020) ★ 7.3 — Action, Sci-Fi
#  2. 🎬 Arrival (2016) ★ 7.9 — Drama, Sci-Fi
#  ...
```

### `movie-cli movie move [directory]`

Browse a directory, select a video file, and move it to a configured destination with a clean filename (`Title (Year).ext`). Supports cross-drive moves with automatic copy+delete fallback.

```bash
movie-cli movie move ~/Downloads
#  1. 🎬 Inception (2010)  [2.4 GB]
#  2. 📺 Breaking Bad S01  [1.1 GB]
# Select file: 1
# Move to: [1] Movies  [2] TV Shows  [3] Archive  [4] Custom
# ✅ Moved → ~/Movies/Inception (2010).mkv
```

### `movie-cli movie rename`

Batch rename all library files to clean format (`Title (Year).ext`). Shows a preview and asks for confirmation before renaming.

```bash
movie-cli movie rename
# Renames to apply:
#   inception.2010.bluray.mkv → Inception (2010).mkv
#   the.dark.knight.mp4 → The Dark Knight (2008).mp4
# Apply 2 renames? [y/N]: y
# ✅ 2/2 files renamed
```

### `movie-cli movie undo`

Revert the most recent move or rename operation. Moves the file back to its original location.

```bash
movie-cli movie undo
# Last operation: ~/Downloads/inception.mkv → ~/Movies/Inception (2010).mkv
# Undo this? [y/N]: y
# ✅ Moved back to ~/Downloads/inception.mkv
```

### `movie-cli movie play <id>`

Open a media file with your system's default video player (macOS: `open`, Linux: `xdg-open`, Windows: `start`).

```bash
movie-cli movie play 1
# ▶️ Playing: Inception (2010)
```

### `movie-cli movie stats`

Display library statistics including counts, storage usage, genre distribution chart, and average ratings.

```bash
movie-cli movie stats
# 📊 Library Stats
# 🎬 Movies: 42  📺 TV Shows: 8  📦 Total: 50
# 💾 Storage: 185.3 GB total, avg 3.7 GB, largest 8.2 GB
# 🎭 Top Genres:
#   Action    ████████████████████████████ 28
#   Sci-Fi    ████████████████████ 20
#   Drama     ████████████████ 16
# ★ Avg IMDb: 7.4 / Avg TMDb: 7.1
```

### `movie-cli movie tag`

Manage user-defined tags on media items.

```bash
movie-cli movie tag add 1 favorite        # Add a tag
movie-cli movie tag remove 1 favorite     # Remove a tag
movie-cli movie tag list 1                # List tags for a media item
movie-cli movie tag list                  # List all tags with counts
# favorite (3), watchlist (7), rewatch (2)
```

## Project Structure

```
movie-cli/
├── main.go                        # Entry point
├── cmd/                           # Cobra commands (one file per command)
│   ├── root.go                    # Root command, registers subcommands
│   ├── hello.go                   # movie-cli hello
│   ├── version.go                 # movie-cli version
│   ├── update.go                  # movie-cli self-update
│   ├── movie.go                   # Parent: movie-cli movie
│   ├── movie_config.go            # config get/set
│   ├── movie_scan.go              # scan folder
│   ├── movie_ls.go                # paginated list
│   ├── movie_search.go            # TMDb search
│   ├── movie_info.go              # detail view + shared fetch helpers
│   ├── movie_suggest.go           # recommendations
│   ├── movie_move.go              # interactive move
│   ├── movie_move_helpers.go      # move utility functions
│   ├── movie_rename.go            # batch rename
│   ├── movie_undo.go              # undo last move/rename
│   ├── movie_play.go              # open with default player
│   ├── movie_stats.go             # library statistics
│   └── movie_resolve.go           # shared ID/title resolver
├── cleaner/cleaner.go             # Filename cleaning + slug generation
├── tmdb/client.go                 # TMDb API client
├── db/
│   ├── db.go                      # SQLite connection + migrations
│   ├── media.go                   # Media CRUD operations
│   ├── config.go                  # Config get/set
│   ├── history.go                 # Move history + scan history
│   └── helpers.go                 # String utilities
├── updater/updater.go             # Git-based self-update
├── version/version.go             # Build-time version variables
├── Makefile                       # Build targets
├── build.ps1                      # PowerShell build + deploy
├── spec.md                        # Full project specification
└── spec/                          # Specs and issue tracking
    ├── 01-app/                    # Application specs
    └── 02-app/issues/             # Issue write-ups
```

## Build

```bash
make build              # Current OS
make build-windows      # Windows amd64
make build-mac-arm      # macOS ARM64
make build-mac-intel    # macOS amd64
make build-linux        # Linux amd64
make install            # Build + copy to /usr/local/bin
```

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/spf13/cobra` | CLI framework |
| `modernc.org/sqlite` | Pure-Go SQLite driver (no CGo) |

## Data Storage

All data lives in `./data/`:

```
./data/
├── movie-cli.db              # SQLite database (WAL mode)
├── thumbnails/           # Downloaded poster images
└── json/history/         # Move operation logs (RFC3339 timestamps)
```

## License

Private project.
