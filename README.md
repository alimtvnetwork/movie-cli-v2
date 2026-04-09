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
- **Tag** media with custom labels for organization
- **Stats** with genre charts, storage usage, and average ratings
- **Self-update** via `git pull --ff-only`

## Commands

```
movie
├── hello                         # Greeting with version
├── version                       # Version, commit, build date, Go, OS/arch
├── changelog [--latest]          # Show changelog (full or latest version)
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

---

## Installation

### Option 1 — Quick Install from GitHub Release

Downloads the latest release binary, verifies SHA256 checksums, installs to your PATH. No Go or Git required.

**Windows (PowerShell)**
```powershell
irm https://github.com/mahin/movie-cli-v2/releases/latest/download/install.ps1 | iex
```

**Linux / macOS**
```bash
curl -fsSL https://github.com/mahin/movie-cli-v2/releases/latest/download/install.sh | bash
```

**Install options:**

| Flag | PowerShell | Bash | Default |
|------|-----------|------|---------|
| Install directory | `-InstallDir C:\tools\movie` | `--dir ~/bin` | `%LOCALAPPDATA%\movie` (Win) / `~/.local/bin` (Unix) |
| Force architecture | `-Arch arm64` | `--arch arm64` | Auto-detect |
| Skip PATH update | `-NoPath` | `--no-path` | Adds to PATH |

### Option 2 — Build from Source

**Prerequisites:**

| Requirement | Minimum | Check |
|---|---|---|
| **Go** | 1.22+ | `go version` |
| **Git** | 2.x | `git --version` |
| **PowerShell** | 5.1+ (Win) / 7+ (Unix) | `$PSVersionTable.PSVersion` |

**Windows (PowerShell)**
```powershell
git clone https://github.com/mahin/movie-cli-v2.git; cd movie-cli-v2; .\run.ps1
```

**macOS / Linux**
```bash
git clone https://github.com/mahin/movie-cli-v2.git && cd movie-cli-v2 && pwsh run.ps1
```

**Using the bootstrap installer:**
```powershell
pwsh install.ps1                      # Fresh install (clone + build + deploy)
pwsh install.ps1 -DeployPath ~/bin    # Custom deploy path
```

### Verify

```bash
movie version
# v1.0.0 (commit: abc1234, built: 2026-04-09)
#   Go:   go1.22.0
#   OS:   linux/amd64
```

> **Tip**: If `movie` is not found, add the deploy path to your `PATH`.
> Default: `E:\bin-run` (Windows) or `/usr/local/bin` (Unix) for source builds.

---

## Quick Start

```bash
# Set your TMDb API key
movie movie config set tmdb_api_key YOUR_KEY

# Scan a folder
movie movie scan ~/Downloads

# Browse your library
movie movie ls

# Search TMDb directly
movie movie search "Inception"

# Get suggestions
movie movie suggest 5
```

---

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

## Release Workflow

Releases are fully automated via GitHub Actions. Pushing to a `release/**` branch or a `v*` tag triggers:

1. **Cross-compilation** — 6 binaries (Windows/Linux/macOS × amd64/arm64)
2. **Packaging** — `.zip` (Windows) and `.tar.gz` (Unix)
3. **SHA256 checksums** — `checksums.txt` with all artifact hashes
4. **Install scripts** — version-pinned `install.ps1` and `install.sh`
5. **GitHub Release** — formatted page with changelog, checksums, and install instructions

### Creating a Release

```bash
# Option A: Push a release branch
git checkout -b release/v1.3.0
git push origin release/v1.3.0

# Option B: Tag directly
git tag v1.3.0
git push origin v1.3.0
```

Both trigger the same pipeline. Version is resolved from the ref name.

See [spec/pipeline/01-release-pipeline.md](spec/pipeline/01-release-pipeline.md) for the full pipeline spec.

---

## Command Reference

### `movie hello`

Print a greeting with the current version.

```bash
movie hello
# 👋 Hello from movie-cli-v2!
#    Running version: v1.2.0
```

### `movie version`

Show version, commit hash, build date, Go version, and OS/architecture.

```bash
movie version
# movie v1.2.0 (commit: abc1234, built: 2024-06-01)
#   Go:   go1.22.0
#   OS:   darwin/arm64
```

### `movie self-update`

Pull latest code via `git pull --ff-only`. Requires a clean working tree.

```bash
movie self-update
# ✅ Updated abc1234 → def5678
```

### `movie movie config`

View or update configuration settings. API keys are masked in output.

```bash
movie movie config                          # Show all settings
movie movie config get movies_dir           # Get a single key
movie movie config set movies_dir ~/Movies  # Set a key
movie movie config set tmdb_api_key KEY     # Set API key
movie movie config set page_size 30         # Items per page
```

| Key | Default | Purpose |
|---|---|---|
| `movies_dir` | `~/Movies` | Movie file destination |
| `tv_dir` | `~/TVShows` | TV show destination |
| `archive_dir` | `~/Archive` | Archive destination |
| `scan_dir` | `~/Downloads` | Default scan source |
| `tmdb_api_key` | *(none)* | TMDb API key |
| `page_size` | `20` | Items per page in `ls` |

### `movie movie scan [folder]`

Scan a directory for video files, clean filenames, fetch TMDb metadata, and save to the database. Falls back to `scan_dir` config if no folder is given. Skips duplicates by file path.

```bash
movie movie scan ~/Downloads
# 📁 Scanning: /home/user/Downloads
# 🎬 Found: Inception (2010) — ★ 8.4
# 📺 Found: Breaking Bad S01E01 — ★ 9.5
# ✅ Done: 15 files, 12 movies, 3 TV shows
```

### `movie movie ls`

Paginated, interactive list of file-backed media (items with a local file on disk). Records from `search` or `info` without files are excluded. Navigate with `N`/`P`/`Q` or enter a number for detail view.

```bash
movie movie ls
#  1. 🎬 Inception (2010)               ★ 8.4
#  2. 🎬 The Dark Knight (2008)         ★ 9.0
#  3. 📺 Breaking Bad (2008)            ★ 9.5
# Page 1/3 — [N]ext [P]rev [Q]uit or number:
```

### `movie movie search <name>`

Search TMDb live, select a result, fetch full details + poster, and save to the database. Does **not** require a local file — catalogs metadata only.

```bash
movie movie search "Inception"
#  1. 🎬 Inception (2010) ★ 8.4
#  2. 🎬 Inception: The Cobol Job (2010) ★ 7.2
# Select (0 to cancel): 1
# ✅ Saved: Inception (2010)
```

### `movie movie info <id|title>`

Show detailed metadata for a media item. Looks up by numeric ID or title string in the local DB first, then falls back to TMDb API search (auto-saves if found).

```bash
movie movie info 1
movie movie info "Inception"
# 🎬 Inception (2010)
# ★ IMDb 8.8 / TMDb 8.4
# 🎭 Action, Sci-Fi, Thriller
# 🎬 Christopher Nolan
# 👥 Leonardo DiCaprio, Joseph Gordon-Levitt, ...
```

### `movie movie suggest [N]`

Get movie/TV recommendations based on your library's genre patterns or TMDb trending. Choose between Movie, TV, or Random categories interactively.

```bash
movie movie suggest 5
# 🎯 Your top genres: Action, Sci-Fi, Thriller
#  1. 🎬 Tenet (2020) ★ 7.3 — Action, Sci-Fi
#  2. 🎬 Arrival (2016) ★ 7.9 — Drama, Sci-Fi
#  ...
```

### `movie movie move [directory]`

Browse a directory, select a video file, and move it to a configured destination with a clean filename (`Title (Year).ext`). Supports cross-drive moves with automatic copy+delete fallback.

```bash
movie movie move ~/Downloads
#  1. 🎬 Inception (2010)  [2.4 GB]
#  2. 📺 Breaking Bad S01  [1.1 GB]
# Select file: 1
# Move to: [1] Movies  [2] TV Shows  [3] Archive  [4] Custom
# ✅ Moved → ~/Movies/Inception (2010).mkv
```

### `movie movie rename`

Batch rename all library files to clean format (`Title (Year).ext`). Shows a preview and asks for confirmation before renaming.

```bash
movie movie rename
# Renames to apply:
#   inception.2010.bluray.mkv → Inception (2010).mkv
#   the.dark.knight.mp4 → The Dark Knight (2008).mp4
# Apply 2 renames? [y/N]: y
# ✅ 2/2 files renamed
```

### `movie movie undo`

Revert the most recent move or rename operation. Moves the file back to its original location.

```bash
movie movie undo
# Last operation: ~/Downloads/inception.mkv → ~/Movies/Inception (2010).mkv
# Undo this? [y/N]: y
# ✅ Moved back to ~/Downloads/inception.mkv
```

### `movie movie play <id>`

Open a media file with your system's default video player (macOS: `open`, Linux: `xdg-open`, Windows: `start`).

```bash
movie movie play 1
# ▶️ Playing: Inception (2010)
```

### `movie movie stats`

Display library statistics including counts, storage usage, genre distribution chart, and average ratings.

```bash
movie movie stats
# 📊 Library Stats
# 🎬 Movies: 42  📺 TV Shows: 8  📦 Total: 50
# 💾 Storage: 185.3 GB total, avg 3.7 GB, largest 8.2 GB
# 🎭 Top Genres:
#   Action    ████████████████████████████ 28
#   Sci-Fi    ████████████████████ 20
#   Drama     ████████████████ 16
# ★ Avg IMDb: 7.4 / Avg TMDb: 7.1
```

### `movie movie tag`

Manage user-defined tags on media items.

```bash
movie movie tag add 1 favorite        # Add a tag
movie movie tag remove 1 favorite     # Remove a tag
movie movie tag list 1                # List tags for a media item
movie movie tag list                  # List all tags with counts
# favorite (3), watchlist (7), rewatch (2)
```

### `movie changelog`

Show the project changelog. Prints the full changelog by default, or only the latest version block with `--latest`.

```bash
movie changelog              # Full changelog
movie changelog --latest     # Latest version only
# ## v1.0.0
# ### Added
# - Batch move, JSON export, genre-based discovery, ...
```

---

## Project Structure

```
movie-cli-v2/
├── main.go                        # Entry point
├── cmd/                           # Cobra commands (one file per command)
│   ├── root.go                    # Root command, registers subcommands
│   ├── hello.go                   # movie hello
│   ├── version.go                 # movie version
│   ├── update.go                  # movie self-update
│   ├── movie.go                   # Parent: movie movie
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
│   ├── movie_tag.go               # tag management
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
├── .github/workflows/
│   └── release.yml                # Release pipeline (cross-compile + GitHub Release)
├── run.ps1                        # PowerShell build + deploy pipeline
├── install.ps1                    # Bootstrap installer (clone + build)
├── CHANGELOG.md                   # Release notes
├── spec.md                        # Full project specification
└── spec/                          # Detailed specs
    ├── pipeline/                  # CI/CD pipeline specs
    ├── 01-coding-guidelines/      # Code style
    ├── 02-error-manage-spec/      # Error handling
    ├── 03-general/                # Build, install, config guides
    └── 08-app/                    # Application specs
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
├── movie.db              # SQLite database (WAL mode)
├── thumbnails/               # Downloaded poster images
└── json/history/             # Move operation logs (RFC3339 timestamps)
```

## License

Private project.
